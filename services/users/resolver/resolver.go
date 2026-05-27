package resolver

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/packages/middleware"
	"backend/packages/service"
	usersservice "backend/services/users/service"
)

type Resolver struct {
	svc *usersservice.Service
}

func New(svc *usersservice.Service) *Resolver {
	return &Resolver{svc: svc}
}

func (r *Resolver) GraphQLHandler(w http.ResponseWriter, req *http.Request) {
	var gqlReq service.GQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		service.WriteGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	q := strings.TrimSpace(gqlReq.Query)

	switch {
	case strings.Contains(q, "me"):
		r.handleMe(w, req)
	case strings.Contains(q, "myBranches"):
		r.handleMyBranches(w, req)
	case strings.Contains(q, "switchBranch"):
		r.handleSwitchBranch(w, req)

	case strings.Contains(q, "myTenant"):
		r.handleMyTenant(w, req)
	case strings.Contains(q, "tenantSettings"):
		r.handleTenantSettings(w, req)
	case strings.Contains(q, "tenant"):
		if strings.Contains(q, "createTenant") {
			r.handleCreateTenant(w, req)
		} else if strings.Contains(q, "updateTenant") {
			r.handleUpdateTenant(w, req)
		} else {
			r.handleTenant(w, req)
		}

	case strings.Contains(q, "myProfile"):
		r.handleMyProfile(w, req)
	case strings.Contains(q, "updateProfile"):
		r.handleUpdateProfile(w, req, gqlReq)
	case strings.Contains(q, "updateTenantSettings"):
		r.handleUpdateTenantSettings(w, req, gqlReq)

	case strings.Contains(q, "inviteUser"):
		r.handleInviteUser(w, req, gqlReq)

	case strings.Contains(q, "roles"):
		r.handleRoles(w, req)
	case strings.Contains(q, "permissions"):
		r.handlePermissions(w, req)

	default:
		service.WriteGQL(w, map[string]interface{}{"__typename": "Query"})
	}
}

// ---- SSO ----

func (r *Resolver) handleMe(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req.Context())
	if userID == "" {
		service.WriteGQLError(w, 401, "authentication required")
		return
	}

	email := middleware.GetEmail(req.Context())
	tenantID := middleware.GetTenantID(req.Context())

	ctx, err := r.svc.GetUserContext(req.Context(), userID, email, "", tenantID)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"me": toMap(ctx)})
}

func (r *Resolver) handleMyBranches(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req.Context())
	if userID == "" {
		service.WriteGQLError(w, 401, "authentication required")
		return
	}

	branches, err := r.svc.GetUserBranches(req.Context(), userID)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"myBranches": toSlice(branches)})
}

func (r *Resolver) handleSwitchBranch(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req.Context())
	if userID == "" {
		service.WriteGQLError(w, 401, "authentication required")
		return
	}

	var gqlReq service.GQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		service.WriteGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vars := gqlReq.Variables
	tenantID, _ := vars["tenantId"].(string)
	if tenantID == "" {
		service.WriteGQLError(w, 400, "tenantId is required")
		return
	}

	email := middleware.GetEmail(req.Context())

	ctx, err := r.svc.SwitchBranch(req.Context(), userID, email, "", tenantID)
	if err != nil {
		service.WriteGQLError(w, 403, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"switchBranch": toMap(ctx)})
}

// ---- Tenant ----

func (r *Resolver) handleTenant(w http.ResponseWriter, req *http.Request) {
	var gqlReq service.GQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		service.WriteGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	vars := gqlReq.Variables
	id, _ := vars["id"].(string)
	if id == "" {
		service.WriteGQLError(w, 400, "id is required")
		return
	}
	tenant, err := r.svc.GetTenant(req.Context(), id)
	if err != nil {
		service.WriteGQLError(w, 404, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"tenant": toMap(tenant)})
}

func (r *Resolver) handleMyTenant(w http.ResponseWriter, req *http.Request) {
	tenantID := middleware.GetTenantID(req.Context())
	if tenantID == "" {
		service.WriteGQLError(w, 403, "tenant context required")
		return
	}
	tenant, err := r.svc.GetTenant(req.Context(), tenantID)
	if err != nil {
		service.WriteGQLError(w, 404, "tenant not found")
		return
	}
	service.WriteGQL(w, map[string]interface{}{"myTenant": toMap(tenant)})
}

