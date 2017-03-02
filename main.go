package plugin

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

// L is the lua state
// This is the VM that runs the plugins
var (
	L       *lua.LState
	plugins map[string]Plugin
)

// SetFn ...
// Plugin ...
type Plugin struct {
	Path string
	Name string
}

// Call a function for the specific plugin
func (p *Plugin) Call(fn string, args ...interface{}) (lua.LValue, error) {
	return Call(p.Name+"."+fn, args...)
}

// Unload a specific plugin
func (p *Plugin) Unload() error {
	return Unload(p.Path)
}
func Set(name string, val interface{}) {
	L.SetGlobal(name, luar.New(L, val))
}

func IsLoaded(path string) bool {
	filePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	_, ok := plugins[filePath]
	return ok
}

// Unload remove plugin functions from lua VM
func Unload(path string) error {
	filePath, _ := filepath.Abs(path)
	_, fileName := filepath.Split(filePath)
	fileExt := filepath.Ext(fileName)
	pluginName := strings.TrimSuffix(fileName, fileExt)
	str := "\nlocal P = {}\n" + pluginName + " = P\nsetmetatable(" + pluginName + ", {__index = _G})\nsetfenv(1, P)\n"
	if err := L.DoString(str); err != nil {
		return errors.Wrap(err, "Unload: Failed to unload plugin")
	}
	delete(plugins, filePath)
	return nil
}

func Load(path string) (Plugin, error) {
	var newPlugin Plugin
	filePath, _ := filepath.Abs(path)
	_, fileName := filepath.Split(filePath)
	fileExt := filepath.Ext(fileName)
	pluginName := strings.TrimSuffix(fileName, fileExt)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return newPlugin, errors.Wrap(err, "Load: Unable to read file")
	}
	pluginDef := "\nlocal P = {}\n" + pluginName + " = P\nsetmetatable(" + pluginName + ", {__index = _G})\nsetfenv(1, P)\n"
	if err := L.DoString(pluginDef + string(data)); err != nil {
		return newPlugin, errors.Wrap(err, "Load: Unable to execute lua string")
	}
	newPlugin = Plugin{filePath, pluginName}
	plugins[filePath] = newPlugin
	return newPlugin, nil
}

func Call(fn string, args ...interface{}) (lua.LValue, error) {
	var luaFunc lua.LValue
	if strings.Contains(fn, ".") {
		plugin := L.GetGlobal(strings.Split(fn, ".")[0])
		if plugin.String() == "nil" {
			return nil, errors.New("function does not exist: " + fn)
		}
		luaFunc = L.GetField(plugin, strings.Split(fn, ".")[1])
	} else {
		luaFunc = L.GetGlobal(fn)
	}
	if luaFunc.String() == "nil" {
		return nil, errors.New("function does not exist: " + fn)
	}
	var luaArgs []lua.LValue
	for _, v := range args {
		luaArgs = append(luaArgs, luar.New(L, v))
	}
	err := L.CallByParam(lua.P{
		Fn:      luaFunc,
		NRet:    1,
		Protect: true,
	}, luaArgs...)
	if err != nil {
		return nil, err
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)         // remove received value
	return ret, nil
}
}

func Init() {
	L = lua.NewState()
	plugins = make(map[string]Plugin)
}

func Close() {
	L.Close()
}
