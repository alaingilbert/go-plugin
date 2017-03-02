# Go-Plugin

## Example

#### main.go
```go
package main

import (
	"fmt"
	"github.com/alaingilbert/go-plugin"
)

func LuaCanCallMe() string {
	return "Hello World!"
}

func main() {
	plugin.Init()
	defer plugin.Close()

	// Make the function available inside the plugins
	plugin.Set("LuaCanCallMe", LuaCanCallMe)

	p, err := plugin.Load("./myPlugin.lua")
	plugin.IsLoaded("./myPlugin.lua") // true

	ret1, err := plugin.Call("myPlugin.HelloWorld")
	ret2, err := p.Call("OnSomeEvent")

	var squared int
	err = p.CallUnmarshal(&squared, "Square", 2)

	// Call each loaded plugins "HelloWorld" function
	plugin.Each(func(p plugin.Plugin) {
		if ret, err := p.Call("HelloWorld"); err == nil {
			fmt.Println(ret)
		}
	})

	fmt.Println(ret1, ret2, err)
}
```

#### myPlugin.lua
```lua
function HelloWorld()
    return LuaCanCallMe()
end

function OnSomeEvent()
    return "Some event process"
end

function Square(x)
    return x * x
end
```

## Documentation

https://godoc.org/github.com/alaingilbert/go-plugin

## Thanking
Heavily inpired by `micro` plugins system  
https://github.com/zyedidia/micro

