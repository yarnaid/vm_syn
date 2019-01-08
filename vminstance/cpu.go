package vminstance

type cpu interface {
	Registry() Registry
	Stack() Stack
}

type cpuData struct {
	reg             Registry
	stack           Stack
	commandsMapping commandMap
}

func newCPU() cpu {
	return &cpuData{
		reg:   Registry{},
		stack: newStack(),
	}
}

func (cpu cpuData) Registry() Registry {
	return cpu.reg
}

func (cpu cpuData) Stack() Stack {
	return cpu.stack
}
