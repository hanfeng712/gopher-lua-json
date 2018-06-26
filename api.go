package json

import (
	"encoding/json"

	lua "github.com/coinstack/gopher-lua"
)

var api = map[string]lua.LGFunction{
	"decode": apiDecode,
	"encode": apiEncode,
}

func apiDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(FromJSON(L, value))
	return 1
}

func apiEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := ToJSON(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}
