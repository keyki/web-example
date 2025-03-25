package types

var (
	ADMIN = Role("ADMIN")
	USER  = Role("USER")
	Roles = []Role{ADMIN, USER}

	HUF = Currency("HUF")
	USD = Currency("USD")
	EUR = Currency("EUR")
)

type Role string

type Currency string
