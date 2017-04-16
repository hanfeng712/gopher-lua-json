package json

import (
	"encoding/json"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func getValue(state *lua.LState, def string) lua.LValue {
	_ = state.DoString(def)
	value := state.Get(-1)
	return value
}

func TestMarshal(t *testing.T) {
	state := lua.NewState()

	stringValue := getValue(state, `return { foo = 314}`)
	// stringValue := getValue(state, `return 314`)

	tempLuaObj := NewJsonValue(stringValue)
	jsonValue, err := tempLuaObj.MarshalJSON()
	if nil != err {
		t.Error(err)
	}

	t.Logf("%v\n", string(jsonValue))

	// var parsedJson map[string]interface{}
	var parsedJson interface{}
	json.Unmarshal(jsonValue, &parsedJson)
	lValue := FromJSON(state, parsedJson)

	tempLuaObj2 := NewJsonValue(lValue)
	jsonValue2, err := tempLuaObj2.MarshalJSON()
	if nil != err {
		t.Error(err)
	}
	t.Logf("%v\n", string(jsonValue2))
}
