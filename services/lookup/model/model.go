package model

import (
	"time"

	"github.com/uptrace/bun"
)

type 	VehicleMake struct {
	bun.BaseModel `bun:"table:lookup.vehicle_makes"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Slug        string    `bun:"slug,notnull" json:"slug"`
	LogoURL     *string   `bun:"logo_url" json:"logoUrl"`
	Country     *string   `bun:"country" json:"country"`
	FoundedYear *int      `bun:"founded_year" json:"foundedYear"`
	IsActive    bool      `bun:"is_active,notnull" json:"isActive"`
	SortOrder   int       `bun:"sort_order,notnull" json:"sortOrder"`
	Metadata    string    `bun:"metadata,type:jsonb" json:"metadata"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`

	Models []*VehicleModel `bun:"rel:has-many,join:id=make_id" json:"models,omitempty"`
}

type 	VehicleModel struct {
	bun.BaseModel `bun:"table:lookup.vehicle_models"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	MakeID      string    `bun:"make_id,notnull,type:uuid" json:"makeId"`
	Name        string    `bun:"name,notnull" json:"name"`
	Slug        string    `bun:"slug,notnull" json:"slug"`
	YearStart   *int      `bun:"year_start" json:"yearStart"`
	YearEnd     *int      `bun:"year_end" json:"yearEnd"`
	VehicleType *string   `bun:"vehicle_type" json:"vehicleType"`
	TrimLevels  string    `bun:"trim_levels,type:jsonb" json:"trimLevels"`
	IsActive    bool      `bun:"is_active,notnull" json:"isActive"`
	Metadata    string    `bun:"metadata,type:jsonb" json:"metadata"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`

	Make *VehicleMake `bun:"rel:has-one,join:make_id=id" json:"make,omitempty"`
}

