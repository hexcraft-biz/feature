package feature

import (
	"reflect"
)

// ================================================================
const (
	Delimiter = " "
)

type Identifiers map[any]bool

func NewIdentifiers(input any) Identifiers {
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		return nil
	}

	items, ok := input.([]any)
	if !ok {
		return nil
	}

	identifiers := Identifiers{}
	for _, i := range items {
		identifiers.Set(i)
	}

	return identifiers
}

func (i *Identifiers) Set(item any) {
	(*i)[item] = true
}

func (i Identifiers) AnySlice() []any {
	ss := []any{}
	for s, has := range i {
		if has {
			ss = append(ss, s)
		}
	}
	return ss
}
