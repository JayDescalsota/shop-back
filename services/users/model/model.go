package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TenantType string

const (
	TenantTypeAutoOwner  TenantType = "AUTO_OWNER"
	TenantTypeRepairShop TenantType = "REPAIR_SHOP"
	TenantTypePartsStore TenantType = "PARTS_STORE"
	TenantTypePlatform   TenantType = "PLATFORM"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
	TenantStatusTrial     TenantStatus = "TRIAL"
	TenantStatusCancelled TenantStatus = "CANCELLED"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenants"`

	ID        string       `bun:"id,pk,type:uuid" json:"id"`
	Name      string       `bun:"name,notnull" json:"name"`
	Type      TenantType   `bun:"type,notnull" json:"type"`
	Status    TenantStatus `bun:"status,notnull" json:"status"`
	Domain    *string      `bun:"domain" json:"domain"`
	Settings  string       `bun:"settings,type:jsonb" json:"settings"`
	CreatedAt time.Time    `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time    `bun:"updated_at,notnull" json:"updatedAt"`
}

type TenantSettings struct {
	bun.BaseModel `bun:"table:tenant_settings"`

	ID                 string    `bun:"id,pk,type:uuid" json:"id"`
	TenantID           string    `bun:"tenant_id,notnull,type:uuid" json:"tenantId"`
	BusinessHours      string    `bun:"business_hours,type:jsonb" json:"businessHours"`
	PaymentConfig      string    `bun:"payment_config,type:jsonb" json:"paymentConfig"`
	NotificationConfig string    `bun:"notification_config,type:jsonb" json:"notificationConfig"`
	Branding           string    `bun:"branding,type:jsonb" json:"branding"`
	Features           string    `bun:"features,type:jsonb" json:"features"`
	CreatedAt          time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt          time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles"`

	ID                string    `bun:"id,pk,type:uuid" json:"id"`
	UserID            string    `bun:"user_id,notnull" json:"userId"`
	Title             *string   `bun:"title" json:"title"`
	Department        *string   `bun:"department" json:"department"`
	Timezone          string    `bun:"timezone,notnull" json:"timezone"`
	Locale            string    `bun:"locale,notnull" json:"locale"`
	NotificationPrefs string    `bun:"notification_prefs,type:jsonb" json:"notificationPrefs"`
	CreatedAt         time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt         time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	Code        string    `bun:"code,notnull" json:"code"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description *string   `bun:"description" json:"description"`
	IsSystem    bool      `bun:"is_system,notnull" json:"isSystem"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
}

type Permission struct {
	bun.BaseModel `bun:"table:permissions"`

	ID          string    `bun:"id,pk,type:uuid" json:"id"`
	Code        string    `bun:"code,notnull" json:"code"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description *string   `bun:"description" json:"description"`
	Module      string    `bun:"module,notnull" json:"module"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
}

type TenantUser struct {
	bun.BaseModel `bun:"table:tenant_users"`

	ID        string    `bun:"id,pk,type:uuid" json:"id"`
	TenantID  string    `bun:"tenant_id,notnull,type:uuid" json:"tenantId"`
	UserID    string    `bun:"user_id,notnull" json:"userId"`
	RoleID    string    `bun:"role_id,notnull,type:uuid" json:"roleId"`
	InvitedAt time.Time `bun:"invited_at,notnull" json:"invitedAt"`
}

// BranchInfo is the user's view of a branch they can access.
type BranchInfo struct {
	TenantID    string     `json:"tenantId"`
	TenantName  string     `json:"tenantName"`
	TenantType  TenantType `json:"tenantType"`
	Role        string     `json:"role"`
	Permissions []string   `json:"permissions"`
}

// UserContext is the full SSO context returned by myContext / switchBranch.
type UserContext struct {
	UserID       string      `json:"userId"`
	Email        string      `json:"email"`
	Name         string      `json:"name"`
	ActiveBranch *BranchInfo `json:"activeBranch"`
	Branches     []BranchInfo `json:"branches"`
}
