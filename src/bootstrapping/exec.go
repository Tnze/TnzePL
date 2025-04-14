package main

import (
	"slices"
	"strconv"
)

type tnScope map[string]any

var (
	rootScope = make(tnScope)
	scopes    = []tnScope{rootScope}
)

func pushScope() { scopes = append(scopes, nil) }
func popScope()  { scopes = scopes[:len(scopes)-1] }

type (
	expr interface {
		eval() (v any, brk bool)
	}
	exprProg []expr
	// Branch
	exprIfCheckItem struct {
		cond   expr
		action expr
	}
	exprIf struct {
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
	exprLoop     struct{ body exprProg }
)

func (e exprIf) eval() (v any, brk bool) {
	for _, checkItem := range e.ifCheckList {
		if c, cbrk := checkItem.cond.eval(); cbrk {
			return c, brk
		} else if c.(bool) {
			pushScope()
			v, brk = checkItem.action.eval()
			popScope()
			return
		}
	}
	if e.elseBranch != nil {
		pushScope()
		v, brk = e.elseBranch.eval()
		popScope()
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

func (s statDefine) eval() (any, bool) {
	v, brk := s.expression.eval()
	if brk {
		return v, brk
	}
	scope := &scopes[len(scopes)-1]
	if *scope == nil {
		*scope = make(tnScope)
	}
	(*scope)[s.identifier] = v
	return nil, false
}

func (s statAssign) eval() (any, bool) {
	v, brk := s.expression.eval()
	if brk {
		return v, brk
	}
	for _, scope := range slices.Backward(scopes) {
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

func (e exprLoad) eval() (any, bool) {
	for _, scope := range slices.Backward(scopes) {
		if scope == nil {
			continue
		}
		if v, ok := scope[e.identifier]; ok {
			return v, false
		}
	}
	return nil, false
}

func (e exprLiteral) eval() (any, bool) {
	return e.value, false
}

func (e exprProg) eval() (v any, brk bool) {
	for _, e := range e {
		v, brk = e.eval()
		if brk {
			return
		}
	}
	return
}

func (e exprLoop) eval() (any, bool) {
	for {
		v, brk := e.body.eval()
		if brk && v != (exprContinue{}) {
			return v, false // eat the flag
		}
	}
}

func (e exprBreak) eval() (any, bool) {
	var v any
	if e.value != nil {
		v, _ = e.value.eval()
	}
	return v, true
}

func (e exprContinue) eval() (any, bool) {
	return e, true
}

func evalLiteral(sym tnSymType) expr {
	literal := string(sym.Value.(unEvaled))
	if literal == "true" {
		return exprLiteral{value: true}
	}
	if literal == "false" {
		return exprLiteral{value: false}
	}
	if v, err := strconv.Atoi(literal); err == nil {
		return exprLiteral{value: v}
	}

	// Do nothing
	panic("Unevaled literal")
}
