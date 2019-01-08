package vminstance

// StackData is datatype of stack values
type StackData cpuArch

// Stack is some stack for cpu
type Stack interface {
	Push(StackData)
	Pop() (StackData, error)
	Len() int
	Head() []StackData
}

// simpleStack is VM stack
type simpleStack struct {
	s []StackData
}

func newStack() Stack {
	return &simpleStack{
		s: make([]StackData, 0),
	}
}

// Push adds an element to stack
func (s *simpleStack) Push(v StackData) {
	s.s = append(s.s, v)
}

// Pop returns new stack with data item
func (s *simpleStack) Pop() (StackData, error) {
	l := len(s.s)
	if l < 1 {
		return 0, emptyStack{s}
	}
	v := s.s[l-1]
	s.s = s.s[:l-1]
	return v, nil
}

func (s simpleStack) Len() int {
	return len(s.s)
}

func (s simpleStack) Head() []StackData {
	var resLen int
	if s.Len() > stackHeadLen {
		resLen = stackHeadLen
	} else {
		resLen = s.Len()
	}
	res := make([]StackData, resLen)
	l := s.Len()
	for i := 0; i < resLen; i++ {
		res[i] = s.s[l-i-1]
	}
	return res
}
