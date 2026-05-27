# AutoEcosystem - Enterprise Backend Architecture

## 1. System Overview

An enterprise-grade, multi-tenant automotive ecosystem platform built with Go microservices, GraphQL Federation, and PostgreSQL with unified database tenant scoping.

### User Types & Capabilities
| Role | Capabilities |
|------|-------------|
| **Auto Owner** | Track vehicle expenses/repairs/costs, find & book repair shops, vehicle history |
| **Repair Shop Owner** | Customer management, accept bookings, staff/payroll, inventory/parts, payments |
| **Auto Parts Store Owner** | Parts requests/quotes, accept orders, inventory, HR, POS, marketplace |

### Architecture Style
- **Pattern**: Microservices with GraphQL Federation (Apollo v2.3)
- **Database**: Unified PostgreSQL with tenant_id scoping (single DB, logical multi-tenancy)
- **API Gateway**: GraphQL Mesh Gateway (federated supergraph)
- **Communication**: Synchronous (GraphQL Federation) + Asynchronous (NATS event bus)
- **Auth**: JWT-based with role-based access control (RBAC) + tenant isolation

---

## 2. Service Boundaries

```
┌─────────────────────────────────────────────────────────────────────┐
│                        API GATEWAY (Supergraph)                     │
│              GraphQL Federation Router + Rate Limiting              │
└───────────────┬─────────────────────────────────────────────────────┘
                │
    ┌───────────┼───────────┬───────────┬───────────┬────────────────┐
    │           │           │           │           │                │
┌───▼───┐ ┌────▼────┐ ┌────▼────┐ ┌────▼────┐ ┌────▼────┐ ┌────────▼──┐
│ auth  │ │ users   │ │vehicles │ │ bookings│ │ repair  │ │ inventory │
│       │ │ &       │ │ &       │ │ &       │ │ shop    │ │ & parts   │
│       │ │ tenant  │ │ history │ │ payment │ │ mgmt    │ │ mgmt      │
└───┬───┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └─────┬────┘
    │         │           │           │           │            │
┌───▼───┐ ┌───▼────┐ ┌───▼────┐ ┌───▼────┐ ┌────▼────┐ ┌─────▼─────┐
│notify │ │ parts  │ │ payroll│ │ search │ │ analytics│ │ lookup    │
│       │ │ market │ │ & HR   │ │ & rec  │ │ & audit  │ │ (refdata) │
└───────┘ └────────┘ └────────┘ └────────┘ └─────────┘ └───────────┘
    │         │           │           │           │
    │         └───────┬───┘           │           │
    │                 │               │           │
    └─────────────────┴───────────────┴───────────┘
                         files/
                        media/
```

### Service Descriptions

| Service | Port | Subgraph | Responsibility |
|---------|------|----------|----------------|
| **gateway** | 4000 | Supergraph | Federation router, auth middleware, rate limiting, CORS |
| **auth** | 8081 | Subgraph | Registration, login, JWT, password reset, MFA |
| **users** | 8082 | Subgraph | Tenant lifecycle, user profiles, roles, permissions, org management |
| **vehicles** | 8083 | Subgraph | Vehicle CRUD, maintenance history, expense tracking, documents |
| **bookings** | 8084 | Subgraph | Appointment scheduling, booking lifecycle, calendar |
| **repair** | 8085 | Subgraph | Repair orders, job cards, diagnostics, shop operations |
| **inventory** | 8086 | Subgraph | Parts inventory, stock levels, reorder alerts, warehouse |
| **parts-marketplace** | 8087 | Subgraph | Parts catalog, quotes, orders, B2B marketplace |
| **payments** | 8088 | Subgraph | Payment processing, invoicing, refunds, payouts |
| **payroll** | 8089 | Subgraph | Staff management, shifts, payroll, attendance |
| **notifications** | 8090 | - | Email, SMS, push notifications (NATS consumer) |
| **lookup** | 8090 | Subgraph | Reference data: vehicle makes/models, service types, diagnostic codes, fuel/transmission types, parts categories, labor rates, geo data. Heavily cached, read-optimized |
| **search** | 8091 | Subgraph | Full-text search (PostgreSQL TSVector), recommendations |
| **analytics** | 8092 | Subgraph | Business intelligence, reports, dashboards |
| **files** | 8093 | Subgraph | File upload, image processing, S3 integration |

---

## 3. Database Architecture

### Multi-Tenancy Strategy: Unified DB with Tenant Scoping

All tables include `tenant_id` column. Every query is scoped to the authenticated user's tenant via middleware-injected context.

```sql
-- Core tenant table
tenants (id, name, type, status, config JSONB, created_at, updated_at)

-- Tenant types: 'auto_owner', 'repair_shop', 'parts_store', 'platform'
```

