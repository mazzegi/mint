package mint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var idfChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(idfChars, r) {
			return false
		}
	}
	return true
}

func extractFuncExpr(s string, expr string) (fe funcExpr, ok bool) {
	if !strings.HasSuffix(s, ")") {
		ok = false
		return
	}
	bo := strings.Index(s, "(")
	if bo < 0 {
		ok = false
		return
	}
	name := s[:bo]
	if !isValidIdentifier(name) {
		ok = false
		return
	}
	fe.name = name
	fe.args = strings.Split(s[bo+1:len(s)-1], ",")
	fe.expr = expr
	ok = true
	return
}

type Result string

type State struct {
	vars map[string]float64
	fncs *funcs
}

func NewState() *State {
	return &State{
		vars: map[string]float64{},
		fncs: newFuncs(),
	}
}

func (s *State) Eval(expr string) (Result, error) {
	fs := strings.Split(expr, "=")
	switch len(fs) {
	case 1:
		return s.evalExpr(fs[0])
	case 2:
		return s.evalAssignment(fs[0], fs[1])
	default:
		return "", errors.Errorf(("an expression may not contain more than 1 assignment"))
	}
}

func (s *State) evalAssignment(id string, expr string) (Result, error) {
	if fe, ok := extractFuncExpr(id, expr); ok {
		err := s.fncs.addFuncExpr(fe, func(id string) (float64, bool) {
			v, ok := s.vars[id]
			return v, ok
		})
		if err != nil {
			return "", err
		}
		return Result(fmt.Sprintf("function: %s(%s) => %s", fe.name, strings.Join(fe.args, ","), fe.expr)), nil
	}
	if !isValidIdentifier(id) {
		return "", errors.Errorf("assignment identifier is invalid")
	}
	v, err := s.eval(expr)
	if err != nil {
		return "", err
	}
	s.vars[id] = v
	return Result(fmt.Sprintf("stored: %s = %f", id, v)), nil
}

func (s *State) evalExpr(expr string) (Result, error) {
	v, err := s.eval(expr)
	if err != nil {
		return "", err
	}
	return Result(fmt.Sprintf("%f", v)), nil
}

func (s *State) eval(expr string) (float64, error) {
	e := newEvaluator(expr,
		func(id string) (float64, bool) {
			v, ok := s.vars[id]
			return v, ok
		},
		s.fncs,
	)
	return e.eval()
}
