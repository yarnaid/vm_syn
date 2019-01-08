package vminstance

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/yarnaid/vm_syn/internal/logger"
)

// VM is implementations of VM
type VM interface {
	Terminal() Terminal
	SetTerminal(t Terminal)
	Run()
	LoadFile(filename string)
	Debug(...interface{})
	Registry() Registry
	CPU() cpu
	Status() vmStatus
	SetStatus(status vmStatus)
	Addr() cpuArch
	GetByAddr(addr cpuArch) cpuArch
	SetByAddr(addr, value cpuArch)
}

type vmData struct {
	cpu      cpu
	Memory   Memory
	terminal Terminal
	err      error
	Ind      cpuArch
	logger   logger.Logger
	status   vmStatus
	counter  int64
}

func (vm vmData) Registry() Registry {
	return vm.cpu.Registry()
}

// New initiate virtual machine
func New() VM {
	logger := logrus.New().WithFields(logrus.Fields{
		"app":    "root",
		"stream": "logger",
	})
	logger.Logger.SetLevel(logrus.DebugLevel)
	vm := &vmData{
		Memory:   newMemory(),
		cpu:      newCPU(),
		terminal: logger.WithField("stream", "terminal").Writer(),
		logger:   logger,
		status:   NewStatus,
	}
	vm.Debug("vm created")
	return vm
}

func (vm *vmData) Terminal() Terminal        { return vm.terminal }
func (vm *vmData) SetTerminal(t Terminal)    { vm.terminal = t }
func (vm *vmData) CPU() cpu                  { return vm.cpu }
func (vm *vmData) Status() vmStatus          { return vm.status }
func (vm *vmData) SetStatus(status vmStatus) { vm.status = status }
func (vm *vmData) Addr() cpuArch             { return vm.Ind }
func (vm *vmData) GetByAddr(addr cpuArch) cpuArch {
	entry := vm.logger.WithFields(logrus.Fields{
		// "addr": addr,
		// "MAX":  MAX,
		// "gt":   addr >= MAX,
	})
	if addr >= MAX {
		entry.Debug(addr)
		return vm.Registry().Get(addr % MAX)
	}
	entry.Debug(addr)
	return vm.Memory.Get(addr)
}
func (vm *vmData) SetByAddr(addr, value cpuArch) {
	if addr >= MAX {
		vm.Registry().Set(addr%MAX, value)
	}
	vm.Memory.Set(addr, value)
}

func (vm *vmData) log(s ...interface{}) {
	// TODO: simplify this
	for _, i := range s {
		x := i.([]interface{})
		for _, v := range x {
			vm.Terminal().Write([]byte(v.(string)))
		}
	}
	vm.Terminal().Write([]byte("\n"))
}

// Debug VM
func (vm *vmData) Debug(s ...interface{}) {
	// vm.log(s)
	vm.logger.Debugln(s)
}

// Debugf VM
func (vm *vmData) Debugf(x string, s ...interface{}) {
	// vm.log(s)
	vm.logger.Debugf(x, s)
}