### Schema Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────┐
│                      POSTGRES (Unified DB)                          │
├──────────────────┬──────────────────┬───────────────────────────────┤
│   AUTH DOMAIN    │   USERS DOMAIN   │      VEHICLES DOMAIN          │
│  ──────────────  │  ──────────────  │  ─────────────────────────   │
│ users            │ tenants          │ vehicles                     │
│ sessions         │ tenant_users     │ vehicle_mileage              │
│ refresh_tokens   │ roles            │ vehicle_maintenance          │
│ otp_codes        │ permissions      │ vehicle_expenses             │
│ password_resets  │ tenant_settings  │ vehicle_documents            │
├──────────────────┼──────────────────┼───────────────────────────────┤
│ BOOKINGS DOMAIN  │ REPAIR DOMAIN    │   INVENTORY DOMAIN           │
│  ──────────────  │  ──────────────  │  ──────────────────────────  │
│ bookings         │ repair_orders    │ inventory_items              │
│ booking_slots    │ repair_line_items│ inventory_movements          │
│ shop_availability│ job_cards        │ inventory_categories         │
│                  │ diagnostics      │ suppliers                    │
│                  │                  │ purchase_orders              │
├──────────────────┼──────────────────┼───────────────────────────────┤
│ PAYMENTS DOMAIN  │ PARTS MARKET     │   PAYROLL DOMAIN             │
│  ──────────────  │  ──────────────  │  ──────────────────────────  │
│ payments         │ parts_catalog    │ employees                    │
│ invoices         │ parts_quotes     │ employee_payrolls            │
│ refunds          │ parts_orders     │ shifts                       │
│ payment_methods  │ quote_items      │ attendance                   │
│ subscriptions    │                  │ leave_requests               │
├──────────────────┼──────────────────┼───────────────────────────────┤
│ FILES DOMAIN     │ ANALYTICS DOMAIN │   LOOKUP DOMAIN              │
│  ──────────────  │  ──────────────  │  ──────────────────────────  │
│ file_uploads     │ dashboard_cache  │ vehicle_makes                │
│ file_presigned   │ report_snapshots │ vehicle_models               │
│                  │                  │ service_types                │
│                  │                  │ diagnostic_codes             │
│                  │                  │ part_categories              │
│                  │                  │ fuel_types                   │
│                  │                  │ transmission_types           │
│                  │                  │ engine_types                 │
│                  │                  │ labor_rate_tiers             │
│                  │                  │ countries                    │
│                  │                  │ currencies                   │
│                  │                  │ part_compatibility           │
│                  │                  │                               │
│ AUDIT DOMAIN     │                  │                               │
│  ──────────────  │                  │                               │
│ audit_logs       │                  │                               │
│ data_access_logs │                  │                               │
└──────────────────┴──────────────────┴───────────────────────────────┘
```

### Key Design Decisions
1. **Row-Level Security (RLS)** enabled on all tenant-scoped tables
2. **Tenant ID** injected from JWT claims into query context
3. **Background job** periodically validates tenant isolation
4. **Audit trail** on all write operations with tenant context

---

## 4. Shared Kernel (Packages)

```
packages/
├── config/         # Configuration loader with validation
├── database/       # PostgreSQL connection pool, migrations runner
├── redis/          # Redis client with connection management
├── natss/          # NATS JetStream client for event bus
├── middleware/
│   ├── auth.go     # JWT validation, user/tenant injection
│   ├── tenant.go   # Tenant isolation enforcement
│   ├── ratelimit.go # Rate limiting per tenant
│   ├── cors.go     # CORS configuration
│   └── logging.go  # Request logging with correlation ID
├── errors/         # Standardized error types and codes
├── logger/         # Structured logging (zerolog)
├── validator/      # Request validation utilities
├── response/       # Standardized API response format
├── events/         # NATS event definitions and publishers
│   ├── types.go    # Event envelope, types
│   ├── publish.go  # Event publisher
│   └── subscribe.go # Event subscriber
├── pagination/     # Cursor-based pagination helpers
├── encryption/     # AES-256 encryption for sensitive data
└── idgen/          # Snowflake ID generator
```

---

## 5. Security Architecture

### Authentication Flow
```
Client → Gateway → Auth Middleware → JWT Validation → Service
                ↓
         Redis (JWT blacklist)
                ↓
         Auth Service (token refresh)
