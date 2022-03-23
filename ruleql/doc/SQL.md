## SQL


### Simple Case
An example SQL statement looks like this:

```
SELECT color AS rgb FROM 'a/b' WHERE temperature > 50
```

An example MQTT message (also called an incoming payload) looks like this:
```
{
	"color":"red", 
	"temperature":50
}
```

If this message is published o n the 'a/b' topic, the rule 
is triggered and the SQL statement is evaluated. The SQL 
statement extracts the value of the color property if the 
"temperature" property is greater than 50. The WHERE clause 
specifies the condition temperature > 50. 
The AS keyword renames the "color" property to "rgb". 
The result (also called an outgoing payload) looks like this:
```
{
    "rgb":"red"
}
```

This data is then forwarded to the rule's action, which sends 
the data for more processing. 


### IoT Case
An example SQL statement looks like this:

```
SELECT
    *,
    temperature + '°F' AS temperature.fahrenheit,
    ((temperature - 32) * 5 / 9.0) + '°C' AS temperature.celsius,
    color as metadata.color,
    CASE "color"
        WHEN 'red' THEN '红色'
        WHEN 'green' THEN '绿色'
        ELSE 'none' 
        AS color_zh,
    ports[2] AS dev,
    metadata.name
FROM a/b
WHERE  color = 'red'

```

An example MQTT message (also called an incoming payload) looks like this:
```
{
	"color":"red", 
	"temperature":50, 
	"ports": ["lo0", "eth1", "eth2"],
	"metadata": {"name": "Light1", "price": 11.05}
}
```

If this message is published o n the 'a/b' topic, the rule 
is triggered and the SQL statement is evaluated. The SQL 
statement extracts the value of the color property if the 
"color" property is equate 'red''. The WHERE clause 
specifies the condition color = 'red'. 

The `*` Operator containing the message data as output. 

The `temperature + '°F'` , Converts the temperature to a string and concatenates the "°F" to the end of the temperature as result,
` AS temperature.fahrenheit ` will set the result as the 'fahrenheit' attribute of temperature JSON object.

The `((temperature - 32) * 5 / 9.0) + '°C'` convert fahrenheit to celsius and append the "°C" to the value as result,
` AS temperature.celsius ` will set the result as the 'celsius' attribute of temperature JSON object.

The Case statements like a switch statement, or if/else statements.

```
 {
	"color":"红色", "temperature":{"fahrenheit":"50°F","celsius":"10.000000°C"},
	"metadata": {"name": "Light1", "price": 11.05,"color":"red"}
 }
```

This data is then forwarded to the rule's action, which sends 
the data for more processing. 



