package fastrlp

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func TestEncodingRawRandom(t *testing.T) {
	for i := 0; i < 10000; i++ {
		v0 := generateRandom()
		buf := v0.MarshalTo(nil)

		p := &Parser{}
		v, err := p.Parse(buf)
		if err != nil {
			t.Fatal(err)
		}

		buf1 := v.MarshalTo(nil)
		if !bytes.Equal(buf, buf1) {
			t.Fatal("bad")
		}
		if !checkRaw(p, v, v0) {
			t.Fatal("bad")
		}
	}
}

func checkRaw(p *Parser, v *Value, v0 *Value) bool {
	if v.Type() == TypeArray {
		elems, err := v.GetElems()
		if err != nil {
			panic(err)
		}
		for indx, elem := range elems {
			if !checkRaw(p, elem, v0.Get(indx)) {
				return false
			}
		}
	}
	buf := p.Raw(v)
	if !bytes.Equal(buf, v.MarshalTo(nil)) {
		return false
	}
	return true
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func generateRandom() *Value {
	rand.Seed(time.Now().UTC().UnixNano())
	return generateRandomImpl(&Arena{}, 0)
}

func generateRandomImpl(a *Arena, depth uint) *Value {
	if randomInt(0, 10) < 5 || depth > 2 {
		// value
		buf := make([]byte, randomInt(1, 100))
		rand.Read(buf)
		return a.NewBytes(buf)
	}
	// array
	v := a.NewArray()
	for i := 0; i < randomInt(1, 5); i++ {
		v.Set(generateRandomImpl(a, depth+1))
	}
	return v
}

func TestInvalidRlp(t *testing.T) {
	// cases from the official spec

	cases := []string{
		"bf0f000000000000021111",
		"ff0f000000000000021111",
		"f80180",
		"f80100",
		"b9002100dc2b275d0f74e8a53e6f4ec61b27f24278820be3f82ea2110e582081b0565df0",
		"f861f83eb9002100dc2b275d0f74e8a53e6f4ec61b27f24278820be3f82ea2110e582081b0565df027b90015002d5ef8325ae4d034df55d4b58d0dfba64d61ddd17be00000b9001a00dae30907045a2f66fa36f2bb8aa9029cbb0b8a7b3b5c435ab331",
		"8100",
		"8101",
		"817F",
	}

	p := &Parser{}
	for _, c := range cases {
		buf, err := hex.DecodeString(c)
		if err != nil {
			t.Fatal(err)
		}

		if _, err = p.Parse(buf); err == nil {
			t.Fatal("it should fail")
		}
	}
}
