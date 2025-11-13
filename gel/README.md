 Pango GEL
=====================================================================

A Go language expression package.

### What is EL?
EL = expression language.  
pango/gel calculate this expression and return it's result.


### Simple usage
```go
gel.Calculate("3+4*5")  // Output: 23
```

### Variable support
```go
m := map[string]any{"a": 10}
gel.Calculate("a*10", m)  // Output: 100 
```

### Supported operator

 | Operator | Operator Number | Priority | Description |
 |----------|-----------------|----------|------------------|
 | ()       |  \*             | 100      | Parenthesis      |
 | ,        |  \*             | 90       | Comma between parameter |
 | .        |  2              | 1        | Property or method accessor |
 | {1,2}    |  \*             | 1        | Array            |
 | ['abc']  |  2              | 1        | Object Property or Map Element |
 | [3]      |  2              | 1        | Number indexed array/collection |
 | \*       |  2              | 3        | Multiply         |
 | /        |  2              | 3        | Divide           |
 | %        |  2              | 3        | Mod              |
 | +        |  2              | 4        | Plus             |
 | -        |  2              | 4        | Minus            |
 | -        |  2              | 2        | Negative         |
 | >        |  2              | 6        | Greater          |
 | <        |  2              | 6        | Less             |
 | >=       |  2              | 6        | Great Equal      |
 | <=       |  2              | 6        | Less Equal       |
 | ==       |  2              | 6        | Equal            |
 | !=       |  2              | 6        | Not Equal        |
 | ~=       |  2              | 6        | Regexp Match     |
 | !        |  1              | 2        | Not              |
 | !!       |  1              | 2        | Ignore exception and return nil |
 | &&       |  2              | 11       | Logical And      |
 | \|\|     |  2              | 12       | Logical Or       |
 | A\|\|\|B |  2              | 12       | Return B if A is empty or false, else return A |
 | ?:       |  2              | 13       | Ternary          |
 | ~        |  1              | 2        | Bit NOT          |
 | &        |  2              | 7        | Bit AND          |
 | ^        |  2              | 8        | Bit XOR          |
 | \|       |  2              | 9        | Bit OR           |
 | <<       |  2              | 5        | Bit Left Shift   |
 | >>       |  2              | 5        | Bit Right Shift  |


### Like Golang

GEL is completely faithful to Golang basic arithmetic rules and does not do some extensions, such as the most common data type conversions.  
In the process of numerical computation in Golang, the type of the operation result is finally determined according to the type of both sides of the operator.

Example:  

```go
7/3            // return int
(1.0 * 7)/3    // return float64
```

### Some simple examples
#### General operation

```go
gel.Calculate("3+2*5") // Output:  13
```

#### String manipulation
```go
gel.Calculate("'a'+'b'+'c'") // Output:  abc
```

#### struct field
```go
m := map[string]any{
	"pet": struct{
		Name string
	}{"GFW"},
}
gel.Calculate("pet.name", m) // Output:  GFW
```

#### Method call
```go
type Pet struct {
	name string
}

func (p *Pet) SetName(name string) {
	p.name = name
}

func (p *Pet) GetName() string {
	return p.name
}

m := map[string]any{
	"pet": &Pet{},
}
gel.Calculate("pet.SetName('XiaoBai')", m)

gel.Calculate("pet.GetName()", m) // Output:  XiaoBai
```

#### Array element
```go
m := map[string]any{
	"x": []string { "A", "B", "C" },
}

gel.Calculate("x[0]", m) // Output:  A
```

#### Map
```go
m := map[string]any{
	"a": map[string]int{
		"x": 10,
		"y": 5,
	}
}

gel.Calculate("a['x'] * a.y", m) // Output:  50
```

#### Logical
```go
m := map[string]any{
	"a": 5,
}

gel.Calculate("a>10", m) // Output:  false

m["a"] = 20
gel.Calculate("a>10", m) // Output:  true
```

#### A or B
```go
m := map[string]any{
	"obj": map[string]any{},
}
gel.Calculate("!!(obj.pet.name) ||| 'cat'", m) // Output:  cat
```

### strict mode
Defautly, EL use none strict mode (call method of nil object will not cause error)  
Example：  
```go
m := map[string]any{
	"obj": map[string]any{},
}
gel.Calculate("(obj.pet.name) == nil", m)  // true
```

Run in strict mode will return error.
Example:  
```go
m := map[string]any{
	"obj": map[string]any{},
}
gel.CalculateStrict("(obj.pet.name) == nil", m)  // error
```


### How about EL's speed?

I think it's not very fast. The principle of its work is such that each parse passes through 3 steps as below:

  - Parse the expression to a suffix expression array
  - Parse the suffix expression array into a operation tree
  - Evaluate the root node of the operation tree

Of course, I also provide a method to improve efficiency, because if each evaluation passes through these 3 steps is certainly slow, we can precompile it first:

```go
el := gel.Compile("a*10")  // Compile a expression and got a EL instance

m := map[string]any{"a": 10}

el.Calculate(m))  // Output: 100 
```

