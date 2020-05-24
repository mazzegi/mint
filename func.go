package mint

import (
	"math"

	"github.com/pkg/errors"
)

type Funcs interface {
	Eval(name string, vs ...float64) (float64, error)
	Contains(name string) bool
}

type funcExpr struct {
	name string
	args []string
	expr string
}

type funcs struct {
	items map[string]func(...float64) (float64, error)
}

func newFuncs() *funcs {
	fs := &funcs{
		items: map[string]func(...float64) (float64, error){},
	}
	fs.items["sqrt"] = func(vs ...float64) (float64, error) {
		if len(vs) != 1 {
			return 0, errors.Errorf("invalid amount or arguments")
		}
		return math.Sqrt(vs[0]), nil
	}
	fs.items["pow"] = func(vs ...float64) (float64, error) {
		if len(vs) != 2 {
			return 0, errors.Errorf("invalid amount or arguments")
		}
		return math.Pow(vs[0], vs[1]), nil
	}

	return fs
}

func (f *funcs) addFuncExpr(fe funcExpr, glookup func(id string) (float64, bool)) error {
	if _, contains := f.items[fe.name]; contains {
		return errors.Errorf("there's already a function %q", fe.name)
	}
	f.items[fe.name] = func(vs ...float64) (float64, error) {
		if len(vs) != len(fe.args) {
			return 0, errors.Errorf("invalid amount of arguments: want %d got %d", len(fe.args), len(vs))
		}
		m := map[string]float64{}
		for i, a := range fe.args {
			m[a] = vs[i]
		}
		return newEvaluator(fe.expr, func(id string) (float64, bool) {
			v, ok := m[id]
			if ok {
				return v, true
			}
			return glookup(id)
		}, f).eval()
	}
	return nil
}

func (f *funcs) Eval(name string, vs ...float64) (float64, error) {
	fnc, ok := f.items[name]
	if !ok {
		return 0, errors.Errorf("no such func %q", name)
	}
	return fnc(vs...)
}

func (f *funcs) Contains(name string) bool {
	_, ok := f.items[name]
	return ok
}
