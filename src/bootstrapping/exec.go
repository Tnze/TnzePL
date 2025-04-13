package main

import (
	"strconv"
)

func eval(expr tnSymType) any {
	bytesLiteral, ok := expr.Value.(unEvaled)
	if !ok {
		return expr.Value // Already evaluated
	}

	literal := string(bytesLiteral)
	if literal == "true" {
		return true
	}
	if literal == "false" {
		return false
	}
	if v, err := strconv.Atoi(literal); err == nil {
		return v
	}

	// Do nothing
	return expr.Value
}

var identifiers = make(map[string]any)

func bind(identifier tnSymType, value any) {
	key := string(identifier.Value.(unEvaled))
	identifiers[key] = value
}

func find(identifier tnSymType) any {
	return identifiers[string(identifier.Value.(unEvaled))]
}
