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
	ContextKeyReqID     ContextKey = "requestID"
	LogKey              ContextKey = "logger"
	HTTPHeaderRequestID            = "X-Request-ID"
)

type ContextKey string

type Role string

type Currency string
