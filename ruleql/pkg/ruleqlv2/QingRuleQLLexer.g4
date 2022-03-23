/*
 * Copyright (C) 2019 Yunify, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this work except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * 1. 支持select
 * 2. 支持where
 * 3. 支持简单表达式
 * 4. 支持action模板
 */
lexer grammar QingRuleQLLexer;

// 1. Tokens & KeyWord
// 1.1 KeyWord
AS:                     A S;
AND:                    A N D;
CASE:                   C A S E;
ELSE:                   E L S E;
END:                    E N D;
FROM:                   F R O M;
NULL:                   N U L L;
NOT:                    N O T;
SELECT:                 S E L E C T;
THEN:                   T H E N;
WHERE:                  W H E R E;
WHEN:                   W H E N;
GROUP:                  G R O U P;
BY:                     B Y;
TUMBLINGWINDOW:         T U M B L I N G W I N D O W;
HOPPINGWINDOW:          H O P P I N G W I N D O W;
SLIDINGWINDOW:          S L I D I N G W I N D O W;
SESSIONWINDOW:          S E S S I O N W I N D O W;
DD:                     D D;
HH:                     H H;
MI:                     M I;
SS:                     S S;
MS:                     M S;


// 1.2 Token

COLON        : ':'                    ;
COMMA        : ','                    ;
SEMI         : ';'                    ;
LPAREN       : '('                    ;
RPAREN       : ')'                    ;
RARROW       : '->'                   ;
LT           : '<'                    ;
GT           : '>'                    ;
QUESTION     : '?'                    ;
STAR         : '*'                    ;
PLUS         : '+'                    ;
ADDR         : '#'                    ;
CONCAT       : '||'                   ;
OR           : '|'                    ;
DOLLAR       : '$'                    ;
DOT		     : '.'                    ;
PERCENT      : '%'                    ;
MINUS        : '-'                    ;
DIVIDE       : '/'                    ;
EQUALS       : '=='                   ;
ASSIGN       : '='                    ;
NOT_EQUALS   : '!='                   ;
NOT_EQUALS2  : '<>'                   ;
LE           : '<='                   ;
GE           : '>='                   ;
LBRAKET      : '['                    ;
RBRAKET      : ']'                    ;
LCURLY       : '{'                    ;
RCURLY       : '}'                    ;
DMARKS       : '"'                    ;
SMARKS       : '\''                    ;


TRUE:               T R U E;
FALSE:              F A L S E;
INDENTIFIER:        [a-zA-Z_#][a-zA-Z_\-#$@0-9]*;
NUMBER:             '0' | [1-9][0-9]* ;
INTEGER:            ('+' | '-')? NUMBER;
FLOAT:              ('+' | '-')? (NUMBER+ '.' NUMBER+ |  NUMBER+ '.' | '.' NUMBER+);
TOPICITEM:          [a-zA-Z_\-#$@0-9]+;
STRING:             '\'' (~'\'' | '\'\'')* '\'';
WHITESPACE:         [ \r\n\t]+ -> skip;


IDENTIFIER
  : [a-zA-Z_] [a-zA-Z_0-9]*
  ;

NUMERIC_LITERAL
 : DIGIT+ ( '.' DIGIT* )? ( E [-+]? DIGIT+ )?
 | '.' DIGIT+ ( E [-+]? DIGIT+ )?
 ;

STRING_LITERAL
 : '\'' ( ~'\'' | '\\\'' )* '\''
 ;

QUOTED_LITERAL
 : '`' ( ~'`' )* '`'
 ;

SPACES
 : [ \u000B\t\r\n] -> channel(HIDDEN)
 ;

UNEXPECTED_CHAR
 : .
 ;

fragment DIGIT : [0-9];


fragment A: [aA];
fragment B: [bB];
fragment C: [cC];
fragment D: [dD];
fragment E: [eE];
fragment F: [fF];
fragment G: [gG];
fragment H: [hH];
fragment I: [iI];
fragment J: [jJ];
fragment K: [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment Q: [qQ];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];
fragment Y: [yY];
fragment Z: [zZ];