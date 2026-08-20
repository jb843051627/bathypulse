package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jb843051627/bathypulse/internal/clock"
	"github.com/jb843051627/bathypulse/internal/handler"
	"github.com/jb843051627/bathypulse/internal/service"
	"github.com/jb843051627/bathypulse/internal/store"
)

func main() {
	dbPath := os.Getenv("BATHYPULSE_DB")
	if dbPath == "" {
		dbPath = "data/bathypulse.db"
	}
	db, err := store.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	svc := service.New(db, clock.System{})
	svc.Start()
	defer svc.Close()

	if len(os.Args) > 1 && os.Args[1] == "--smoke-test" {
		if err := smoke(context.Background(), svc); err != nil {
			panic(err)
		}
		fmt.Println("smoke test passed")
		return
	}

	addr := os.Getenv("BATHYPULSE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8097"
	}
	srv := handler.New(svc)
	server := &http.Server{Addr: addr, Handler: srv.Routes(), ReadHeaderTimeout: 5 * time.Second}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func smoke(ctx context.Context, svc *service.ObservatoryService) error {
	station := service.SeedStation("smoke-station", "north-trench")
	if _, err := svc.RegisterStation(ctx, station); err != nil {
		return err
	}
	return svc.Health(ctx)
}
