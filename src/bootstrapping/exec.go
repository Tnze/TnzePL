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
		eval() any
	}
	exprIfCheckItem struct {
		cond   expr
		action expr
	}
	exprIf struct {
		ifCheckList []exprIfCheckItem
		elseBranch  expr
	}
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
	exprProg    []expr
)

func (e exprIf) eval() (v any) {
	for _, checkItem := range e.ifCheckList {
		if checkItem.cond.eval().(bool) {
			pushScope()
			v = checkItem.action.eval()
			popScope()
			return
		}
	}
	if e.elseBranch != nil {
		pushScope()
		v = e.elseBranch.eval()
		popScope()
	}
	return
}

func (s statDefine) eval() any {
	v := s.expression.eval()
	scope := &scopes[len(scopes)-1]
	if *scope == nil {
		*scope = make(tnScope)
	}
	(*scope)[s.identifier] = v
	return nil
}

func (s statAssign) eval() any {
	v := s.expression.eval()
	for _, scope := range slices.Backward(scopes) {
		if scope == nil {
			continue
		}
		if _, ok := scope[s.identifier]; ok {
			scope[s.identifier] = v
			return nil
		}
	}
	panic("Assign to undefined identifier: " + s.identifier)
}

func (e exprLoad) eval() any {
	for _, scope := range slices.Backward(scopes) {
		if scope == nil {
			continue
		}
		if v, ok := scope[e.identifier]; ok {
			return v
		}
	}
	return nil
}

func (e exprLiteral) eval() any {
	return e.value
}

func (e exprProg) eval() (v any) {
	for _, e := range e {
		v = e.eval()
	}
	return
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
