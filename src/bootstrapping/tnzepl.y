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
%token RARROW // =>
%token BREAK CONTINUE // { }

%token EQ NE OROR ANDAND OR AND XOR // == != || && | ^ &
%token LE LEEQ BG BGEQ  // < <= > >=
%token LSHIFT RSHIFT // << >>

%%

file        : program       { tnRoot = append($1.Value.(exprProg), exprEmpty{}) }
            | program expr1 { tnRoot = append($1.Value.(exprProg), $2.Value.(expr)) }
            ;

program     : /* empty */   { $$.Token = &lexmachine.Token{ Value: exprProg{} } }
            | program stmt  { $$.Value = append($1.Value.(exprProg), $2.Value.(expr)) }
            ;

atom        : atom_nb
            | block_stmt
            | if_stmt
            | for_stmt
            ;

atom_nb     : '(' expr1 ')'     { $$.Value = $2.Value }
            | LITERAL           { $$.Value = evalLiteral($1) }
            | IDENTIFIER        { $$.Value = exprLoad{ identifier: string($1.Value.(unEvaled)) } }
            ;

stmt        : assign_stmt
            | break_stmt
            | continue_stmt
            | expr1 ';'
            ;

assign_stmt : LET IDENTIFIER '=' expr1 ';' {
                $$.Value = statDefine {
                    identifier: string($2.Value.(unEvaled)),
                    expression: $4.Value.(expr),
                }
            }
            // | LET IDENTIFIER type_anno '=' assign ';' {
            //     $$.Value = statDefine {
            //         identifier: string($2.Value.(unEvaled)),
            //         expression: $4.Value.(expr),
            //     }
            // }
            | IDENTIFIER '=' expr1 ';' {
                $$.Value = statAssign {
                    identifier: string($1.Value.(unEvaled)),
                    expression: $3.Value.(expr),
                }
            }
            ;
break_stmt      : BREAK expr1 ';' { $$.Value = exprBreak{ value: $2.Value.(expr) } } ;
continue_stmt   : CONTINUE ';' { $$.Value = exprContinue{} }

expr1       : expr1 OROR expr2
            | expr2
            ;
expr1_nb    : expr1_nb OROR expr2_nb
            | expr2_nb
            ;

expr2       : expr2 ANDAND expr3
            | expr3
            ;
expr2_nb    : expr2_nb ANDAND expr3_nb
            | expr3_nb
            ;

expr3       : expr3 OR expr4
            | expr4
            ;
expr3_nb    : expr3_nb OR expr4_nb
            | expr4_nb
            ;
expr4       : expr4 XOR expr5
            | expr5
            ;
expr4_nb    : expr4_nb XOR expr5_nb
            | expr5_nb
            ;

expr5       : expr5 AND expr6
            | expr6
            ;
expr5_nb    : expr5_nb AND expr6_nb
            | expr6_nb
            ;

expr6       : expr6 EQ expr7 { $$.Value = exprEq{ $1.Value.(expr), $3.Value.(expr) } }
            | expr6 NE expr7 { $$.Value = exprNe{ $1.Value.(expr), $3.Value.(expr) } }
            | expr7
            ;
expr6_nb    : expr6_nb EQ expr7_nb { $$.Value = exprEq{ $1.Value.(expr), $3.Value.(expr) } }
            | expr6_nb NE expr7_nb { $$.Value = exprNe{ $1.Value.(expr), $3.Value.(expr) } }
            | expr7_nb
            ;

expr7       : expr7 LE expr8
            | expr7 LEEQ expr8
            | expr7 BG expr8
            | expr7 BGEQ expr8
            | expr8
            ;
expr7_nb    : expr7_nb LE   expr8_nb
            | expr7_nb LEEQ expr8_nb
            | expr7_nb BG   expr8_nb
            | expr7_nb BGEQ expr8_nb
            | expr8_nb
            ;

expr8       : expr8 LSHIFT expr9
            | expr8 RSHIFT expr9
            | expr9
            ;
expr8_nb    : expr8_nb LSHIFT expr9_nb
            | expr8_nb RSHIFT expr9_nb
            | expr9_nb
            ;

expr9       : expr9 '+' expr10 { $$.Value = exprAdd{ $1.Value.(expr), $3.Value.(expr) } }
            | expr9 '-' expr10 { $$.Value = exprSub{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10
            ;
expr9_nb    : expr9_nb '+' expr10_nb { $$.Value = exprAdd{ $1.Value.(expr), $3.Value.(expr) } }
            | expr9_nb '-' expr10_nb { $$.Value = exprSub{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10_nb
            ;

expr10      : expr10 '*' expr11 { $$.Value = exprMul{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10 '/' expr11 { $$.Value = exprDiv{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10 '%' expr11 { $$.Value = exprMod{ $1.Value.(expr), $3.Value.(expr) } }
            | expr11
            ;
expr10_nb   : expr10_nb '*' expr11_nb { $$.Value = exprMul{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10_nb '/' expr11_nb { $$.Value = exprDiv{ $1.Value.(expr), $3.Value.(expr) } }
            | expr10_nb '%' expr11_nb { $$.Value = exprMod{ $1.Value.(expr), $3.Value.(expr) } }
            | expr11_nb
            ;

expr11      : atom '(' params_list ')' { $$.Value = exprFuncCall{ fn: $1.Value.(expr), args: $3.Value.([]expr) } }
            | atom
            ;
expr11_nb   : atom_nb '(' params_list ')' { $$.Value = exprFuncCall{ fn: $1.Value.(expr), args: $3.Value.([]expr) } }
            | atom_nb
            ;

params_list : /* empty */ { $$.Token = &lexmachine.Token{ Value: []expr{} } }
            | params
            ;
params      : params ',' expr1 { $$.Value = append($1.Value.([]expr), $3.Value.(expr)) }
            | expr1              { $$.Value = []expr{ $1.Value.(expr) } }
            ;

if_stmt     : if_only
            | if_else
            ;

if_else     : if_only ELSE block_stmt           { $$.Value = $1.Value.(exprIf).addElse($3) }
            ;

if_only     : IF expr1 block_stmt                { $$.Value = exprIf{}.addIf($2, $3) }
            | if_only ELSE IF expr1 block_stmt   { $$.Value = $1.Value.(exprIf).addIf($4, $5) }
            ;

for_stmt    : FOR block_stmt                    { $$.Value = exprLoop{ body: $2.Value.(expr) } }
            | finite_for                        { $$.Value = $1.Value }
            | finite_for ELSE block_stmt        { $$.Value = $1.Value.(exprWhile).addElse($3) }
            ;
finite_for  : FOR expr1_nb block_stmt            { $$.Value = exprWhile{ cond: $2.Value.(expr), body: $3.Value.(expr) } }
            | FOR expr1_nb ';' atom ';' expr1_nb block_stmt
            | FOR LET IDENTIFIER ':' atom_nb block_stmt
            ;

block_stmt  : '{' program '}'       { $$.Value = append($1.Value.(exprProg), exprEmpty{}) }
            | '{' program expr1 '}' { $$.Value = append($1.Value.(exprProg), $3.Value.(expr)) }
            ;
