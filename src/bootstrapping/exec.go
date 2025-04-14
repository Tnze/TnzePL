package main

import (
	"strconv"
)

var identifiers = make(map[string]any)

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
	statAssign struct {
		identifier string
		expression expr
	}
	exprLoad struct {
		identifier string
	}
	exprLiteral struct {
		value any
	}
	exprProg []expr
	exprUnit struct{}
)

var unit = exprUnit{}

func (e exprIf) eval() any {
	for _, checkItem := range e.ifCheckList {
		if checkItem.cond.eval().(bool) {
			return checkItem.action.eval()
		}
	}
	if e.elseBranch != nil {
		return e.elseBranch
	}
	return nil
}

func (s statAssign) eval() any {
	identifiers[s.identifier] = s.expression.eval()
	return nil
}

func (e exprLoad) eval() any {
	return identifiers[e.identifier]
}

func (e exprLiteral) eval() any {
	return e.value
}

func (e exprProg) eval() (v any) {
	for _, e := range e {
		v = e.eval()
	}
	return v
}

func (e exprUnit) eval() any {
	return nil
}

func evalLiteral(sym tnSymType) expr {
	bytesLiteral, ok := sym.Value.(unEvaled)
	if !ok {
		return sym.Value.(expr) // Already evaluated
	}

	literal := string(bytesLiteral)
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
