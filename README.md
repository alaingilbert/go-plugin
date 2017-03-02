# Go-Plugin

## Example

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
	p, _ := plugin.Load("./myPlugin.lua")
	plugin.IsLoaded("./myPlugin.lua") // true
	ret, _ := plugin.Call("myPlugin.HelloWorld")
	ret, _ := p.Call("OnSomeEvent")
	ret, _ := p.Call("Square", 2)
	
	// Call each loaded plugins "HelloWorld" function
	plugin.Each(func(p plugin.Plugin) {
		if ret, err := p.Call("HelloWorld"); err == nil {
			fmt.Println(ret)
		}
	})
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

