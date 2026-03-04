package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"rate-limiter/internal/auth/service"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("failed to write"); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("email and password required"); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	resp, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("failed to write json response")
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("failed to write"); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("email and password required"); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	userAgent := r.UserAgent()

	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	}

	resp, err := h.authService.Login(
		req.Email,
		req.Password,
		userAgent,
		ipAddress,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Println("failed to write json response")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}); err != nil {
		log.Println("failed to write json response")
	}
}
