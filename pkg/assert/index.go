package assert

import (
	"bytes"
	"reflect"
)

func IsNil(actual any, msg string) {
	if actual != nil {
		panic(msg)
	}
}

func ErrorNil(err error, msg string) {
	if err != nil {
		panic(err.Error() + " " + msg)
	}
}

func IsEqual(expected, actual interface{}, msg string) {
	if actual == nil || expected == nil {
		if actual != expected {
			panic(msg)
		}
	}

	exp, ok := expected.([]byte)
	if !ok {
		if reflect.DeepEqual(expected, actual) {
			panic(msg)
		}
	}

	act, ok := actual.([]byte)
	if !ok {
		panic(msg)
	}

	if act == nil && exp != nil {
		panic(msg)
	} else if act != nil && exp == nil {
		panic(msg)
	} else if act == nil && exp == nil {
		return
	}

	if !bytes.Equal(act, exp) {
		panic(msg)
	}
}
