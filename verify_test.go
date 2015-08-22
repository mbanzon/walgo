package walgo

import (
	"testing"
)

func TestSomeStuff(t *testing.T) {
	s := struct {
		Foo string `walgo:"skip"`
	}{}

	ok, err := verifyData([]byte(`{"Foo":"bar"}`), &s)
	t.Log(ok, err, s)
}
