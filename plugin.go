package plugin

import (
	errorsPkg "errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golang/go/src/sort"
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

// Plugin ...
type Plugin struct {
	Path string
	Name string
}

// Call a function for the specific plugin
func (p *Plugin) Call(fn string, args ...interface{}) (lua.LValue, error) {
	return Call(p.Name+"."+fn, args...)
}

// Call function, put result in v
func (p *Plugin) CallUnmarshal(v interface{}, fn string, args ...interface{}) error {
	return CallUnmarshal(v, p.Name+"."+fn, args...)
}

// Unload a specific plugin
func (p *Plugin) Unload() error {
	return Unload(p.Path)
}

// Reload a specific plugin
func (p *Plugin) Reload() error {
	_, err := Load(p.Path)
	return err
}

// Each Call clb with each plugins sorted by name
func Each(clb func(p Plugin)) {
	var loadedPlugins []Plugin
	for _, v := range plugins {
		loadedPlugins = append(loadedPlugins, v)
	}
	sort.Slice(loadedPlugins, func(i, j int) bool { return loadedPlugins[i].Name < loadedPlugins[j].Name })
	for _, v := range loadedPlugins {
		clb(v)
	}
}

// Set a global variable in lua VM
func Set(name string, val interface{}) {
	L.SetGlobal(name, luar.New(L, val))
}

// IsLoaded returns whether the file is loaded or not in the lua VM
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

// Load plugin functions in lua VM
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

// Call a function in the lua VM
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

func CallUnmarshal(v interface{}, fn string, args ...interface{}) error {
	rv := reflect.ValueOf(v)
	nv := reflect.Indirect(rv)
	lv, err := Call(fn, args...)
	switch v.(type) {
	case *string:
		if str, ok := lv.(lua.LString); ok {
			nv.SetString(string(str))
		}
	case *int:
		if nb, ok := lv.(lua.LNumber); ok {
			nv.SetInt(int64(nb))
		}
	case *bool:
		if b, ok := lv.(lua.LBool); ok {
			nv.SetBool(bool(b))
		}
	default:
		err = errorsPkg.New("Invalid type")
	}
	return err
}

// Init the lua VM
func Init() {
	L = lua.NewState()
	plugins = make(map[string]Plugin)
}

// Close lua VM
func Close() {
	L.Close()
}
