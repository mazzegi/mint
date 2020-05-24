package mint

import "testing"

func TestEval(t *testing.T) {
	s := NewState()
	r, err := s.evalExpr("sqrt(pow(8,2))")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestSplitArgs(t *testing.T) {
	expr := "1,2+4"
	sl := splitArgs(expr)
	t.Log(sl)

	expr = "12,"
	sl = splitArgs(expr)
	t.Log(sl)

	expr = "pow(a,b)"
	sl = splitArgs(expr)
	t.Log(sl)

	expr = "pow,(a,b)"
	sl = splitArgs(expr)
	t.Log(sl)
}

func TestFuncDecl(t *testing.T) {
	s := "my(a,b)"
	fe, ok := extractFuncExpr(s, "")
	t.Log(fe, ok)

	s = "21(a,b)"
	fe, ok = extractFuncExpr(s, "")
	t.Log(fe, ok)

	s = "(a,b)xy"
	fe, ok = extractFuncExpr(s, "")
	t.Log(fe, ok)

	s = "a(b)"
	fe, ok = extractFuncExpr(s, "")
	t.Log(fe, ok)
}
