package plugin_test

import (
	"testing"

	"github.com/alaingilbert/go-plugin"
)

func TestCall(t *testing.T) {
	plugin.Init()
	defer plugin.Close()
	p, _ := plugin.Load("plugin_test.lua")
	ret, err := p.Call("Square", 3)
	if err != nil {
		t.Error("Unable to call Square", err)
	}
	if ret.String() != "9" {
		t.Error("Should be 9", ret.String())
	}
}

func TestCallUnmarshal(t *testing.T) {
	plugin.Init()
	defer plugin.Close()
	p, _ := plugin.Load("plugin_test.lua")

	// Test with number
	var squared int
	if err := p.CallUnmarshal(&squared, "Square", 3); err != nil {
		t.Error("There should be no error", err)
	}
	if squared != 9 {
		t.Error("squared should be 9", squared)
	}

	// Test with string
	var stringTest string
	p.CallUnmarshal(&stringTest, "StringTest")
	if stringTest != "a string" {
		t.Error("stringTest should be 'a string'", stringTest)
	}

	// Test with bool
	var boolTest bool
	p.CallUnmarshal(&boolTest, "BoolTest")
	if boolTest != true {
		t.Error("boolTest should be true", boolTest)
	}
}
