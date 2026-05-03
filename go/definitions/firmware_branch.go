package definitions

// FirmwareBranch is the firmware branch type
type FirmwareBranch string

const (
	// Development firmware branch
	Development FirmwareBranch = "1"
	// Stable firmware branch
	Stable FirmwareBranch = "0"
)
