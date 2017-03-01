package plugin

import (
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

// L is the lua state
// This is the VM that runs the plugins
var L *lua.LState

// SetFn ...
func Set(name string, val interface{}) {
	L.SetGlobal(name, luar.New(L, val))
}

func LoadPlugin(path string) error {
	if err := L.DoFile(path); err != nil {
		return err
	}
	return nil
}

func Call(fn string, args ...interface{}) (ret lua.LValue, err error) {
	var luaArgs []lua.LValue
	for _, v := range args {
		luaArgs = append(luaArgs, luar.New(L, v))
	}
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal(fn),
		NRet:    1,
		Protect: true,
	}, luaArgs...)
	if err != nil {
		return
	}
	ret = L.Get(-1) // returned value
	L.Pop(1)        // remove received value
	return
}

func Init() {
	L = lua.NewState()
}

func Close() {
	L.Close()
}
