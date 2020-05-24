package mint

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var normalizeReplacer = strings.NewReplacer(" ", "", "\r", "", "\n", "", "\t", "")

func normalize(s string) string {
	return normalizeReplacer.Replace(s)
}

type evaluator struct {
	expr   string
	pos    int
	lookup func(id string) (float64, bool)
	fncs   Funcs
}

func newEvaluator(expr string, lookup func(id string) (float64, bool), fncs Funcs) *evaluator {
	return &evaluator{
		expr:   normalize(expr),
		pos:    0,
		lookup: lookup,
		fncs:   fncs,
	}
}

const (
	plus   byte = '+'
	minus  byte = '-'
	times  byte = '*'
	div    byte = '/'
	exp    byte = '^'
	bopen  byte = '('
	bclose byte = ')'
)

func (e *evaluator) eval() (float64, error) {
	var curr string
	var lastOp byte

	num := func(s string) (float64, error) {
		if s == "" {
			return 0, errors.Errorf("empty identifier")
		}
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v, nil
		} else {
			v, ok := e.lookup(s)
			if !ok {
				return 0, errors.Errorf("no such identifier %q", curr)
			}
			return v, nil
		}
	}
	var stack *Stack
	push := func(v float64) error {
		defer func() {
			curr = ""
			lastOp = 0
		}()
		if stack == nil {
			switch lastOp {
			case minus:
				v = -v
			case times, div, exp:
				return errors.Errorf("invalid operator at beginning of expr")
			}
			stack = NewStack(OpPlus, Value(v))
			return nil
		}
		if lastOp == 0 {
			return errors.Errorf("missing operator")
		}
		switch lastOp {
		case plus:
			stack.Push(OpPlus, v)
		case minus:
			stack.Push(OpPlus, -v)
		case times:
			stack.Push(OpTimes, v)
		case div:
			stack.Push(OpTimes, 1.0/v)
		case exp:
			stack.Push(OpExp, v)
		}
		return nil
	}

	pushCurr := func() error {
		if curr == "" {
			return nil
		}
		if v, err := num(curr); err == nil {
			return push(v)
		} else {
			return err
		}
	}

	//fmt.Printf("eval: %q\n", e.expr)
	for e.pos < len(e.expr) {
		r := e.expr[e.pos]
		switch r {
		case bopen:
			ic, ok := findClosingBraceIdx(e.expr[e.pos+1:])
			if !ok {
				return 0, errors.Errorf("no closing brace found for open brace at %d", e.pos)
			}
			bexpr := e.expr[e.pos+1 : e.pos+1+ic]
			var v float64
			if curr != "" {
				sl := splitArgs(bexpr)
				args := []float64{}
				for _, s := range sl {
					varg, err := newEvaluator(s, e.lookup, e.fncs).eval()
					if err != nil {
						return 0, err
					}
					args = append(args, varg)
				}
				var err error
				v, err = e.fncs.Eval(curr, args...)
				if err != nil {
					return 0, err
				}
			} else {
				var err error
				v, err = newEvaluator(bexpr, e.lookup, e.fncs).eval()
				if err != nil {
					return 0, err
				}
			}
			push(v)
			e.pos = e.pos + 1 + ic + 1
		case bclose:
			return 0, errors.Errorf("unexpected closing brace at %d", e.pos)
		case plus, minus, times, div, exp:
			err := pushCurr()
			if err != nil {
				return 0, err
			}
			lastOp = r
			e.pos++
		default:
			curr += string(r)
			e.pos++
		}
	}
	err := pushCurr()
	if err != nil {
		return 0, err
	}
	if stack == nil {
		return 0, errors.Errorf("no elements on the eval-stack")
	}

	v := stack.Eval()
	return v, nil
}

func splitArgs(expr string) []string {
	args := []string{}
	curr := ""
	nopen := 0
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case bopen:
			nopen++
		case bclose:
			nopen--
		case ',':
			if nopen == 0 {
				args = append(args, curr)
				curr = ""
				continue
			}
		}
		curr += string(expr[i])
	}
	args = append(args, curr)
	return args
}

func findClosingBraceIdx(s string) (int, bool) {
	nopen := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case bopen:
			nopen++
		case bclose:
			if nopen == 0 {
				return i, true
			}
			nopen--
		}
	}
	return -1, false
}
