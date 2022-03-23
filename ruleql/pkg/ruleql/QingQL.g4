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

grammar QingQL;

// 1. Tokens & KeyWord
// 1.1 KeyWord
AS:                     A S;
AND:                    A N D;
CASE:                   C A S E;
ELSE:                   E L S E;
END:                    E N D;
EQ:                     E Q     | '=';
FROM:                   F R O M;
GT:                     G T     | '>';
GTE:                    G T E   | '>' '=';
LT:                     L T     | '<';
LTE:                    L T E   | '<' '=';
NE:                     N E     | '!' '=' | '<' '>';
NOT:                    N O T   | '!';
NULL:                   N U L L;
OR:                     O R;
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
MUL:                '*';
DIV:                '/';
MOD:                '%';
ADD:                '+';
SUB:                '-';
DOT:                '.';
TRUE:               T R U E;
FALSE:              F A L S E;
INDENTIFIER:        [a-zA-Z_#][a-zA-Z_\-#$@0-9]*;
NUMBER:             '0' | [1-9][0-9]* ;
INTEGER:            ('+' | '-')? NUMBER;
FLOAT:              ('+' | '-')? (NUMBER+ DOT NUMBER+ |  NUMBER+ DOT | DOT NUMBER+);
TOPICITEM:          [a-zA-Z_\-#$@0-9]+;
PATHITEM:           TOPICITEM (ARRAYITEM)? (DOT TOPICITEM (ARRAYITEM)?)*;
ARRAYITEM:          '[' NUMBER ']' | '[' '#' ']';
STRING:             '\'' (~'\'' | '\'\'')* '\'';
WHITESPACE:         [ \r\n\t]+ -> skip;




// 2. Rules
root
    : SELECT
      fields
      FROM topic
      where?
      group?
    ;

// 2.1 Select
fields
    : field_elem (',' field_elem)*
    ;



field_elem
    : expr AS xpath_name                            # ExprField
    | xpath_name                                    # XpathField
    ;


// 2.2 Topic
topic
    : (DIV)? topic_item (DIV topic_item)*
    ;

topic_item
    : '+' | '#' | TOPICITEM | INDENTIFIER | NUMBER
    ;

// 2.3 Where
where: WHERE filter;

filter
    : filter_condition
    ;

filter_condition
    : filter_condition_or (AND filter_condition_or)*
    ;

filter_condition_or
    : filter_condition_not (OR filter_condition_not)*
    ;

filter_condition_not
    : NOT? expr
    ;




// 2.4 Group
//`SELECT * FROM /topic/xxx GROUP BY TUMBLINGWINDOW(hh, 10)`,
//`SELECT * FROM /topic/xxx GROUP BY HOPPINGWINDOW(mi, 5, 1)`,
//`SELECT * FROM /topic/xxx GROUP BY SLIDINGWINDOW(ms, 5)`,
//`SELECT * FROM /topic/xxx GROUP BY SESSIONWINDOW(ss, 5, 1)`,
group
    : GROUP BY group_by_elements;

group_by_elements
    : group_by_element (',' group_by_element)*;

group_by_element
    : group_windows                                  # GroupWindows
    | xpath_name                                     # GroupField
    ;

group_window_unit
    : DD | HH | MI | SS | MS
    ;

group_window_length
    : NUMBER
    ;

group_window_interval
    : NUMBER
    ;

group_windows
    : group_tumbling_window
    | group_hopping_window
    | group_sliding_window
    | group_session_window
    | group_hopping_window
    ;

group_tumbling_window
    : TUMBLINGWINDOW '(' group_window_unit ','  group_window_interval ')'
    ;

group_hopping_window
    : HOPPINGWINDOW '(' group_window_unit ','  group_window_interval  ','  group_window_length ')'
    ;

group_sliding_window
    : SLIDINGWINDOW '(' group_window_unit ','  group_window_interval ')'
    ;

group_session_window
    : SESSIONWINDOW '(' group_window_unit ','  group_window_interval ')'
    ;






// 2.4 expr
expr
   : constant                                       # Braces
   | '(' expr ')'                                   # Braces
   | expr op=('*'|'/'|'%') expr                     # Binary
   | expr op=('+'|'-') expr                         # Binary
   | expr op=(EQ | GT | LT | GTE | LTE | NE) expr   # Binary
   | call_expr                                      # Function
   | switch_stmt                                    # Switch
   ;

constant
    : TRUE                                           # Boolean
    | FALSE                                          # Boolean
    | NUMBER                                         # Integer
    | INTEGER                                        # Integer
    | FLOAT                                          # Float
    | STRING                                         # String
    | xpath_name                                     # XPath
    ;


switch_stmt
    : CASE expr
    WHEN expr THEN expr (WHEN expr THEN expr)*
    (ELSE expr)?
    ;

call_expr
    : key=INDENTIFIER '(' (expr (',' expr)*)? ')';

/*
CASE v WHEN t[1] THEN r[1]
  WHEN t[2] THEN r[2] ...
  WHEN t[n] THEN r[n]
  ELSE r[e] END
*/

// 2.5 json
asterisk: '*';

xpath_name
        : asterisk
        | dotnotation
        | '"' + (dotnotation) + '"'
        ;


dotnotation
    : INDENTIFIER
    | PATHITEM
    ;

identifierWithTOPICITEM
    :  PATHITEM '[' ']'
    | PATHITEM '[' NUMBER ']'
    | PATHITEM '[' '#'  ']'
    | PATHITEM
    |  FLOAT
    ;

identifierWithQualifier : INDENTIFIER '[]'
                        | INDENTIFIER '[' NUMBER ']'
                        | INDENTIFIER '[#]'
                        | INDENTIFIER
                        ;



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