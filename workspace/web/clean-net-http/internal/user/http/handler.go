package httpuser

import (
	"clean-net-http/internal/user/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.UserService
}

func NewHandler(service *service.UserService) *Handler {
	return &Handler{service: service}
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.createUser)
	mux.HandleFunc("GET /users", h.listUsers)
	mux.HandleFunc("GET /users/{id}", h.getUser)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if req.Email == "" || req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and name are required"})
		return
	}

	user, err := h.service.CreateUser(ctx, req.Email, req.Name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		if h.service.IsNotFound(err) || errors.Is(err, serviceErrNotFoundCompat()) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := int32(50)
	offset := int32(0)

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n <= 0 || n > 200 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid limit"})
			return
		}
		limit = int32(n)
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid offset"})
			return
		}
		offset = int32(n)
	}

	users, err := h.service.ListUsers(ctx, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Tiny compatibility shim in case you later wrap not-found differently.
func serviceErrNotFoundCompat() error {
	return nil
}
