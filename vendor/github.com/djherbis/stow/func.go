package stow

import (
	"fmt"
	"reflect"
)

type funcCall struct {
	s *Store

	Value reflect.Value
	Type  reflect.Type

	hasKey   bool
	isKeyPtr bool
	keyType  reflect.Type

	isValPtr bool
	valType  reflect.Type
}

func newFuncCall(s *Store, fn interface{}) (fc funcCall, err error) {
	fc.s = s
	fc.Value = reflect.ValueOf(fn)
	fc.Type = fc.Value.Type()
	if fc.Value.Kind() != reflect.Func {
		return fc, fmt.Errorf("fn is not a func()")
	}

	if fc.Type.NumIn() == 1 {
		fc.setValue(fc.Type.In(0))
	} else if fc.Type.NumIn() == 2 {
		fc.setKey(fc.Type.In(0))
		fc.setValue(fc.Type.In(1))
	} else {
		return fc, fmt.Errorf("bad number of args in ForEach func()")
	}

	return fc, nil
}

func isPtr(typ reflect.Type) bool { return typ.Kind() == reflect.Ptr }

func (fc *funcCall) setKey(typ reflect.Type) {
	fc.hasKey = true
	fc.keyType = typ
	if fc.isKeyPtr = isPtr(fc.keyType); fc.isKeyPtr {
		fc.keyType = fc.keyType.Elem()
	}
}

func (fc *funcCall) setValue(typ reflect.Type) {
	fc.valType = typ
	if fc.isValPtr = isPtr(fc.valType); fc.isValPtr {
		fc.valType = fc.valType.Elem()
	}
}

func (fc *funcCall) getKey(v []byte) (key reflect.Value, err error) {
	if fc.keyType.Kind() == reflect.String {
		return reflect.ValueOf(string(v)), nil
	} else if fc.keyType.Kind() == reflect.Slice && fc.keyType.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(v), nil
	}

	key = reflect.New(fc.keyType)

	if err := fc.s.unmarshal(v, key.Interface()); err != nil {
		return key, err
	}

	if !fc.isKeyPtr {
		key = deref(key)
	}

	return key, err
}

func (fc *funcCall) getValue(v []byte) (val reflect.Value, err error) {
	val = reflect.New(fc.valType)

	if err := fc.s.unmarshal(v, val.Interface()); err != nil {
		return val, err
	}

	if !fc.isValPtr {
		val = deref(val)
	}

	return val, err
}

func (fc *funcCall) call(k, v []byte) error {
	val, err := fc.getValue(v)
	if err != nil {
		return err
	}

	if !fc.hasKey {
		fc.Value.Call([]reflect.Value{val})
		return nil
	}

	key, err := fc.getKey(k)
	if err != nil {
		return err
	}
	fc.Value.Call([]reflect.Value{key, val})
	return nil
}

func deref(val reflect.Value) reflect.Value {
	if val.IsValid() {
		return val.Elem()
	}
	return reflect.Zero(val.Type())
}
