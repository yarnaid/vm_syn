package vminstance

import "fmt"

const regLen = 8

// Registry is for VM registries
type Registry [8]cpuArch

// ToStringList converts values to strings
func (r Registry) ToStringList() []string {
	res := make([]string, regLen)
	for i := 0; i < regLen; i++ {
		res[i] = fmt.Sprint(r[i])
	}
	return res
}

func (r Registry) Get(addr cpuArch) cpuArch {
	return r[addr]
}

func (r Registry) Set(addr, value cpuArch) {
	r[addr] = value
}
