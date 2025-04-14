package main

import (
	"errors"
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
		log.Print("EOF")
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
		return scan.Token(id, unEvaled(match.Bytes), match), nil
	}
}

func singleCharToken(scan *lexmachine.Scanner, match *machines.Match) (any, error) {
	if len(match.Bytes) != 1 {
		log.Panic(match)
		return scan.Token(0, nil, match), errors.New("not a single charactor")
	}
	return scan.Token(int(match.Bytes[0]), unEvaled(match.Bytes), match), nil
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
	lexer.Add([]byte(`=>`), token(RARROW))
	lexer.Add([]byte("break"), token(BREAK))
	lexer.Add([]byte("continue"), token(CONTINUE))
	lexer.Add([]byte(`-?[0-9]+|true|false|"[^"]*"`), token(LITERAL))
	lexer.Add([]byte(`[a-zA-Z_][a-zA-Z0-9_]*`), token(IDENTIFIER))
	lexer.Add([]byte(`:|;|\{|\}|=|\(|\)|,`), singleCharToken)
	lexer.Add([]byte(`//[^\n]*\r?\n`), token(COMMENT))
	lexer.Add([]byte(` |\t|\r?\n`), ignoreToken)

	if err := lexer.Compile(); err != nil {
		log.Panic(err)
	}

	source, err := os.ReadFile("./testfiles/main.tn")
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
	tnDebug = 0
	tnParse(&tnLex{Scanner: scanner})

	log.Printf("%#v", tnRoot)

	tnRoot.eval()
	log.Printf("%v", identifiers)
}
