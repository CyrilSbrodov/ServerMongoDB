package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"server/internal/user"
	"server/internal/user/repository"
	"server/pkg/client/mongodb"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	return &App{
		}
}

func (a *App) Run() {

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	client, err := mongodb.NewClient(context.Background(), "server")
	if err != nil {
		log.Fatal(err)
	}
	storage := repository.NewStorage(client, "users")
	service := user.NewService(storage)
	handler := user.NewHandler(service)
	handler.Register(router)

	a.server = &http.Server{
		Addr: ":9090",
		Handler: router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server is listen on port 9090")

	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}