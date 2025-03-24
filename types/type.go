package types

var (
    ADMIN = Role("ADMIN")
    USER  = Role("USER")
    Roles = []Role{ADMIN, USER}
)

type Role string
