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

program     : /* empty */       { $$.Token = &lexmachine.Token { Value: exprProg{} } }
            | program expr      { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            | program COMMENT   { $$.Value = $1.Value }
            | program assign    { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            | program break     { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            | program continue  { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
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

if_expr     : IF expr block                 { $$.Value = exprIf{}.addIf($2, $3) }
            | if_expr ELSE IF expr block    { $$.Value = $1.Value.(exprIf).addIf($4, $5) }
            ;
if_else_expr: if_expr ELSE block { $$.Value = $1.Value.(exprIf).addElse($3) }
            ;

finite_for  : FOR expr block
            | FOR expr ';' expr ';' expr block
            | FOR LET IDENTIFIER ':' expr block
            ;

for_expr    : FOR block             { $$.Value = exprLoop { body: $2.Value.(exprProg) } }
            | finite_for
            | finite_for ELSE block
            ;

break       : BREAK ';'         { $$.Value = exprBreak {} }
            | BREAK expr ';'    { $$.Value = exprBreak { value: $2.Value.(expr) } }
            ;

continue    : CONTINUE ';' { $$.Value = exprContinue{} }
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
