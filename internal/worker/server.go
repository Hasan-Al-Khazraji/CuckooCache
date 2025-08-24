package worker

import (
	"bufio"
	"net"
	"time"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/proto"
)

type KVStore interface {
	Get(key []byte) (val []byte, ok bool)
	Put(key, val []byte) (evictedKey, evictedVal []byte, evicted bool)
	Delete(key []byte) bool
	Len() int
}

type Config struct {
	ListenAddr    string
	IdleTimeout   time.Duration
	MaxKeyBytes   uint16
	MaxValueBytes uint32
}

type Server struct {
	cfg   Config
	store KVStore
	ln    net.Listener
}

func New(cfg Config, store KVStore) *Server {
	return &Server{cfg: cfg, store: store}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	for {
		c, err := s.ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(c)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	for {
		if err := c.SetReadDeadline(time.Now().Add(s.cfg.IdleTimeout)); err != nil {
			return
		}

		req, err := proto.ReadRequest(br, s.cfg.MaxKeyBytes, s.cfg.MaxValueBytes)
		if err != nil {
			return
		}

		switch req.Op {
		case proto.OpGet:
			if val, ok := s.store.Get(req.Key); ok {
				if err := proto.WriteResponse(bw, proto.OK(val)); err != nil {
					return
				}
			} else {
				if err := proto.WriteResponse(bw, proto.NotFound()); err != nil {
					return
				}
			}
		case proto.OpSet:
			s.store.Put(req.Key, req.Value)
			if err := proto.WriteResponse(bw, proto.OK(nil)); err != nil {
				return
			}
		default:
			if err := proto.WriteResponse(bw, proto.Err()); err != nil {
				return
			}
		}
		if err := bw.Flush(); err != nil {
			return
		}
	}
}
