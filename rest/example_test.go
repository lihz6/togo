package rest_test

import (
	"errors"
	"togo/rest"
)

func ExampleWithRecover() {
	go rest.WithRecover(func() {
		rest.Panic(errors.New("rest.WithRecover panic"))
	})
}

func ExampleRecover() {
	defer rest.Recover()
	rest.Panic(errors.New("rest.Recover panic"))
}
