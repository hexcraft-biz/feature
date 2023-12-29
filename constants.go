package feature

const (
	ScopesDelimiter = " "
)

const (
	OwnershipOrganization = "ORGANIZATION"
	OwnershipPrivate      = "PRIVATE"
	OwnershipPublic       = "PUBLIC"
)

const (
	EnumOwnershipUndef int = iota
	EnumOwnershipOrganization
	EnumOwnershipPrivate
	EnumOwnershipPublic
)
