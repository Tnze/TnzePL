%{
package main

import "github.com/timtadh/lexmachine"
// import "log"

var tnRoot expr
%}

%token IF ELSE LET FOR
%token COMMENT
%token IDENTIFIER LITERAL
%token FN
%token RARROW
%token BREAK CONTINUE

%%

file        : program { tnRoot = $$.Value.(expr) } ;

program     : /* empty */ { $$.Token = &lexmachine.Token { Value: exprProg{} } }
            | program expr { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            | program COMMENT { $$.Value = $1.Value }
            | program assign { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            | program break { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            ;

expr        : if_expr
            | if_else_expr
            | for_expr
            | func
            | func_call
            | IDENTIFIER    { $$.Value = exprLoad { identifier: string($1.Value.(unEvaled)) } }
            | LITERAL       { $$.Value = evalLiteral($1) }
            ;

type        : IDENTIFIER ;
type_anno   : ':' type
            ;

if_expr     : IF expr block {
                checkItem := exprIfCheckItem{ cond: $2.Value.(expr), action: $3.Value.(expr) }
                $$.Value = exprIf {
                    ifCheckList: []exprIfCheckItem{checkItem},
                    elseBranch: nil,
                }
            }
            | if_expr ELSE IF expr block {
                ev := $1.Value.(exprIf)
                checkItem := exprIfCheckItem{ cond: $4.Value.(expr), action: $5.Value.(expr) }
                ev.ifCheckList = append(ev.ifCheckList, checkItem)
                $$.Value = ev
            }
            ;
if_else_expr: if_expr ELSE block {
                ev := $1.Value.(exprIf)
                ev.elseBranch = $3.Value.(expr)
                $$.Value = ev
            }
            ;

finite_for  : FOR expr block
            | FOR expr ';' expr ';' expr block
            | FOR LET IDENTIFIER ':' expr block
            ;

for_expr    : FOR block
            | finite_for
            | finite_for ELSE block
            ;

break       : BREAK ';'
            | BREAK expr ';'
            ;

block       : '{' program '}' { $$.Value = $2.Value }
            ;

assign      : LET IDENTIFIER '=' expr ';' {
                $$.Value = statDefine {
                    identifier: string($2.Value.(unEvaled)),
                    expression: $4.Value.(expr),
                }
            }
            | LET IDENTIFIER type_anno '=' expr ';' {
                $$.Value = statDefine {
                    identifier: string($2.Value.(unEvaled)),
                    expression: $4.Value.(expr),
                }
            }
            | IDENTIFIER '=' expr ';' {
                $$.Value = statAssign {
                    identifier: string($1.Value.(unEvaled)),
                    expression: $3.Value.(expr),
                }
            }
            ;

func        : FN '(' args ')' ret_anno block
            | FN '(' ')' ret_anno block
            ;

arg         : IDENTIFIER type_anno ;
args        : arg
            | args ',' arg
            ;


func_call   : IDENTIFIER '(' ')'
            | IDENTIFIER '(' params ')'
            ;

params      : params ',' expr
            | expr
            ;

ret_anno    : /* empty */
            | RARROW type
            ;
