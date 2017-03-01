# Go-Plugin


#### main.go
```go
package main

import "github.com/alaingilbert/go-plugin"

func LuaCanCallMe() string {
	return "Hello World!"
}

func main() {
	plugin.Init()
	defer plugin.Close()
	plugin.Set("LuaCanCallMe", LuaCanCallMe)
	plugin.LoadPlugin("./myPlugin.lua")
	ret, _ := plugin.Call("myPlugin.HelloWorld")
	ret, _ := plugin.Call("myPlugin.OnSomeEvent")
	ret, _ := plugin.Call("myPlugin.Square", 2)
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
