package main

import (
	"log"
	"slices"
	"strconv"
	"strings"
)

type tnScope map[string]any

type ExecCtx struct {
	rootScope tnScope
	scopes    []tnScope
}

func NewExecCtx() *ExecCtx {
	rootScope := make(tnScope)
	return &ExecCtx{
		rootScope: rootScope,
		scopes:    []tnScope{rootScope},
	}
}

func (ctx *ExecCtx) pushScope()            { ctx.scopes = append(ctx.scopes, nil) }
func (ctx *ExecCtx) popScope()             { ctx.scopes = ctx.scopes[:len(ctx.scopes)-1] }
func (ctx *ExecCtx) outterScope() *tnScope { return &ctx.scopes[len(ctx.scopes)-1] }
func (ctx *ExecCtx) clone() *ExecCtx {
	return &ExecCtx{
		rootScope: ctx.rootScope,
		scopes:    slices.Clone(ctx.scopes),
	}
}

type (
	expr interface {
		eval(ctx *ExecCtx) (v any, brk bool)
	}
	exprProg  []expr
	exprEmpty struct{}
	// Branch
	exprIfCheckItem struct{ cond, action expr }
	exprIf          struct {
		ifCheckList []exprIfCheckItem
		elseBranch  expr
	}
	// Assign
	statDefine struct {
		identifier string
		expression expr
	}
	statAssign struct {
		identifier string
		expression expr
	}
	exprLoad    struct{ identifier string }
	exprLiteral struct{ value any }
	// Loop
	exprBreak    struct{ value expr }
	exprContinue struct{}
	exprLoop     struct{ body expr }
	exprWhile    struct{ cond, body, elseBranch expr }
	// Operations
	exprFuncCall struct {
		fn   expr
		args []expr
	}
	exprFuncDefine struct {
		args   []funcArgAnno
		retTyp string
		body   expr
	}
	funcArgAnno struct{ arg, typ string }
)

func (e exprEmpty) eval(ctx *ExecCtx) (any, bool) {
	return nil, false
}

func (e exprIf) eval(ctx *ExecCtx) (v any, brk bool) {
	for _, checkItem := range e.ifCheckList {
		if c, cbrk := checkItem.cond.eval(ctx); cbrk {
			return c, brk
		} else if c.(bool) {
			ctx.pushScope()
			v, brk = checkItem.action.eval(ctx)
			ctx.popScope()
			return
		}
	}
	if e.elseBranch != nil {
		ctx.pushScope()
		v, brk = e.elseBranch.eval(ctx)
		ctx.popScope()
	}
	return
}

func (e exprIf) addIf(cond, action tnSymType) expr {
	item := exprIfCheckItem{
		cond:   cond.Value.(expr),
		action: action.Value.(expr),
	}
	e.ifCheckList = append(e.ifCheckList, item)
	return e
}

func (e exprIf) addElse(action tnSymType) expr {
	e.elseBranch = action.Value.(expr)
	return e
}

func (s statDefine) eval(ctx *ExecCtx) (any, bool) {
	v, brk := s.expression.eval(ctx)
	if brk {
		return v, brk
	}
	scope := ctx.outterScope()
	if *scope == nil {
		*scope = make(tnScope)
	}
	(*scope)[s.identifier] = v
	return nil, false
}

func (s statAssign) eval(ctx *ExecCtx) (any, bool) {
	v, brk := s.expression.eval(ctx)
	if brk {
		return v, brk
	}
	for _, scope := range slices.Backward(ctx.scopes) {
		if scope == nil {
			continue
		}
		if _, ok := scope[s.identifier]; ok {
			scope[s.identifier] = v
			return nil, false
		}
	}
	panic("Assign to undefined identifier: " + s.identifier)
}

func (e exprLoad) eval(ctx *ExecCtx) (any, bool) {
	for _, scope := range slices.Backward(ctx.scopes) {
		if scope == nil {
			continue
		}
		if v, ok := scope[e.identifier]; ok {
			return v, false
		}
	}
	return nil, false
}

func (e exprLiteral) eval(ctx *ExecCtx) (any, bool) {
	return e.value, false
}

func (e exprProg) eval(ctx *ExecCtx) (v any, brk bool) {
	for _, e := range e {
		v, brk = e.eval(ctx)
		if brk {
			return
		}
	}
	return
}

func (e exprBreak) eval(ctx *ExecCtx) (any, bool) {
	var v any
	if e.value != nil {
		v, _ = e.value.eval(ctx)
	}
	return v, true
}

func (e exprContinue) eval(ctx *ExecCtx) (any, bool) {
	return e, true
}

func (e exprLoop) eval(ctx *ExecCtx) (any, bool) {
	for {
		v, brk := e.body.eval(ctx)
		if brk && v != (exprContinue{}) {
			return v, false // eat the flag
		}
	}
}

func (e exprWhile) eval(ctx *ExecCtx) (any, bool) {
	for {
		c, brk := e.cond.eval(ctx)
		if brk && c != (exprContinue{}) {
			// brak inside the condition also stop the loop
			return c, false
		}
		if !c.(bool) {
			if e.elseBranch != nil {
				return e.elseBranch.eval(ctx)
			} else {
				return nil, false
			}
		}
		v, brk := e.body.eval(ctx)
		if brk && v != (exprContinue{}) {
			return v, false
		}
	}
}

func (e exprWhile) addElse(action tnSymType) expr {
	e.elseBranch = action.Value.(expr)
	return e
}

func (e exprFuncCall) eval(ctx *ExecCtx) (any, bool) {
	fn, brk := e.fn.eval(ctx)
	if brk {
		return fn, brk
	}
	args := make([]any, len(e.args))
	for i, e := range e.args {
		v, brk := e.eval(ctx)
		if brk {
			return v, brk
		}
		args[i] = v
	}

	ctx.pushScope()
	v := fn.(func(args []any) any)(args)
	ctx.popScope()
	return v, false
}

func (e exprFuncDefine) eval(ctx *ExecCtx) (any, bool) {
	ctx = ctx.clone()
	return func(args []any) any {
		// TODO: check argument types
		if paramsLen, argsLen := len(args), len(e.args); paramsLen != argsLen {
			log.Panicf("wrong number of args: want %d got %d", argsLen, paramsLen)
		}

		// Create scope
		ctx.pushScope()
		defer ctx.popScope()

		// Inject arguments
		scope := ctx.outterScope()
		*scope = make(tnScope)
		for i, arg := range args {
			(*scope)[e.args[i].arg] = arg
		}

		// Execute function body
		v, brk := e.body.eval(ctx)
		if brk {
			panic("not allowed break here")
		}
		return v
	}, false
}

func evalLiteral(sym tnSymType) expr {
	literal := string(sym.Lexeme)
	if literal == "true" {
		return exprLiteral{value: true}
	}
	if literal == "false" {
		return exprLiteral{value: false}
	}
	if v, err := strconv.Atoi(literal); err == nil {
		return exprLiteral{value: v}
	}
	if literal[0] == '"' {
		return exprLiteral{value: strings.ReplaceAll(literal[1:len(literal)-1], `\"`, `"`)}
	}

	// Do nothing
	panic("Unevaled literal")
}
