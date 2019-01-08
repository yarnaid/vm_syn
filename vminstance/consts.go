package vminstance

type cpuArch uint16
type opcode uint8
type vmStatus uint

//go:generate stringer -type=opcode
const (
	haltOp opcode = 0
	jumpOp        = 6
	jtOp          = 7
	jfOp          = 8
	outOp         = 19
	noopOp        = 21
)

// MAX is module base for math
const MAX = 32768

//go:generate stringer -type=vmStatus
// VM states
const (
	NewStatus vmStatus = iota
	RunningStatus
	FinishedStatus
	ErrorStatus
)

const (
	stackHeadLen = 10
)
