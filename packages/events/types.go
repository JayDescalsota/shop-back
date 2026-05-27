package events

// Event types for the AutoEcosystem platform
const (
	// Auth events
	EventUserRegistered    = "auth.user.registered"
	EventUserLoggedIn      = "auth.user.logged_in"
	EventUserPasswordReset = "auth.user.password_reset"
	EventUserInvited       = "auth.user.invited"

	// Tenant events
	EventTenantCreated  = "tenant.created"
	EventTenantUpdated  = "tenant.updated"
	EventTenantSuspended = "tenant.suspended"

	// Vehicle events
	EventVehicleCreated  = "vehicle.created"
	EventVehicleUpdated  = "vehicle.updated"
	EventVehicleDeleted  = "vehicle.deleted"
	EventMaintenanceDue  = "vehicle.maintenance.due"

	// Booking events
	EventBookingCreated   = "booking.created"
	EventBookingConfirmed = "booking.confirmed"
	EventBookingStarted   = "booking.started"
	EventBookingCompleted = "booking.completed"
	EventBookingCancelled = "booking.cancelled"

	// Repair events
	EventRepairCreated   = "repair.created"
	EventRepairStarted   = "repair.started"
	EventRepairCompleted = "repair.completed"
	EventRepairInvoiced  = "repair.invoiced"

	// Inventory events
	EventInventoryLowStock  = "inventory.low_stock"
	EventInventoryReceived = "inventory.received"
	EventInventoryAdjusted = "inventory.adjusted"

	// Payment events
	EventPaymentSucceeded = "payment.succeeded"
	EventPaymentFailed    = "payment.failed"
	EventPaymentRefunded  = "payment.refunded"
	EventInvoiceCreated   = "invoice.created"
	EventInvoicePaid      = "invoice.paid"

	// Parts marketplace events
	EventPartsQuoteRequested = "parts.quote.requested"
	EventPartsQuoteSubmitted = "parts.quote.submitted"
	EventPartsOrderPlaced    = "parts.order.placed"
	EventPartsOrderShipped   = "parts.order.shipped"

	// Notification events
	EventNotificationSend = "notification.send"

	// Payroll events
	EventPayrollProcessed = "payroll.processed"
	EventLeaveRequested   = "leave.requested"
	EventLeaveApproved    = "leave.approved"

	// Lookup events
	EventLookupMakeCreated    = "lookup.make.created"
	EventLookupMakeUpdated   = "lookup.make.updated"
	EventLookupModelCreated  = "lookup.model.created"
	EventLookupModelUpdated  = "lookup.model.updated"
	EventLookupServiceTypeCreated = "lookup.service_type.created"
	EventLookupServiceTypeUpdated = "lookup.service_type.updated"
	EventLookupDiagCodeCreated    = "lookup.diag_code.created"
	EventLookupDiagCodeUpdated    = "lookup.diag_code.updated"
	EventLookupCacheInvalidated   = "lookup.cache.invalidated"

	// Audit events
	EventAuditLogCreated = "audit.log.created"
)
