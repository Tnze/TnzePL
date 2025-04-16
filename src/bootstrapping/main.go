package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

//go:generate goyacc -o tnzepl.go -p "tn" tnzepl.y

type (
	tnLex     struct{ *lexmachine.Scanner }
	tnSymType struct {
		yys int // yy State
		*lexmachine.Token
	}
	unEvaled []byte
)

func (t *tnLex) Lex(lval *tnSymType) int {
	tok, err, eos := t.Scanner.Next()
	if eos {
		return 0 // EOF
	}
	if err != nil {
		log.Print(err)
		return -1
	}

	lval.Token = tok.(*lexmachine.Token)
	if tnDebug > 0 {
		log.Print(lval)
	}
	return lval.Type
}

func (t *tnLex) Error(s string) {
	log.Print(s)
}

func token(id int) lexmachine.Action {
	return func(scan *lexmachine.Scanner, match *machines.Match) (any, error) {
		return scan.Token(id, nil, match), nil
	}
}

func singleCharToken(scan *lexmachine.Scanner, match *machines.Match) (any, error) {
	if len(match.Bytes) != 1 {
		log.Panic(match)
		return scan.Token(0, nil, match), errors.New("not a single charactor")
	}
	return scan.Token(int(match.Bytes[0]), nil, match), nil
}

func ignoreToken(scan *lexmachine.Scanner, match *machines.Match) (any, error) {
	return nil, nil
}

func main() {
	lexer := lexmachine.NewLexer()
	lexer.Add([]byte("if"), token(IF))
	lexer.Add([]byte("else"), token(ELSE))
	lexer.Add([]byte("let"), token(LET))
	lexer.Add([]byte("for"), token(FOR))
	lexer.Add([]byte("fn"), token(FN))
	lexer.Add([]byte(`->`), token(RARROW))
	lexer.Add([]byte(`==`), token(EQ))
	lexer.Add([]byte(`!=`), token(NE))
	lexer.Add([]byte("break"), token(BREAK))
	lexer.Add([]byte("continue"), token(CONTINUE))
	lexer.Add([]byte(`-?[0-9]+|true|false|"[^"]*"`), token(LITERAL))
	lexer.Add([]byte(`[a-zA-Z_][a-zA-Z0-9_]*`), token(IDENTIFIER))
	lexer.Add([]byte(`:|;|\{|\}|=|\(|\)|,|\+|\-|\*|\/|%`), singleCharToken)
	lexer.Add([]byte(`//[^\n]*\r?\n`), ignoreToken)
	lexer.Add([]byte(` |\t|\r?\n`), ignoreToken)

	if err := lexer.Compile(); err != nil {
		log.Panic(err)
	}

	source, err := os.ReadFile("./testfiles/fn.tn")
	if err != nil {
		log.Panic(err)
	}

	scanner, err := lexer.Scanner(source)
	if err != nil {
		log.Panic(err)
	}

	// for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	log.Println(tok.(tnSymType).id)
	// }
	tnErrorVerbose = true
	tnDebug = 1
	errcode := tnParse(&tnLex{Scanner: scanner})
	log.Println("Parse result:", errcode)

	// log.Printf("%#v", tnRoot)

	ctx := NewExecCtx()
	ctx.rootScope["print"] = func(args []any) any {
		n, _ := fmt.Println(args...)
		return n
	}
	ctx.rootScope["assert"] = func(args []any) any {
		if len(args) != 1 {
			log.Panic("assert expect 1 boolean argument")
		}
		if v, ok := args[0].(bool); !ok {
			log.Panicf("assert expression type %T is not boolean", args[0])
		} else if !v {
			log.Panic("assert failed")
		}
		return nil
	}
	tnRoot.eval(ctx)
	log.Printf("Root Scope: %v", ctx.rootScope)
}
