package json

import (
	"encoding/json"
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

var (
	ErrFunction = errors.New("cannot encode function to JSON")
	ErrChannel  = errors.New("cannot encode channel to JSON")
	ErrState    = errors.New("cannot encode state to JSON")
	ErrUserData = errors.New("cannot encode userdata to JSON")
	ErrNested   = errors.New("cannot encode recursively nested tables to JSON")
)

type JsonValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j JsonValue) MarshalJSON() ([]byte, error) {
	return ToJSON(j.LValue, j.visited)
}

func ToJSON(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = json.Marshal(converted)
	case lua.LChannel:
		err = ErrChannel
	case lua.LNumber:
		data, err = json.Marshal(converted)
	case *lua.LFunction:
		err = ErrFunction
	case *lua.LNilType:
		data, err = json.Marshal(converted)
	case *lua.LState:
		err = ErrState
	case lua.LString:
		data, err = json.Marshal(converted)
	case *lua.LTable:
		var arr []JsonValue
		var obj map[string]JsonValue

		if visited[converted] {
			panic(ErrNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]JsonValue)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = JsonValue{v, visited}
					return
				}
				arr = append(arr, JsonValue{v, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]JsonValue)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = JsonValue{v, visited}
		})
		if obj != nil {
			data, err = json.Marshal(obj)
		} else {
			data, err = json.Marshal(arr)
		}
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = ErrUserData
	}
	return
}

func FromJSON(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(FromJSON(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), FromJSON(L, item))
		}
		return tbl
	}
	return lua.LNil
}

func NewJsonValue(value lua.LValue) *JsonValue {
	return &JsonValue{
		value,
		make(map[*lua.LTable]bool),
	}
}
