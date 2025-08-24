package main

import (
	"flag"
	"log"
	"time"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/store/memstore"
	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/worker"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:7101", "listen address")
	capacity := flag.Int("cap", 100000, "cache capacity (items)")
	idleMs := flag.Int("idle_ms", 60000, "idle timeout ms")
	maxKey := flag.Int("max_key", 4096, "max key bytes")
	maxVal := flag.Int("max_val", 1<<20, "max value bytes")
	flag.Parse()

	store := memstore.New(*capacity)
	srv := worker.New(worker.Config{
		ListenAddr:    *listen,
		IdleTimeout:   time.Duration(*idleMs) * time.Millisecond,
		MaxKeyBytes:   uint16(*maxKey),
		MaxValueBytes: uint32(*maxVal),
	}, store)

	log.Printf("worker listening on %s", *listen)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
