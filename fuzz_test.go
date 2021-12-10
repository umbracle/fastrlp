package fastrlp

import (
	"fmt"
	"testing"
)

func TestFuzzFramework(t *testing.T) {
	obj := &Simple{}
	if err := Fuzz(100, obj); err != nil {
		t.Fatal(err)
	}
}

type Simple struct {
	Data1 []byte
	Data2 [][]byte
	Data3 uint64
}

func (s *Simple) MarshalRLPTo(dst []byte) ([]byte, error) {
	return MarshalRLP(s)
}

func (s *Simple) MarshalRLPWith(ar *Arena) (*Value, error) {
	vv := ar.NewArray()

	// Data1
	if len(s.Data1) == 0 {
		vv.Set(ar.NewNull())
	} else {
		vv.Set(ar.NewBytes(s.Data1))
	}

	// Data2
	if len(s.Data2) == 0 {
		vv.Set(ar.NewNullArray())
	} else {
		committed := ar.NewArray()
		for _, a := range s.Data2 {
			if len(a) == 0 {
				committed.Set(ar.NewNull())
			} else {
				committed.Set(ar.NewBytes(a))
			}
		}
		vv.Set(committed)
	}

	// Data3
	vv.Set(ar.NewUint(s.Data3))

	return vv, nil
}

func (s *Simple) UnmarshalRLP(buf []byte) error {
	return UnmarshalRLP(buf, s)
}

func (s *Simple) UnmarshalRLPWith(v *Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}
	if num := len(elems); num != 3 {
		return fmt.Errorf("not enough elements to decode extra, expected 3 but found %d", num)
	}

	// Data1
	{
		if s.Data1, err = elems[0].GetBytes(s.Data1); err != nil {
			return err
		}
	}

	// Data2
	{
		vals, err := elems[1].GetElems()
		if err != nil {
			return fmt.Errorf("list expected for committed")
		}
		if len(vals) == 0 {
			s.Data2 = nil
		} else {
			s.Data2 = make([][]byte, len(vals))
			for indx, val := range vals {
				if s.Data2[indx], err = val.GetBytes(s.Data2[indx]); err != nil {
					return err
				}
			}
		}
	}

	// Data3
	if s.Data3, err = elems[2].GetUint64(); err != nil {
		return err
	}

	return nil
}
