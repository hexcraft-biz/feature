package feature

import (
	"reflect"
)

// ================================================================
type Identifiers map[any]bool

func NewIdentifiers(items any) Identifiers {
	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		return nil
	}

	identifiers := Identifiers{}
	for i := 0; i < val.Len(); i += 1 {
		identifiers.Set(val.Index(i).Interface())
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
