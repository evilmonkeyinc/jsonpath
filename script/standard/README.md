# Standard Script Engine

The standard script engine is a basic implementation of the script.Engine interface that is used as the default script engine

## Supported Operations

|operator|name|supported types|
|-|-|-|
|`\|\|`|logical OR|boolean|
|`&&`|logical AND|boolean|
|`==`|equals|number\|string|
|`!=`|not equals|number\|string|
|`<=`|less than or equal to|number|
|`>=`|greater than or equal to|number|
|`<`|less than|number|
|`>`|greater than|number|
|`=~`|regex|string|
|`+`|plus/addition|number|
|`-`|minus/subtraction|number|
|`**`|power|number|
|`*`|multiplication|number|
|`/`|division|number|
|`%`|modulus|integer|

### Regex

The regex operator will perform a regex match check using the left side argument as the input and the right as the regex pattern.

The right side pattern should be passed as a string, between single or double quotes, to ensure that no characters are mistaken for other operators.

> the regex operation is handled by the standard [`regexp`](https://pkg.go.dev/regexp) golang library `Match` function.

## Special Parameters

The following symbols/tokens have special meaning when used in script expressions and will be replaced before the expression is evaluated. The symbols used within a string, between single or double quotes, will not be replaced.

|symbol|name|replacement|
|-|-|-|
|`$`|root|the root json data node|
|`@`|current|the current json data node|
|`nil`|nil|`nil`|
|`null`|null|`nil`|

Using the root or current symbol allows to embed a JSONPath selector within an expression and it is expected that any argument that includes these characters should be a valid selector.

The nil and null tokens can be used interchangeably to represent a nil value.

> remember that the @ character has different meaning in subscripts than it does in filters.
