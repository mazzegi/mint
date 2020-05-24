package mint

import (
	"fmt"
	"math"
)

type Op int

const (
	OpPlus Op = iota
	OpTimes
	OpExp
)

func (op Op) String() string {
	switch op {
	case OpPlus:
		return "+"
	case OpTimes:
		return "*"
	default:
		return "^"
	}
}

type Evaler interface {
	Eval() float64
	Push(op Op, val float64) Evaler
	Dump() string
}

type Value float64

func (v Value) Eval() float64 {
	return float64(v)
}

func (v Value) Push(op Op, val float64) Evaler {
	s := NewStack(op, Value(v))
	s.Push(op, val)
	return s
}

func (v Value) Dump() string {
	return fmt.Sprintf("value %f", v)
}

type Stack struct {
	op   Op
	elts []Evaler
}

func NewStack(op Op, ev Evaler) *Stack {
	s := &Stack{
		op:   op,
		elts: []Evaler{ev},
	}
	return s
}

func (s *Stack) Eval() float64 {
	var v float64
	for i, e := range s.elts {
		if i == 0 {
			v = e.Eval()
			continue
		}
		switch s.op {
		case OpPlus:
			v += e.Eval()
		case OpTimes:
			v *= e.Eval()
		case OpExp:
			v = math.Pow(v, e.Eval())
		}
	}
	return v
}

func (s *Stack) Push(op Op, v float64) Evaler {
	if op == s.op {
		s.elts = append(s.elts, Value(v))
	} else if op > s.op {
		last := s.elts[len(s.elts)-1]
		last = last.Push(op, v)
		s.elts[len(s.elts)-1] = last
	} else if op < s.op {
		relts := s.elts
		rop := s.op
		s.op = op
		s.elts = []Evaler{
			&Stack{
				elts: relts,
				op:   rop,
			},
			Value(v),
		}
	}
	return s
}

func (s *Stack) Dump() string {
	d := fmt.Sprintf("--- stack (%s):\n", s.op)
	for _, e := range s.elts {
		d += e.Dump() + "\n"
	}
	d += fmt.Sprintf("--- stack \n")
	return d
}
