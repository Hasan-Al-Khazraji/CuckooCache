package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/orchestrator"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:7000", "listen address")
	workers := flag.String("workers", "127.0.0.1:7101,127.0.0.1:7102", "comma-separated worker addrs")
	reqMs := flag.Int("req_ms", 300, "request timeout ms")
	idleMs := flag.Int("idle_ms", 60000, "idle timeout ms")
	maxKey := flag.Int("max_key", 4096, "max key bytes")
	maxVal := flag.Int("max_val", 1<<20, "max value bytes")
	flag.Parse()

	cfg := orchestrator.Config{
		ListenAddr:    *listen,
		RequestTO:     time.Duration(*reqMs) * time.Millisecond,
		IdleTimeout:   time.Duration(*idleMs) * time.Millisecond,
		MaxKeyBytes:   uint16(*maxKey),
		MaxValueBytes: uint32(*maxVal),
		Workers:       strings.Split(*workers, ","),
	}

	o := orchestrator.New(cfg)
	log.Printf("orchestrator listening on %s", cfg.ListenAddr)
	if err := o.ListenAndServe(); err != nil {
		log.Fatalf("orchestrator error: %v", err)
	}
}
