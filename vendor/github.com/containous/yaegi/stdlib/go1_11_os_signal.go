// Code generated by 'goexports os/signal'. DO NOT EDIT.

// +build go1.11,!go1.12

package stdlib

import (
	"os/signal"
	"reflect"
)

func init() {
	Symbols["os/signal"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Ignore":  reflect.ValueOf(signal.Ignore),
		"Ignored": reflect.ValueOf(signal.Ignored),
		"Notify":  reflect.ValueOf(signal.Notify),
		"Reset":   reflect.ValueOf(signal.Reset),
		"Stop":    reflect.ValueOf(signal.Stop),
	}
}
