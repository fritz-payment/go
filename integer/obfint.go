package integer

import (
	"fmt"
	"math"
)

var maxId uint64 = math.MaxInt64

type Coder struct {
	prime, primeInv, xor int64
}

func NewCoder(prime, xor int64) (*Coder, error) {
	c := &Coder{
		prime: prime,
		xor:   xor,
	}
	primeInv, err := modinv(uint64(prime), maxId+1)
	if err != nil {
		return nil, err
	}
	c.primeInv = int64(primeInv)
	return c, nil
}

func (c *Coder) Hide(val int64) int64 {
	hidden := ((uint64(val) * uint64(c.prime)) & maxId) ^ uint64(c.xor)
	return int64(hidden)
}

func (c *Coder) Show(hidden int64) int64 {
	shown := ((uint64(hidden) ^ uint64(c.xor)) * uint64(c.primeInv)) & maxId
	return int64(shown)
}

func modinv(a, m uint64) (uint64, error) {
	x, _, g := egcd(a, m)
	if g == 1 {
		return x % m, nil
	}
	return 0, fmt.Errorf("Could not determine modinv of %d, %d.", a, m)
}

func egcd(a, b uint64) (uint64, uint64, uint64) {
	var x, x1, y, y1 uint64 = 1, 0, 0, 1
	for b != 0 {
		q := a / b
		x, x1 = x1, x-q*x1
		y, y1 = y1, y-q*y1
		a, b = b, a-q*b
	}
	return x, y, a
}

type ObfInt struct {
	i int64
	c *Coder
}

func NewObfInt(plain int64, c *Coder) *ObfInt {
	return &ObfInt{plain, c}
}

func NewObfIntFromHidden(hidden int64, c *Coder) *ObfInt {
	return &ObfInt{c.Show(hidden), c}
}

func (o ObfInt) Int64() int64 {
	return o.i
}

func (o ObfInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", o.Hide())), nil
}

func (o ObfInt) Hide() int64 {
	if o.c == nil {
		panic("Coder not set up.")
	}
	return o.c.Hide(o.i)
}
