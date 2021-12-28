package fastrlp

import (
	"bytes"
	"reflect"

	fuzz "github.com/google/gofuzz"
)

type FuzzObject interface {
	Marshaler
	Unmarshaler
}

type FuzzError struct {
	Source, Target interface{}
}

func (f *FuzzError) Error() string {
	return "failed to encode fuzz object"
}

type FuzzOption func(f *fuzz.Fuzzer) *fuzz.Fuzzer

func WithFuncts(fuzzFuncts ...interface{}) FuzzOption {
	return func(f *fuzz.Fuzzer) *fuzz.Fuzzer {
		return f.Funcs(fuzzFuncts...)
	}
}

func Fuzz(num int, obj FuzzObject, opts ...FuzzOption) error {
	fuzzImpl := func() error {
		f := fuzz.New()
		for _, opt := range opts {
			f = opt(f)
		}
		f.Fuzz(obj)

		data, err := obj.MarshalRLPTo(nil)
		if err != nil {
			return err
		}
		obj2 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(FuzzObject)
		if err := obj2.UnmarshalRLP(data); err != nil {
			return err
		}

		// instead of relying on DeepEqual and issues with zero arrays and so on
		// we use the rlp marshal values to compare
		data2, err := obj2.MarshalRLPTo(nil)
		if err != nil {
			return err
		}
		if !bytes.Equal(data, data2) {
			return &FuzzError{Source: obj, Target: obj2}
		}
		return nil
	}

	for i := 0; i < num; i++ {
		if err := fuzzImpl(); err != nil {
			return err
		}
	}
	return nil
}
