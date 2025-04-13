%{
package main

// import "log"
%}

%token IF ELSE LET FOR
%token COMMENT
%token IDENTIFIER LITERAL
%token FN
%token RARROW
%token BREAK CONTINUE

%%

program     : /* empty */
            | program expr
            | program COMMENT
            | program assign
            | program break
            ;

expr        : if_expr
            | for_expr
            | func
            | func_call
            | IDENTIFIER { $$.Value = find($1) }
            | LITERAL { $$.Value = eval($1) }
            ;

type        : IDENTIFIER ;
type_anno   : ':' type
            ;

if_expr     : IF expr block
            | IF expr block ELSE block
            | IF expr block ELSE if_expr
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

block       : '{' program '}'
            ;

assign      : LET IDENTIFIER '=' expr ';' { bind($2, eval($4)) }
            | LET IDENTIFIER type_anno '=' expr ';' { bind($2, eval($5)) }
            | IDENTIFIER '=' expr ';' { bind($1, eval($3)) }
            ;

func        : FN '(' args ')' ret_anno block
            | FN '(' ')' ret_anno block
            ;

args        : args ',' arg
            | arg
            ;

arg         : IDENTIFIER type_anno
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
