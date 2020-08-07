package main

import (
	"chat/internal/postgres"
	"chat/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
)

const (
	port  = ":9000"
	delay = 5
)

func main() {
	newLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Could not instantiate log %+s", err)
	}

	db := postgres.New(newLogger)

	defer db.Close()

	repoUser, err := postgres.NewUserStorage(db)
	if err != nil {
		newLogger.Fatalf("failed to create user storage %+s", err)
	}

	repoMsg, err := postgres.NewMsgStorage(db)
	if err != nil {
		newLogger.Fatalf("failed to create session storage %+s", err)
	}

	repoChat, err := postgres.NewChatStorage(db)
	if err != nil {
		newLogger.Fatalf("failed to create session storage %+s", err)
	}

	//templates := ParseTemplates()
	handler := newHandler(newLogger, repoUser, repoMsg, repoChat)

	r := chi.NewRouter()

	handler.Routers(r)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			newLogger.Fatalf("server stopped %+s", err)
		}
	}()

	newLogger.Debugf("server started")

	<-ctx.Done()

	newLogger.Debugf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), delay*time.Second)

	defer func() {
		cancel()
	}()

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		newLogger.Fatalf("server shutdown failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
