package wasmgo

// Undefined isn't defined, but we have to make it _something_.
var Undefined = &struct{}{}

// Object is a generic JS-like object.
type Object struct {
	Props map[string]interface{}
	New   func(args []interface{}) interface{}
}

// ArrayBuffer is a wrapper around arrays for "JS" interop
type ArrayBuffer struct {
	data []byte
}

// TypedArray is an ArrayBuffer with an offset and length
type TypedArray struct {
	Buffer *ArrayBuffer
	Offset int
	Length int
}

func (a *TypedArray) contents() []byte {
	return a.Buffer.data[a.Offset : a.Offset+a.Length]
}

var typedArrayClass = &Object{
	New: func(args []interface{}) interface{} {
		return &TypedArray{
			Buffer: args[0].(*ArrayBuffer),
			Offset: int(args[1].(float64)),
			Length: int(args[2].(float64)),
		}
	},
}

type FuncWrapper struct {
	id interface{}
}
