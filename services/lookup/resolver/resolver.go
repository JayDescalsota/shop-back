package resolver

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/services/lookup/model"
	"backend/services/lookup/repository"

	"github.com/google/uuid"
)

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

type GraphQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type Resolver struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Resolver {
	return &Resolver{repo: repo}
}

func writeGQL(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GraphQLResponse{Data: data})
}

func writeGQLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(GraphQLResponse{
		Errors: []GraphQLError{{Message: message}},
	})
}

func variables(req GraphQLRequest) map[string]interface{} {
	if req.Variables != nil {
		return req.Variables
	}
	return make(map[string]interface{})
}

func varString(vars map[string]interface{}, key string) string {
	if v, ok := vars[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func varID(vars map[string]interface{}, key string) string {
	return varString(vars, key)
}

func varInt(vars map[string]interface{}, key string) int {
	if v, ok := vars[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func varBool(vars map[string]interface{}, key string) *bool {
	if v, ok := vars[key]; ok {
		if b, ok := v.(bool); ok {
			return &b
		}
	}
	return nil
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

func (r *Resolver) GraphQLHandler(w http.ResponseWriter, req *http.Request) {
	var gqlReq GraphQLRequest
	if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
		writeGQLError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	q := strings.TrimSpace(gqlReq.Query)

	switch {
	case strings.Contains(q, "vehicleMakes"):
		r.handleVehicleMakes(w, req, gqlReq)
	case strings.Contains(q, "vehicleModels"):
		r.handleVehicleModels(w, req, gqlReq)
	case strings.Contains(q, "vehicleMake"):
		r.handleVehicleMake(w, req, gqlReq)
	case strings.Contains(q, "vehicleModel"):
		r.handleVehicleModel(w, req, gqlReq)
	case strings.Contains(q, "searchVehicles"):
		r.handleSearchVehicles(w, req, gqlReq)

	case strings.Contains(q, "serviceTypes"):
		r.handleServiceTypes(w, req, gqlReq)
	case strings.Contains(q, "serviceType"):
		r.handleServiceType(w, req, gqlReq)
	case strings.Contains(q, "diagnosticCodes"):
		r.handleDiagnosticCodes(w, req, gqlReq)
	case strings.Contains(q, "diagnosticCode"):
		r.handleDiagnosticCode(w, req, gqlReq)
	case strings.Contains(q, "searchDiagnosticCodes"):
		r.handleSearchDiagnosticCodes(w, req, gqlReq)

	case strings.Contains(q, "partCategories"):
		r.handlePartCategories(w, req, gqlReq)
	case strings.Contains(q, "partCategory"):
		r.handlePartCategory(w, req, gqlReq)
	case strings.Contains(q, "checkPartCompatibility"):
		r.handlePartCompatibility(w, req, gqlReq)

	case strings.Contains(q, "fuelTypes"):
		r.handleFuelTypes(w, req, gqlReq)
	case strings.Contains(q, "transmissionTypes"):
		r.handleTransmissionTypes(w, req, gqlReq)
	case strings.Contains(q, "engineTypes"):
		r.handleEngineTypes(w, req, gqlReq)

	case strings.Contains(q, "laborRateTiers"):
		r.handleLaborRateTiers(w, req, gqlReq)
	case strings.Contains(q, "laborRateTier"):
		r.handleLaborRateTier(w, req, gqlReq)
	case strings.Contains(q, "countries"):
		r.handleCountries(w, req, gqlReq)
	case strings.Contains(q, "currencies"):
		r.handleCurrencies(w, req, gqlReq)

	default:
		writeGQL(w, map[string]interface{}{"__typename": "Query"})
	}
}

// ---- Vehicle handlers ----

func (r *Resolver) handleVehicleMakes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	search := varString(vars, "search")
	isActive := varBool(vars, "isActive")

	makes, err := r.repo.GetVehicleMakes(req.Context(), search, isActive)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch vehicle makes")
		return
	}

	result := make([]map[string]interface{}, len(makes))
	for i, m := range makes {
		result[i] = toMap(m)
	}
	writeGQL(w, map[string]interface{}{"vehicleMakes": result})
}

func (r *Resolver) handleVehicleMake(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	make, err := r.repo.GetVehicleMake(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "vehicle make not found")
		return
	}
	writeGQL(w, map[string]interface{}{"vehicleMake": toMap(make)})
}

func (r *Resolver) handleVehicleModels(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	makeID := varID(vars, "makeId")
	search := varString(vars, "search")
	year := varInt(vars, "year")
	vehicleType := varString(vars, "vehicleType")

	if makeID == "" {
		writeGQLError(w, 400, "makeId is required")
		return
	}

	models, err := r.repo.GetVehicleModels(req.Context(), makeID, search, year, vehicleType)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch vehicle models")
		return
	}

	result := make([]map[string]interface{}, len(models))
	for i, m := range models {
		result[i] = toMap(m)
	}
	writeGQL(w, map[string]interface{}{"vehicleModels": result})
}

func (r *Resolver) handleVehicleModel(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	m, err := r.repo.GetVehicleModel(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "vehicle model not found")
		return
	}
	writeGQL(w, map[string]interface{}{"vehicleModel": toMap(m)})
}

