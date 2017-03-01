# Go-Plugin


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
	plugin.Set("LuaCanCallMe", LuaCanCallMe)
	plugin.LoadPlugin("./plugin.lua")
	ret, _ := plugin.Call("HelloWorld")
	ret, _ := plugin.Call("OnSomeEvent")
	ret, _ := plugin.Call("Square", 2)
}
```

#### plugin.lua
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