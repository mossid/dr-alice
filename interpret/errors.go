package radicle

import (
	"errors"
	"fmt"
)

func newError(ty string, fn string, format string, args ...interface{}) error {
	if fn != "" {
		fn = "(" + fn + ")"
	}
	return errors.New(ty + fn + ": " + fmt.Sprintf(format, args...))
}

func SpecialFormError(fn string, desc string) error {
	return newError("SpecialForm", fn, desc)
}

func WrongNumberArgsError(fn string, exp int, act int) error {
	return newError("WrongNumberArgs", fn, "expected %d but %d", exp, act)
}

func TypeError(fn string, exp Value, act Value) error {
	return newError("TypeError", fn, "expected %T but %T", exp, act)
}

func PatternMatchError(fn string, desc string) error {
	return newError("PatternMatch", fn, desc)
}

func NonFunctionCalledError(fun Value) error {
	return newError("NonFunctionCalled", "", "%+v", fun)
}

func UnknownIdentifierError(id Ident) error {
	return newError("UnknownIdentifier", "", id)
}

func ImpossibleError(fn string, desc string) error {
	return newError("Impossible", fn, desc)
}

func OtherError(fn string, desc string) error {
	return newError("Other", fn, desc)
}