func (r *Resolver) handleSearchVehicles(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	query := varString(vars, "query")
	limit := varInt(vars, "limit")
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	models, err := r.repo.SearchVehicles(req.Context(), query, limit)
	if err != nil {
		writeGQLError(w, 500, "search failed")
		return
	}

	result := make([]map[string]interface{}, len(models))
	for i, m := range models {
		item := toMap(m)
		if m.Make != nil {
			item["make"] = toMap(m.Make)
		}
		result[i] = item
	}
	writeGQL(w, map[string]interface{}{"searchVehicles": result})
}

func (r *Resolver) handleCreateVehicleMake(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	input, ok := vars["input"].(map[string]interface{})
	if !ok {
		writeGQLError(w, 400, "input is required")
		return
	}

	make := &model.VehicleMake{
		ID:   uuid.New().String(),
		Name: input["name"].(string),
		Slug: input["slug"].(string),
	}
	if v, ok := input["logoUrl"]; ok {
		s := v.(string)
		make.LogoURL = &s
	}
	if v, ok := input["country"]; ok {
		s := v.(string)
		make.Country = &s
	}
	if v, ok := input["foundedYear"]; ok {
		i := int(v.(float64))
		make.FoundedYear = &i
	}
	if v, ok := input["sortOrder"]; ok {
		make.SortOrder = int(v.(float64))
	}
	make.Metadata = "{}"

	err := r.repo.CreateVehicleMake(req.Context(), make)
	if err != nil {
		writeGQLError(w, 500, "failed to create vehicle make")
		return
	}
	writeGQL(w, map[string]interface{}{"createVehicleMake": toMap(make)})
}

// ---- Service handlers ----

func (r *Resolver) handleServiceTypes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	category := varString(vars, "category")
	isGlobal := varBool(vars, "isGlobal")
	tenantID := varString(vars, "tenantId")

	isGlobalBool := tenantID == "" && (isGlobal == nil || *isGlobal)

	types, err := r.repo.GetServiceTypes(req.Context(), category, isGlobalBool, tenantID)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch service types")
		return
	}

	result := make([]map[string]interface{}, len(types))
	for i, t := range types {
		result[i] = toMap(t)
	}
	writeGQL(w, map[string]interface{}{"serviceTypes": result})
}

func (r *Resolver) handleServiceType(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	st, err := r.repo.GetServiceType(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "service type not found")
		return
	}
	writeGQL(w, map[string]interface{}{"serviceType": toMap(st)})
}

func (r *Resolver) handleDiagnosticCodes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	system := varString(vars, "system")
	search := varString(vars, "search")

	codes, err := r.repo.GetDiagnosticCodes(req.Context(), system, search)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch diagnostic codes")
		return
	}

	result := make([]map[string]interface{}, len(codes))
	for i, c := range codes {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"diagnosticCodes": result})
}

func (r *Resolver) handleDiagnosticCode(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	dc, err := r.repo.GetDiagnosticCode(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "diagnostic code not found")
		return
	}
	writeGQL(w, map[string]interface{}{"diagnosticCode": toMap(dc)})
}

func (r *Resolver) handleSearchDiagnosticCodes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	query := varString(vars, "query")

	codes, err := r.repo.SearchDiagnosticCodes(req.Context(), query)
	if err != nil {
		writeGQLError(w, 500, "search failed")
		return
	}

	result := make([]map[string]interface{}, len(codes))
	for i, c := range codes {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"searchDiagnosticCodes": result})
}

