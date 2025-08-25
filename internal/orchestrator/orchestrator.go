package orchestrator

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/partition"
	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/partition/consistent"
	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/partition/static"
	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/proto"
)

type Config struct {
	ListenAddr    string
	RequestTO     time.Duration
	IdleTimeout   time.Duration
	MaxKeyBytes   uint16
	MaxValueBytes uint32
	Workers       []string
	Partitioner   string // static, consistent
	Vnodes        int
}

type Orchestrator struct {
	cfg  Config
	part partition.Partitioner
}

func New(cfg Config) *Orchestrator {
	var p partition.Partitioner
	switch cfg.Partitioner {
	case "consistent":
		p = consistent.New(cfg.Workers, cfg.Vnodes)
	default:
		p = static.New(cfg.Workers)
	}
	return &Orchestrator{
		cfg:  cfg,
		part: p,
	}
}

func (o *Orchestrator) ListenAndServe() error {
	ln, err := net.Listen("tcp", o.cfg.ListenAddr)
	if err != nil {
		return err
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		go o.handleConn(c)
	}
}

func (o *Orchestrator) handleConn(c net.Conn) {
	defer c.Close()
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	for {
		if err := c.SetReadDeadline(time.Now().Add(o.cfg.IdleTimeout)); err != nil {
			return
		}

		req, err := proto.ReadRequest(br, o.cfg.MaxKeyBytes, o.cfg.MaxValueBytes)
		if err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), o.cfg.RequestTO)
		resp := o.route(ctx, req)
		cancel()

		if err := proto.WriteResponse(bw, resp); err != nil {
			return
		}

		if err := bw.Flush(); err != nil {
			return
		}
	}
}

func (o *Orchestrator) route(ctx context.Context, req *proto.Request) *proto.Response {
	owners := o.part.OwnersFor(string(req.Key))
	switch req.Op {
	case proto.OpGet:
		if r := o.callWorker(ctx, owners.Primary, req); r.Status == proto.StatusOK || r.Status == proto.StatusNotFound {
			return r
		}
		return o.callWorker(ctx, owners.Secondary, req)
	case proto.OpSet:
		r1c := make(chan *proto.Response, 1)
		r2c := make(chan *proto.Response, 1)
		go func() { r1c <- o.callWorker(ctx, owners.Primary, req) }()
		go func() { r2c <- o.callWorker(ctx, owners.Secondary, req) }()
		r1, r2 := <-r1c, <-r2c
		if r1.Status == proto.StatusOK && r2.Status == proto.StatusOK {
			return proto.OK(nil)
		}
		return proto.Err()
	default:
		return proto.Err()
	}
}

func (o *Orchestrator) callWorker(ctx context.Context, addr string, req *proto.Request) *proto.Response {
	d := &net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return proto.Err()
	}
	defer conn.Close()

	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	if err := proto.WriteRequest(bw, req); err != nil {
		return proto.Err()
	}
	if err := bw.Flush(); err != nil {
		return proto.Err()
	}

	resp, err := proto.ReadResponse(br, o.cfg.MaxValueBytes)
	if err != nil {
		return proto.Err()
	}
	return resp
}
