package handler

import (
	"encoding/json"
	"net/http"

	"backend/services/auth/model"
	authservice "backend/services/auth/service"
	"github.com/go-chi/chi/v5"
)

type RESTHandler struct {
	auth authservice.AuthService
}

func NewRESTHandler(auth authservice.AuthService) *RESTHandler {
	return &RESTHandler{auth: auth}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     *string `json:"name,omitempty"`
	App      string  `json:"app,omitempty"`
}

type changePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type addAppRequest struct {
	Role string `json:"role"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *RESTHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	input := model.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.auth.Login(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *RESTHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	input := model.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	resp, err := h.auth.Register(r.Context(), input, req.App)
	if err != nil {
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *RESTHandler) AddUserTenant(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	tenantID := chi.URLParam(r, "tenantId")
	if email == "" || tenantID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing email or tenant id"})
		return
	}
	if err := h.auth.AddUserTenant(r.Context(), email, tenantID); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "tenant added"})
}

func (h *RESTHandler) RemoveUserTenant(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	tenantID := chi.URLParam(r, "tenantId")
	if email == "" || tenantID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing email or tenant id"})
		return
	}
	if err := h.auth.RemoveUserTenant(r.Context(), email, tenantID); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "tenant removed"})
}

func (h *RESTHandler) AddUserApp(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	appID := chi.URLParam(r, "appId")
	if userID == "" || appID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing user id or app id"})
		return
	}
	var req addAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	role := req.Role
	if role == "" {
		role = "user"
	}
	if err := h.auth.AddUserApp(r.Context(), userID, appID, role); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "app added"})
}

func (h *RESTHandler) RemoveUserApp(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	appID := chi.URLParam(r, "appId")
	if userID == "" || appID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing user id or app id"})
		return
	}
	if err := h.auth.RemoveUserApp(r.Context(), userID, appID); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "app removed"})
}

func (h *RESTHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing user id"})
		return
	}

	if err := h.auth.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "current password is incorrect" || err.Error() == "user not found" {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "password changed"})
}
