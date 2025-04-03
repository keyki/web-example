package types

var (
	ADMIN = Role("ADMIN")
	USER  = Role("USER")
	Roles = []Role{ADMIN, USER}

	HUF = Currency("HUF")
	USD = Currency("USD")
	EUR = Currency("EUR")
)

const (
	ContextKeyReqID     ContextKey = "request-id"
	HTTPHeaderRequestID            = "X-Request-ID"
	AuditServerPort                = 8071
	WebServerPort                  = 8070
)

type ContextKey string

type Role string

type Currency string
