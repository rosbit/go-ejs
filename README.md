# Embeding JS for golang

go-ejs is a wrapper of [goja](https://github.com/dop251/goja). go-ejs is intended to
embed js in golang applications easily.

### Usage

The package is fully go-getable, So, just type

  `go get github.com/rosbit/go-ejs`

to install.

```go
package main

import "fmt"
import "github.com/rosbit/go-ejs"

func main() {
  ctx := ejs.NewVM(nil)

  res, _ := ctx.Eval("1 + 2")
  fmt.Println("result is:", res)
}
```

### Go calls JavaScript function

Suppose there's a JavaScript file named `a.js` like this:

```javascript
function add(a, b) {
    return a+b
}
```

one can call the JavaScript function add() in Go code like the following:

```go
package main

import "fmt"
import "github.com/rosbit/go-ejs"

var add func(int, int)int

func main() {
  ctx := ejs.NewVM(nil)
  if err := ctx.LoadFile("a.js", nil); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  if err := ctx.BindFunc("add", &add); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  res := add(1, 2)
  fmt.Println("result is:", res)
}
```

### JavaScript calls Go function

JavaScript calling Go function is also easy. Just bind a golang func with a var name
with `AddVar("funcname", function)`. There's the example:

```go
package main

import "github.com/rosbit/go-ejs"

// function to be called by JavaScript
func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  ctx := ejs.NewVM(nil)

  ctx.AddVar("adder", adder)
  ctx.EvalFile("b.js", nil)  // b.js containing code calling "adder"
}
```

In JavaScript code, one can call the registered name directly. There's the example `b.js`.

```javascript
r = adder(1, 100)   // the function "adder" is implemented in Go
console.log(r)
```

### add more than 1 variables and functions at one time

```go
package main

import "github.com/rosbit/go-ejs"
import "fmt"

func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  vars := map[string]interface{}{
     "adder": adder,    // to JavaScript built-in function
     "a": []int{1,2,3}, // to JavaScript array
  }

  ctx := js.NewVM(nil)
  if err := ctx.LoadFile("file.js", vars); err != nil {
     fmt.Printf("%v\n", err)
     return
  }
  // or call ctx.AddVars(vars) to add variables.

  res, err := ctx.GetGlobals("global_var_name") // get the value of var global_var_name
  if err != nil {
     fmt.Printf("%v\n", err)
     return
  }
  fmt.Printf("res:", res)
}
```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.
__Convention:__ fork the repository and make changes on your fork in a feature branch.
