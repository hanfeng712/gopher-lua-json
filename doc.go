// Package json is a simple JSON encode/decoder for gopher-lua.
//
// Documentation
//
// The following functions are exposed by the library:
//  decode(string): Decodes a JSON string. Returns nil and an error string if
//                  the string could not be decoded.
//  encode(value):  Encodes a value into a JSON string. Returns nil and an error
//                  string if the value could not be encoded.
//
// Example
//
// Below is an example usage of the library:
//  L := lua.NewState()
//  luajson.Preload(s)
package json