// LoadFile read file into VM memory
func (vm *vmData) LoadFile(filename string) {
	if vm.err != nil {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		vm.err = errors.Wrap(err, "os loading file error")
		vm.logger.Error(vm.err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var addr cpuArch
	for {
		if vm.err != nil {
			break
		}
		var d cpuArch
		err = binary.Read(r, binary.LittleEndian, &d)
		if err == io.EOF {
			break
		}
		vm.SetByAddr(addr, d)
		if err != nil {
			vm.err = errors.Wrap(err, "command error")
			vm.logger.Error(vm.err)
		}
		addr++
	}
	vm.Debug("file loaded")
	if vm.err != nil {
		panic(vm.err)
	}
}

// Run starts code loading and execution
func (vm *vmData) Run() {
	vm.status = RunningStatus
	vm.Debug("started")

	ticker := time.NewTicker(time.Nanosecond * time.Duration(10000))
	for (vm.status != FinishedStatus) || (vm.err != nil) {
		cmd, err := vm.NewCommand()
		if err != nil {
			vm.err = errors.Wrap(err, "cmd create error")
			vm.logger.Error(vm.err)
			return
		}
		output, err := execLoggingMW(cmd.exec)(cmd)
		vm.counter++
		if err == (STOPCODE{}) {
			vm.Terminal().Write([]byte(fmt.Sprintf("TERMINATING: mem ind = %v, %v\n", vm.Ind, cmd)))
			vm.SetStatus(FinishedStatus)
			return
		}
		if err != nil {
			vm.err = errors.Wrap(err, "cmd exec error")
			vm.logger.Error(vm.err)
		}
		// if cmd.postExec != nil {
		// 	cmd.postExec(cmd)
		// }
		vm.terminal.Write(output)

		select {
		case <-ticker.C:
			// vm.Terminal().Write([]byte(fmt.Sprintf("mem ind = %v, %v\n", vm.Ind, cmd)))
		}
	}
	vm.status = FinishedStatus
	return
}

func (vm *vmData) NewCommand() (commandData, error) {
	code := opcode(vm.GetByAddr(vm.Ind))
	vm.Ind++
	cmd, ok := commandMapping[code]
	cmd.vm = vm
	cmd.Opcode = code
	if !ok {
		return commandMapping[noopOp], commandNotFound{Opcode: code}
	}
	if cmd.Operands != nil {
		for i := 0; i < len(cmd.Operands); i++ {
			cmd.Operands[i] = vm.GetByAddr(vm.Ind + cpuArch(i))
		}
		vm.Ind += cpuArch(len(cmd.Operands))
	}

	return cmd, nil
}

type commandData struct {
	exec     cmdExec
	postExec cmdExec
	Operands []cpuArch
	Opcode   opcode
	vm       *vmData
}

func (c commandData) String() string {
	return c.Opcode.String()
}

type commandMap map[opcode]commandData

var commandMapping commandMap

type cmdExec func(c commandData) ([]byte, error)

func execLoggingMW(f cmdExec) cmdExec {
	return cmdExec(func(c commandData) ([]byte, error) {
		entry := c.vm.logger.WithFields(logrus.Fields{
			"addr": c.vm.Ind,
			"cnt":  c.vm.counter,
			"ops":  c.Operands,
		})
		// if c.Opcode == outOp {
		// 	entry = entry.WithField("c", string(rune(c.Operands[0])))
		// }
		entry.Debug(c)
		return c.exec(c)
	})
}

func init() {
	commandMapping = commandMap{
		haltOp: commandData{
			exec: func(c commandData) ([]byte, error) {
				return nil, STOPCODE{}
			},
		},
		jumpOp: commandData{
			exec: func(c commandData) ([]byte, error) {
				c.vm.Ind = (c.Operands[0] % MAX)
				return nil, nil
			},
			postExec: func(c commandData) ([]byte, error) {
				// c.vm.Ind--
				return nil, nil
			},
			Operands: make([]cpuArch, 1),
		},
		jtOp: commandData{
			exec: func(c commandData) ([]byte, error) {
				if (c.Operands[0] % MAX) != 0 {
					c.vm.Ind = (c.Operands[1] % MAX)
				}
				return nil, nil
			},
			postExec: func(c commandData) ([]byte, error) {
				// c.vm.Ind--
				return nil, nil
			},
			Operands: make([]cpuArch, 2),
		},
		jfOp: commandData{
			exec: func(c commandData) ([]byte, error) {
				if (c.Operands[0] % MAX) == cpuArch(0) {
					c.vm.Ind = (c.Operands[1] % MAX)
				}
				return nil, nil
			},
			postExec: func(c commandData) ([]byte, error) {
				// c.vm.Ind--
				return nil, nil
			},
			Operands: make([]cpuArch, 2),
		},
		outOp: commandData{
			Operands: make([]cpuArch, 1),
			exec: func(c commandData) ([]byte, error) {
				return []byte{byte(c.Operands[0])}, nil
			},
		},
		noopOp: commandData{
			exec: func(c commandData) ([]byte, error) { return nil, nil },
		},
	}
}
