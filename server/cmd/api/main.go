package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/julianstephens/connectify/server/internal/config"
	"github.com/julianstephens/connectify/server/internal/db"
	"github.com/julianstephens/connectify/server/internal/handlers"
	"github.com/julianstephens/connectify/server/internal/middleware"
	"github.com/julianstephens/connectify/server/internal/store"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	config.Load()
	log.Printf("%+v", config.AppConfig)
	addr := "0.0.0.0:" + config.AppConfig.Port

	loggerType := config.AppConfig.LogType
	config.SetLogger(loggerType)

	var zapLogger *zap.Logger
	var logrusLogger *logrus.Logger

	switch loggerType {
	case "logrus":
		l := logrus.New()
		if config.AppConfig.LogFormat == "human" {
			l.SetFormatter(&logrus.TextFormatter{})
		} else {
			l.SetFormatter(&logrus.JSONFormatter{})
		}
		l.SetLevel(logrus.InfoLevel)
		logrusLogger = l
	default:
		var zl *zap.Logger
		var err error
		if config.AppConfig.LogFormat == "human" {
			zl, err = zap.NewDevelopment()
		} else {
			zl, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatalf("failed to initialize zap logger: %v", err)
		}
		zapLogger = zl
		defer func() {
			_ = zapLogger.Sync()
		}()
	}

	r := mux.NewRouter()

	// Request ID middleware first (ensures other middleware/handlers see it)
	r.Use(middleware.RequestIDMiddleware)

	// Add structured logging middleware depending on selected logger
	if logrusLogger != nil {
		r.Use(middleware.LoggingMiddlewareLogrus(logrusLogger))
	} else if zapLogger != nil {
		r.Use(middleware.LoggingMiddlewareZap(zapLogger))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := db.NewDB(ctx, fmt.Sprintf("host=localhost port=5432 user=dev password=%s dbname=connectify sslmode=disable", config.AppConfig.DBPassword))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthGuard)

	api.HandleFunc("/health", HealthHandler).Methods(http.MethodGet)

	// Post routes group
	postRouter := api.PathPrefix("/posts").Subrouter()
	postStore := store.NewPostStore(db)
	postHandler := handlers.NewPostHandler(postStore)
	postsHandler := handlers.NewPostsHandler(postStore)
	postRouter.HandleFunc("/", postHandler.CreatePost).Methods(http.MethodPost)
	postRouter.HandleFunc("/{id}", postHandler.GetPost).Methods(http.MethodGet)
	postRouter.HandleFunc("/", postsHandler.PostsCursorHandler).Methods(http.MethodGet)

	// Post media routes group
	postMediaRouter := api.PathPrefix("/assets").Subrouter()
	postMediaStore := store.NewPostMediaStore(db)
	postMediaHandler := handlers.NewPostMediaHandler(postMediaStore)
	postMediaRouter.HandleFunc("/", postMediaHandler.UploadPostMedia).Methods(http.MethodPost)
	postMediaRouter.HandleFunc("/{id}", postMediaHandler.GetPostMedia).Methods(http.MethodGet)
	postMediaRouter.HandleFunc("/{id}", postMediaHandler.DeletePostMedia).Methods(http.MethodDelete)

	allowedOrigins := []string{"*"}
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: c.Handler(r),
	}

	// graceful shutdown
	go func() {
		log.Printf("starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	serverCtx, serverCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer serverCancel()
	if err := srv.Shutdown(serverCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}
