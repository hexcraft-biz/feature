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

const (
	ActionAssign int = iota
	ActionGrant
	ActionRevoke
)

const (
	WriteBehaviorCreate     = "CREATE"
	WriteBehaviorIdempotent = "IDEMPOTENT"
	WriteBehaviorOverwrite  = "OVERWRITE"
)
