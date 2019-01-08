package vminstance

import "fmt"

// STOPCODE is halt signal
type STOPCODE struct {
	error
}

type emptyStack struct {
	s Stack
}

func (e emptyStack) Error() string {
	return "stack is empty"
}

type commandNotFound struct {
	error
	Opcode opcode
}

func (e commandNotFound) Error() string {
	return fmt.Sprintf("command with opcode [%v] not found", e.Opcode)
}