// ---- Part handlers ----

func (r *Resolver) handlePartCategories(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	parentID := varString(vars, "parentId")

	cats, err := r.repo.GetPartCategories(req.Context(), parentID)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch part categories")
		return
	}

	result := make([]map[string]interface{}, len(cats))
	for i, c := range cats {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"partCategories": result})
}

func (r *Resolver) handlePartCategory(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	cat, err := r.repo.GetPartCategory(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "part category not found")
		return
	}
	writeGQL(w, map[string]interface{}{"partCategory": toMap(cat)})
}

func (r *Resolver) handlePartCompatibility(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	partSKU := varString(vars, "partSku")
	makeID := varString(vars, "makeId")
	modelID := varString(vars, "modelId")
	year := varInt(vars, "year")

	if partSKU == "" {
		writeGQLError(w, 400, "partSku is required")
		return
	}

	compats, err := r.repo.CheckPartCompatibility(req.Context(), partSKU, makeID, modelID, year)
	if err != nil {
		writeGQLError(w, 500, "compatibility check failed")
		return
	}

	result := make([]map[string]interface{}, len(compats))
	for i, c := range compats {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"checkPartCompatibility": result})
}

// ---- Enum handlers ----

func (r *Resolver) handleFuelTypes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	types, err := r.repo.GetFuelTypes(req.Context())
	if err != nil {
		writeGQLError(w, 500, "failed to fetch fuel types")
		return
	}
	result := make([]map[string]interface{}, len(types))
	for i, t := range types {
		result[i] = toMap(t)
	}
	writeGQL(w, map[string]interface{}{"fuelTypes": result})
}

func (r *Resolver) handleTransmissionTypes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	types, err := r.repo.GetTransmissionTypes(req.Context())
	if err != nil {
		writeGQLError(w, 500, "failed to fetch transmission types")
		return
	}
	result := make([]map[string]interface{}, len(types))
	for i, t := range types {
		result[i] = toMap(t)
	}
	writeGQL(w, map[string]interface{}{"transmissionTypes": result})
}

func (r *Resolver) handleEngineTypes(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	types, err := r.repo.GetEngineTypes(req.Context())
	if err != nil {
		writeGQLError(w, 500, "failed to fetch engine types")
		return
	}
	result := make([]map[string]interface{}, len(types))
	for i, t := range types {
		result[i] = toMap(t)
	}
	writeGQL(w, map[string]interface{}{"engineTypes": result})
}

func (r *Resolver) handleLaborRateTiers(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	tenantID := varString(vars, "tenantId")
	isGlobal := varBool(vars, "isGlobal")

	tiers, err := r.repo.GetLaborRateTiers(req.Context(), tenantID, isGlobal)
	if err != nil {
		writeGQLError(w, 500, "failed to fetch labor rate tiers")
		return
	}

	result := make([]map[string]interface{}, len(tiers))
	for i, t := range tiers {
		result[i] = toMap(t)
	}
	writeGQL(w, map[string]interface{}{"laborRateTiers": result})
}

func (r *Resolver) handleLaborRateTier(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	vars := variables(gqlReq)
	id := varID(vars, "id")

	tier, err := r.repo.GetLaborRateTier(req.Context(), id)
	if err != nil {
		writeGQLError(w, 404, "labor rate tier not found")
		return
	}
	writeGQL(w, map[string]interface{}{"laborRateTier": toMap(tier)})
}

func (r *Resolver) handleCountries(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	countries, err := r.repo.GetCountries(req.Context())
	if err != nil {
		writeGQLError(w, 500, "failed to fetch countries")
		return
	}
	result := make([]map[string]interface{}, len(countries))
	for i, c := range countries {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"countries": result})
}

func (r *Resolver) handleCurrencies(w http.ResponseWriter, req *http.Request, gqlReq GraphQLRequest) {
	currencies, err := r.repo.GetCurrencies(req.Context())
	if err != nil {
		writeGQLError(w, 500, "failed to fetch currencies")
		return
	}
	result := make([]map[string]interface{}, len(currencies))
	for i, c := range currencies {
		result[i] = toMap(c)
	}
	writeGQL(w, map[string]interface{}{"currencies": result})
}