```

### Authorization Model
- **RBAC** (Role-Based Access Control) with tenant-scoped roles
- **ABAC** (Attribute-Based Access Control) for fine-grained permissions
- Role hierarchy: `super_admin > tenant_admin > manager > staff > viewer`

### Data Protection
- PII encrypted at rest (AES-256-GCM)
- TLS 1.3 for all service-to-service communication
- Secrets managed via environment injection (Vault in production)

---

## 6. Event-Driven Architecture

### NATS JetStream Topics
```
auth.*           # User events (created, updated, deleted)
tenant.*         # Tenant lifecycle events
vehicle.*        # Vehicle events
booking.*        # Booking lifecycle (created, confirmed, completed, cancelled)
repair.*         # Repair order events
inventory.*      # Stock level changes, low-stock alerts
payment.*        # Payment processed, refunded, failed
parts.*          # Quote created, order placed
lookup.*         # Reference data updated, cache invalidated
notification.*   # Notification delivery events
audit.*          # Audit trail events
```

### Event Envelope
```go
type Event struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    TenantID    string                 `json:"tenant_id"`
    Timestamp   time.Time              `json:"timestamp"`
    Source      string                 `json:"source"`
    Payload     map[string]interface{} `json:"payload"`
    CorrelationID string               `json:"correlation_id"`
}
```

---

## 7. Deployment Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                         KUBERNETES CLUSTER                           │
│                                                                      │
│  ┌─────────────┐   ┌────────────────────────────────────────────┐   │
│  │   Ingress   │   │              Services                      │   │
│  │   (Traefik) │   │  ┌──────┐ ┌──────┐ ┌──────┐ ... ┌──────┐ │   │
│  │             │───│──│GW   │ │Auth  │ │Users │     │Files │ │   │
│  └─────────────┘   │  └──────┘ └──────┘ └──────┘     └──────┘ │   │
│                    └────────────────────────────────────────────┘   │
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  │
│  │   PostgreSQL     │  │   Redis Cluster   │  │   NATS Cluster   │  │
│  │   (Primary +     │  │   (Session,       │  │   (Event Bus +   │  │
│  │    Read Replica) │  │    Cache, Rate    │  │    JetStream)    │  │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘  │
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐                        │
│  │   MinIO/S3       │  │   Prometheus +    │                        │
│  │   (Object Store) │  │   Grafana (Obs)   │                        │
│  └──────────────────┘  └──────────────────┘                        │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 8. API Design Conventions

### GraphQL Federation
- Each service exposes a subgraph with its own schema
- Gateway composes supergraph from subgraph schemas
- Entities shared via `@key` directive

### REST Endpoints (Internal/Admin)
- `/healthz` - Health check
- `/readyz` - Readiness probe
- `/metrics` - Prometheus metrics
- `/admin/*` - Admin-only operations

### Standard Response Envelope
```json
{
  "data": {},
  "errors": [],
  "meta": {
    "requestId": "uuid",
    "tenantId": "uuid",
    "timestamp": "2026-05-27T00:00:00Z"
  }
}
```

---

## 9. Migration Strategy

### Phase 1: Foundation (Weeks 1-4)
- [x] Auth service with basic user registration
- [ ] Tenant management and multi-tenancy layer
- [ ] Shared packages (middleware, errors, events)
- [ ] Database schema migration framework
- [ ] API Gateway with federation

### Phase 2: Core Services + Reference Data (Weeks 5-10)
- [ ] Users & tenant management
- [ ] **Lookup service** (vehicle makes/models, service types, diagnostic codes, enums, geo)
- [ ] Vehicles & history tracking
- [ ] Bookings & scheduling
- [ ] Repair shop management
- [ ] Inventory management

### Phase 3: Marketplace & Payments (Weeks 11-16)
- [ ] Parts marketplace
- [ ] Payment processing
- [ ] Quotation system
- [ ] Order management

### Phase 4: Operations (Weeks 17-20)
- [ ] Payroll & HR
- [ ] Notifications
- [ ] Search & recommendations
- [ ] Analytics & reporting

### Phase 5: Hardening (Weeks 21-24)
- [ ] Security audit
- [ ] Load testing
- [ ] Monitoring & alerting
- [ ] Documentation
- [ ] CI/CD pipeline

---

## 10. Technology Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Language | Go 1.25 | Performance, concurrency, single binary |
| ORM | Bun | Type-safe, fast, PostgreSQL-optimized |
| GraphQL | gqlgen + Apollo Federation | Type-safe codegen, subgraph federation |
| Database | PostgreSQL 16 | ACID, JSONB, full-text search, RLS |
| Cache | Redis 7 | Sessions, rate limiting, pub/sub |
| Event Bus | NATS JetStream | Lightweight, persistent streams |
| Router | go-chi/chi | Lightweight, middleware-friendly |
| Logging | zerolog | Structured, fast, low overhead |
| Config | godotenv + viper | Environment-based config |
| Testing | testify + gqlgen testing | Comprehensive test coverage |
| CI/CD | GitHub Actions | Automated build, test, deploy |
| Container | Docker | Consistent environments |
| Orchestration | Kubernetes (prod) | Scaling, self-healing |
