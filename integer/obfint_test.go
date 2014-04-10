package integer

import (
	"math/rand"
	"testing"
)

var testPrime int64 = 982450871

func TestHide(t *testing.T) {
	coder, err := NewCoder(testPrime, rand.Int63())
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 100; i++ {
		id := NewObfInt(rand.Int63(), coder)
		hidden := id.Hide()
		id2 := NewObfIntFromHidden(hidden, coder)
		if id.Int64() != id2.Int64() {
			t.Error("Hidden and reveal error.")
		}
	}
}
