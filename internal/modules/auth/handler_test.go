// internal/modules/auth/handler/handler.go
package auth_test

import (
	"encoding/json"
	"errors"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/auth"
	"net/http"

	dtoAuth "github.com/kiminodare/HOVARLAY-BE/internal/modules/auth/dto"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
)

type Handler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dtoAuth.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(r.Context(), &req) // ✅ Pass pointer
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dtoUser.Request // ✅ Bukan pointer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), &req) // ✅ Pass pointer
	if err != nil {
		if errors.Is(err, utils.ErrEmailAlreadyExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}
