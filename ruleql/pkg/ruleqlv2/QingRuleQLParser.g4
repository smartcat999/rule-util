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

parser grammar QingRuleQLParser;

options {
	tokenVocab=QingRuleQLLexer;
}


// 1. Tokens & KeyWord
parse
 : ( root ) EOF
 ;


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
    : field_elem ( COMMA field_elem ) *
    ;

field_elem
    : expr AS xpath_name                            # ExprField
    | xpath_name                                    # XpathField
    ;

// 2.2 Topic
topic
    : (DIVIDE)? topic_item (DIVIDE topic_item)*
    ;

topic_item
    : PLUS | ADDR | INDENTIFIER | TOPICITEM
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
    : NOT_EQUALS? expr
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
   | expr op=( '*' | '/' | '%' ) expr                     # Binary
   | expr op=( '+' | '-' ) expr                         # Binary
   | expr op=( EQ | GT | LT | GTE | LTE | NE )  expr   # Binary
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
    : dotnotation_item ('.' dotnotation_item)*
    ;

dotnotation_item
    : identifierWithQualifier | identifierWithTOPICITEM
    ;

identifierWithTOPICITEM
    :  TOPICITEM | INDENTIFIER | NUMBER
    ;

identifierWithQualifier : INDENTIFIER LBRAKET RBRAKET
                        | INDENTIFIER LBRAKET NUMBER RBRAKET
                        | INDENTIFIER LBRAKET ADDR RBRAKET
                        ;



