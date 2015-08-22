package walgo

import (
	"testing"
)

func TestSomeStuff(t *testing.T) {
	s := struct {
		Foo string `bingo:"yeah,baby" json:"foobar,required"`
	}{
		"oh boy",
	}

	RequireBody(nil, nil, s, func() {
		t.Log("good!")
	})
}
