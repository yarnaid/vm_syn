package vminstance

// Memory interface
type Memory interface {
	Get(addr cpuArch) cpuArch
	Set(addr, value cpuArch)
}

type memoryData [MAX]cpuArch

// simpleMemory is VM RAM
type simpleMemory struct {
	m memoryData
}

func newMemory() Memory {
	return &simpleMemory{}
}

func (m *simpleMemory) Get(addr cpuArch) cpuArch {
	return m.m[addr]
}

func (m *simpleMemory) Set(addr, value cpuArch) {
	m.m[addr] = value
}
