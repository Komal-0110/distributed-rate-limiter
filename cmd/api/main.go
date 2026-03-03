package main

import (
	"log"
	"net/http"
	"rate-limiter/internal/auth/infrastructure"
	"rate-limiter/internal/auth/repository"
	"rate-limiter/internal/auth/service"
	"rate-limiter/internal/handler"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}).Methods("GET")

	jwtProvider := infrastructure.NewJWTProvider("super-secret-key", "rate-limiter")

	userRepo := repository.NewUsers()
	sessionRepo := repository.NewSessions()
	credentialRepo := repository.NewCredentials()

	authService := service.NewAuthService(
		userRepo,
		credentialRepo,
		sessionRepo,
		jwtProvider,
	)

	authHandler := handler.NewAuthHandler(authService)

	r.HandleFunc("/register", authHandler.Register).Methods("POST")

	log.Println("server is started")

	err := http.ListenAndServe(":8080", r)

	log.Println("existed server", err)
}
