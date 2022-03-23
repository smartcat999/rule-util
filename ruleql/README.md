# MDMP ruleQL


## Specification

### SQL REFERENCE

In MDMP ruleQL , rules are defined using an SQL-like syntax. 

[SQL Sample](doc/SQL.md)

#### General SQL

In MDMP ruleQL , SQL statements are composed of three types of clauses:


* SELECT

    Required. Extracts information from the incoming message payload and performs transformations.

* FROM

    The MQTT message topic filter. The rule is triggered for each message sent to an MQTT topic 
    that matches the filter specified here. Required for rules that are triggered by messages 
    that pass through the message broker. Optional for rules that are only triggered using 
    the Basic Ingest feature.

* WHERE

    Optional. Adds conditional logic that determines if the actions specified by a rule are carried out.

#### Data Types

The MDMP ruleQL supports all JSON data types.
    
Supported Data Types
    
|  Type   |                        Meaning                          |
|---------|---------------------------------------------------------|	
|Int      |A discrete Int. 34 digits maximum.                       |
|Decimal  |A Decimal with a precision of 34 digits                  |
|Boolean  |True or False.                                           |
|String	  |A UTF-8 string.                                          |
|Array	  |A series of values that don't have to have the same type.|
|Object	  |A JSON value consisting of a key and a value. Keys must be strings. Values can be any type.|
|Null	  |Null as defined by JSON.                                 |
|Undefined|Not a value.                                             |



#### Operators
 
 * AND 运算符 
 * OR 运算符
 * NOT 运算符
 * \> 运算符
 * \>= 运算符
 * \< 运算符
 * \<= 运算符
 * \<> 运算符
 * \= 运算符
 * \+ 运算符
 * \- 运算符
 * \* 运算符
 * \/ 运算符
 * \% 运算符

#### Functions

#### JSON

The MDMP ruleQL use the following extensions to ANSI SQL syntax to 
make it easier to work with nested JSON objects.

 * "." Operator
 
    This operator accesses members in embedded JSON objects and functions 
    identically to ANSI SQL and JavaScript. 

 * "*" Operator
    
    This functions in the same way as the * wildcard in ANSI SQL. It's 
    used in the SELECT clause only and creates a new JSON object containing 
    the message data. 
     

## Development

Generate ruleQL parser `make generate`:

```bash 
$ make generate
```


Run the files by running `make run`:

```bash 
$ make run
```


Test the files by running `make test`:

```bash 
$ make test
```

