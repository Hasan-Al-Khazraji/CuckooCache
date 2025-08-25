package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/proto"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:7000", "orchestrator address")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalf("usage: ccctl [--addr host:port] get <key> | set <key> <value>")
	}

	op := args[0]
	key := []byte(args[1])
	var value []byte
	if op == "set" {
		if len(args) < 3 {
			log.Fatalf("usage: ccctl set <key> <value>")
		}
		value = []byte(args[2])
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	var req *proto.Request
	switch op {
	case "get":
		req = &proto.Request{Op: proto.OpGet, Key: key}
	case "set":
		req = &proto.Request{Op: proto.OpSet, Key: key, Value: value}
	default:
		log.Fatalf("unknown op: %s", op)
	}

	if err := proto.WriteRequest(bw, req); err != nil {
		log.Fatalf("write request: %v", err)
	}
	if err := bw.Flush(); err != nil {
		log.Fatalf("flush: %v", err)
	}

	res, err := proto.ReadResponse(br, ^uint32(0))
	if err != nil {
		log.Fatalf("read response: %v", err)
	}

	switch res.Status {
	case proto.StatusOK:
		if op == "get" {
			fmt.Printf("%s\n", string(res.Value))
		} else {
			fmt.Println("OK")
		}
	case proto.StatusNotFound:
		fmt.Println("NOT_FOUND")
	default:
		fmt.Printf("ERROR: %d\n", res.Status)
	}
}
