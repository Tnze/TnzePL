%{
package main

import "github.com/timtadh/lexmachine"
// import "log"

var tnRoot expr
%}

%token IF ELSE LET FOR
// %token COMMENT
%token IDENTIFIER LITERAL
%token FN
%token RARROW
%token BREAK CONTINUE

%%

file        : expr_stmt { tnRoot = $$.Value.(expr) } ;

// "inline" exprs and "block" exprs
atom_stmt   : /* empty */       { $$.Token = &lexmachine.Token{ Value: exprProg{} } }
            | atom_stmt atom    { $$.Value = append($$.Value.(exprProg), $2.Value.(expr)) }
            ;

atom        : atom_nb
            | '{' expr_stmt '}' { $$.Value = $2.Value }
            ;

atom_nb     : '(' atom_stmt ')' { $$.Value = $2.Value }
            | LITERAL           { $$.Value = evalLiteral($1) }
            | IDENTIFIER        { $$.Value = exprLoad{ identifier: string($1.Value.(unEvaled)) } }
            | if_stmt
            | for_stmt
            ;

expr_stmt   : /* empty */           { $$.Token = &lexmachine.Token{ Value: exprProg{ } } }
            | expr_list
            ;
expr_list   : expr                  { $$.Token = &lexmachine.Token{ Value: exprProg{ $1.Value.(expr) } } }
            | expr_list ';' expr    { $$.Value = append($1.Value.(exprProg), $3.Value.(expr)) }
            | expr_list ';'         { $$.Value = append($1.Value.(exprProg), exprEmpty{} ) }
            ;


expr        : LET IDENTIFIER '=' assign {
                $$.Value = statDefine {
                    identifier: string($2.Value.(unEvaled)),
                    expression: $4.Value.(expr),
                }
            }
            // | LET IDENTIFIER type_anno '=' assign {
            //     $$.Value = statDefine {
            //         identifier: string($2.Value.(unEvaled)),
            //         expression: $4.Value.(expr),
            //     }
            // }
            | IDENTIFIER '=' assign {
                $$.Value = statAssign {
                    identifier: string($1.Value.(unEvaled)),
                    expression: $3.Value.(expr),
                }
            }
            | BREAK assign { $$.Value = exprBreak{ value: $2.Value.(expr) } }
            | CONTINUE { $$.Value = exprContinue{} }
            | assign
            ;

assign      : assign '+' term
            | assign '-' term
            | term
            ;

term        : term '*' call
            | term '/' call
            | term '%' call
            | call
            ;

call        : atom '(' params_list ')' { $$.Value = exprFuncCall{ fn: $1.Value.(expr), args: $3.Value.([]expr) } }
            | atom
            ;

params_list : /* empty */ { $$.Token = &lexmachine.Token{ Value: []expr{} } }
            | params
            ;
params      : params ',' expr   { $$.Value = append($1.Value.([]expr), $3.Value.(expr)) }
            | expr              { $$.Value = []expr{ $1.Value.(expr) } }
            ;

if_stmt     : if_only
            | if_else
            ;

if_else     : if_only ELSE block_stmt           { $$.Value = $1.Value.(exprIf).addElse($3) }
            ;

if_only     : IF atom block_stmt                { $$.Value = exprIf{}.addIf($2, $3) }
            | if_only ELSE IF atom block_stmt   { $$.Value = $1.Value.(exprIf).addIf($4, $5) }
            ;

for_stmt    : FOR block_stmt                    { $$.Value = exprLoop{ body: $2.Value.(expr) } }
            | finite_for                        { $$.Value = $1.Value }
            | finite_for ELSE block_stmt        { $$.Value = $1.Value.(exprWhile).addElse($3) }
            ;
finite_for  : FOR atom_nb block_stmt            { $$.Value = exprWhile{ cond: $2.Value.(expr), body: $3.Value.(expr) } }
            | FOR atom_nb ';' atom ';' expr block_stmt
            | FOR LET IDENTIFIER ':' expr block_stmt
            ;

block_stmt  : '{' expr_stmt '}' { $$.Value = $2.Value }
            ;