type 	ServiceType struct {
	bun.BaseModel `bun:"table:lookup.service_types"`

	ID             string    `bun:"id,pk,type:uuid" json:"id"`
	TenantID       *string   `bun:"tenant_id,type:uuid" json:"tenantId"`
	Code           string    `bun:"code,notnull" json:"code"`
	Name           string    `bun:"name,notnull" json:"name"`
	Description    *string   `bun:"description" json:"description"`
	Category       string    `bun:"category,notnull" json:"category"`
	System         string    `bun:"system,notnull" json:"system"`
	EstimatedHours *float64  `bun:"estimated_hours" json:"estimatedHours"`
	IsGlobal       bool      `bun:"is_global,notnull" json:"isGlobal"`
	IsActive       bool      `bun:"is_active,notnull" json:"isActive"`
	Metadata       string    `bun:"metadata,type:jsonb" json:"metadata"`
	CreatedAt      time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt      time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type 	DiagnosticCode struct {
	bun.BaseModel `bun:"table:lookup.diagnostic_codes"`

	ID                 string    `bun:"id,pk,type:uuid" json:"id"`
	Code               string    `bun:"code,notnull" json:"code"`
	System             string    `bun:"system,notnull" json:"system"`
	Description        string    `bun:"description,notnull" json:"description"`
	Severity           *string   `bun:"severity" json:"severity"`
	PossibleCauses     string    `bun:"possible_causes,type:jsonb" json:"possibleCauses"`
	RecommendedActions string    `bun:"recommended_actions,type:jsonb" json:"recommendedActions"`
	RelatedCodes       string    `bun:"related_codes,array" json:"relatedCodes"`
	IsActive           bool      `bun:"is_active,notnull" json:"isActive"`
	CreatedAt          time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt          time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type 	PartCategory struct {
	bun.BaseModel `bun:"table:lookup.part_categories"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	Code        string    `bun:"code,notnull" json:"code"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description *string   `bun:"description" json:"description"`
	ParentID    *string   `bun:"parent_id,type:uuid" json:"parentId"`
	Icon        *string   `bun:"icon" json:"icon"`
	SortOrder   int       `bun:"sort_order,notnull" json:"sortOrder"`
	IsActive    bool      `bun:"is_active,notnull" json:"isActive"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`

	Children []*PartCategory `bun:"rel:has-many,join:id=parent_id" json:"children,omitempty"`
}

type 	PartCompatibility struct {
	bun.BaseModel `bun:"table:lookup.part_compatibility"`

	ID         string    `bun:"id,pk,type:uuid" json:"id"`
	PartSKU    string    `bun:"part_sku,notnull" json:"partSku"`
	PartName   string    `bun:"part_name,notnull" json:"partName"`
	MakeID     *string   `bun:"make_id,type:uuid" json:"makeId"`
	ModelID    *string   `bun:"model_id,type:uuid" json:"modelId"`
	YearStart  *int      `bun:"year_start" json:"yearStart"`
	YearEnd    *int      `bun:"year_end" json:"yearEnd"`
	EngineType *string   `bun:"engine_type" json:"engineType"`
	Position   *string   `bun:"position" json:"position"`
	Notes      *string   `bun:"notes" json:"notes"`
	IsOem      bool      `bun:"is_oem,notnull" json:"isOem"`
	CreatedAt  time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type LookupPart struct {
	bun.BaseModel `bun:"table:lookup.parts"`

	ID        string    `bun:"id,pk,type:uuid" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Category  *string   `bun:"category" json:"category"`
	IsActive  bool      `bun:"is_active,notnull" json:"isActive"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type StorageLocation struct {
	bun.BaseModel `bun:"table:lookup.storage_locations"`

	ID        string    `bun:"id,pk,type:uuid" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Code      string    `bun:"code,notnull" json:"code"`
	IsActive  bool      `bun:"is_active,notnull" json:"isActive"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type 	FuelType struct {
	bun.BaseModel `bun:"table:lookup.fuel_types"`

	ID          string `bun:"id,pk,type:uuid" json:"id"`
	Code        string `bun:"code,notnull" json:"code"`
	Name        string `bun:"name,notnull" json:"name"`
	Description string `bun:"description" json:"description"`
	IsActive    bool   `bun:"is_active,notnull" json:"isActive"`
	SortOrder   int    `bun:"sort_order,notnull" json:"sortOrder"`
}

type 	TransmissionType struct {
	bun.BaseModel `bun:"table:lookup.transmission_types"`

	ID          string `bun:"id,pk,type:uuid" json:"id"`
	Code        string `bun:"code,notnull" json:"code"`
	Name        string `bun:"name,notnull" json:"name"`
	Description string `bun:"description" json:"description"`
	IsActive    bool   `bun:"is_active,notnull" json:"isActive"`
	SortOrder   int    `bun:"sort_order,notnull" json:"sortOrder"`
}

type 	EngineType struct {
	bun.BaseModel `bun:"table:lookup.engine_types"`

	ID          string  `bun:"id,pk,type:uuid" json:"id"`
	Code        string  `bun:"code,notnull" json:"code"`
	Name        string  `bun:"name,notnull" json:"name"`
	FuelType    *string `bun:"fuel_type" json:"fuelType"`
	Description string  `bun:"description" json:"description"`
	IsActive    bool    `bun:"is_active,notnull" json:"isActive"`
	SortOrder   int     `bun:"sort_order,notnull" json:"sortOrder"`
}

type 	LaborRateTier struct {
	bun.BaseModel `bun:"table:lookup.labor_rate_tiers"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	TenantID    *string   `bun:"tenant_id,type:uuid" json:"tenantId"`
	Name        string    `bun:"name,notnull" json:"name"`
	HourlyRate  float64   `bun:"hourly_rate,notnull" json:"hourlyRate"`
	Description *string   `bun:"description" json:"description"`
	IsGlobal    bool      `bun:"is_global,notnull" json:"isGlobal"`
	IsActive    bool      `bun:"is_active,notnull" json:"isActive"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type 	Country struct {
	bun.BaseModel `bun:"table:lookup.countries"`

	Code      string `bun:"code,pk" json:"code"`
	Name      string `bun:"name,notnull" json:"name"`
	PhoneCode string `bun:"phone_code,notnull" json:"phoneCode"`
	Currency  string `bun:"currency,notnull" json:"currency"`
	Timezone  string `bun:"timezone,notnull" json:"timezone"`
	IsActive  bool   `bun:"is_active,notnull" json:"isActive"`
	SortOrder int    `bun:"sort_order,notnull" json:"sortOrder"`
}

type 	Currency struct {
	bun.BaseModel `bun:"table:lookup.currencies"`

	Code          string `bun:"code,pk" json:"code"`
	Name          string `bun:"name,notnull" json:"name"`
	Symbol        string `bun:"symbol,notnull" json:"symbol"`
	DecimalPlaces int    `bun:"decimal_places,notnull" json:"decimalPlaces"`
	IsActive      bool   `bun:"is_active,notnull" json:"isActive"`
	SortOrder     int    `bun:"sort_order,notnull" json:"sortOrder"`
}
