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
	HeaderViaEndpointId = "X-Via-Endpoint-Id"
	HeaderByCustodianId = "X-By-Custodian-Id"
)

const (
	ActionAssign int = iota
	ActionGrant
	ActionRevoke
)

const (
	WriteBehaviorIdempotent = "IDEMPOTENT"
	WriteBehaviorOverwrite  = "OVERWRITE"
)
