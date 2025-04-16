package main

type (
	exprAdd [2]expr
	exprSub [2]expr
	exprMul [2]expr
	exprDiv [2]expr
	exprMod [2]expr

	exprEq [2]expr
	exprNe [2]expr
)

func (e exprAdd) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs.(int) + lhs.(int), false
}

func (e exprSub) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs.(int) - lhs.(int), false
}

func (e exprMul) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs.(int) * lhs.(int), false
}

func (e exprDiv) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs.(int) / lhs.(int), false
}

func (e exprMod) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs.(int) % lhs.(int), false
}

func (e exprEq) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs == lhs, false
}

func (e exprNe) eval(ctx *ExecCtx) (any, bool) {
	rhs, brk := e[0].eval(ctx)
	if brk {
		return rhs, brk
	}
	lhs, brk := e[1].eval(ctx)
	if brk {
		return lhs, brk
	}
	return rhs != lhs, false
}