func (r *Resolver) handleCreateTenant(w http.ResponseWriter, req *http.Request) {
	var gqlReq service.GQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		service.WriteGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vars := gqlReq.Variables
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		service.WriteGQLError(w, 400, "input is required")
		return
	}

	tenant, err := r.svc.CreateTenant(req.Context(), input)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"createTenant": toMap(tenant)})
}

func (r *Resolver) handleUpdateTenant(w http.ResponseWriter, req *http.Request) {
	var gqlReq service.GQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		service.WriteGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vars := gqlReq.Variables
	id, _ := vars["id"].(string)
	if id == "" {
		service.WriteGQLError(w, 400, "id is required")
		return
	}
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		service.WriteGQLError(w, 400, "input is required")
		return
	}

	tenant, err := r.svc.UpdateTenant(req.Context(), id, input)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"updateTenant": toMap(tenant)})
}

func (r *Resolver) handleTenantSettings(w http.ResponseWriter, req *http.Request) {
	tenantID := middleware.GetTenantID(req.Context())
	if tenantID == "" {
		service.WriteGQLError(w, 403, "tenant context required")
		return
	}
	settings, err := r.svc.GetTenantSettings(req.Context(), tenantID)
	if err != nil {
		service.WriteGQLError(w, 404, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"tenantSettings": toMap(settings)})
}

func (r *Resolver) handleUpdateTenantSettings(w http.ResponseWriter, req *http.Request, gqlReq service.GQLRequest) {
	tenantID := middleware.GetTenantID(req.Context())
	if tenantID == "" {
		service.WriteGQLError(w, 403, "tenant context required")
		return
	}

	vars := gqlReq.Variables
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		service.WriteGQLError(w, 400, "input is required")
		return
	}

	settings, err := r.svc.UpdateTenantSettings(req.Context(), tenantID, input)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"updateTenantSettings": toMap(settings)})
}

// ---- Profile ----

func (r *Resolver) handleMyProfile(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req.Context())
	if userID == "" {
		service.WriteGQLError(w, 401, "authentication required")
		return
	}

	profile, err := r.svc.GetProfile(req.Context(), userID)
	if err != nil {
		service.WriteGQLError(w, 404, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"myProfile": toMap(profile)})
}

func (r *Resolver) handleUpdateProfile(w http.ResponseWriter, req *http.Request, gqlReq service.GQLRequest) {
	userID := middleware.GetUserID(req.Context())
	if userID == "" {
		service.WriteGQLError(w, 401, "authentication required")
		return
	}

	vars := gqlReq.Variables
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		service.WriteGQLError(w, 400, "input is required")
		return
	}

	profile, err := r.svc.UpdateProfile(req.Context(), userID, input)
	if err != nil {
		service.WriteGQLError(w, 500, err.Error())
		return
	}
	service.WriteGQL(w, map[string]interface{}{"updateProfile": toMap(profile)})
}

// ---- Invite ----

func (r *Resolver) handleInviteUser(w http.ResponseWriter, req *http.Request, gqlReq service.GQLRequest) {
	service.WriteGQL(w, map[string]interface{}{"inviteUser": true})
}

// ---- Roles & Permissions ----

func (r *Resolver) handleRoles(w http.ResponseWriter, req *http.Request) {
	roles, err := r.svc.ListRoles(req.Context())
	if err != nil {
		service.WriteGQLError(w, 500, "failed to fetch roles")
		return
	}
	result := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		result[i] = toMap(role)
	}
	service.WriteGQL(w, map[string]interface{}{"roles": result})
}

func (r *Resolver) handlePermissions(w http.ResponseWriter, req *http.Request) {
	perms, err := r.svc.ListPermissions(req.Context())
	if err != nil {
		service.WriteGQLError(w, 500, "failed to fetch permissions")
		return
	}
	result := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		result[i] = toMap(p)
	}
	service.WriteGQL(w, map[string]interface{}{"permissions": result})
}

func toMap(v interface{}) map[string]interface{} {
	data, _ := json.Marshal(v)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return m
}

func toSlice(v interface{}) []interface{} {
	data, _ := json.Marshal(v)
	var s []interface{}
	json.Unmarshal(data, &s)
	return s
}
