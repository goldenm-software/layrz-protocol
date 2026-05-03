package definitions

// CommandDefinition defines the command structure based on the Layrz Protocol v2 specification
type CommandDefinition struct {
	// Is the command id, this value is unique and should be used to
	// send the ACK packet PdPacket to the server
	CommandId int

	// Is the command name, this value is used to identify the command
	CommandName *string

	// Is the command arguments, may contain any value depending of
	// the command definition
	Args map[string]any
}
