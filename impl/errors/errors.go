package errors

import "reflect"

func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Unwrap(err error) error {
	u, ok := err.(interface{ Unwrap() error })
	if ok {
		return u.Unwrap()
	}
	return nil
}

func Is(err, target error) bool {
	if target == nil {
		return err == nil
	}
	cmp := reflect.TypeOf(target).Comparable()
	for {
		if cmp && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func As(err error, target any) bool {
	if target == nil {
		panic("target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}
	ele := typ.Elem()
	if ele.Kind() != reflect.Interface && !ele.Implements(errorType) {
		panic("*target must be interface or implement error")
	}
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(ele) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
