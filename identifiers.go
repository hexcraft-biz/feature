package feature

import (
	"reflect"
	"strings"
)

// ================================================================
const (
	Delimiter = " "
)

type Identifiers map[string]bool

func NewIdentifiers(input any) Identifiers {
	items := []string{}
	if reflect.TypeOf(input).Kind() == reflect.Slice {
		items = input.([]string)
	} else {
		items = strings.Split(input.(string), Delimiter)
	}

	identifiers := Identifiers{}
	for _, i := range items {
		identifiers.Set(i)
	}

	return identifiers
}

func (i *Identifiers) Set(item string) {
	(*i)[item] = true
}

func (i Identifiers) HasOneOf(sub Identifiers) bool {
	for item, has := range sub {
		if has {
			if val, ok := i[item]; ok && val {
				return true
			}
		}
	}
	return false
}

func (i Identifiers) Contains(sub Identifiers) bool {
	for item, has := range sub {
		if has {
			if val, ok := i[item]; !ok || !val {
				return false
			}
		}
	}
	return true
}

func (i Identifiers) StringSlice() []string {
	ss := []string{}
	for endpoint, has := range i {
		if has {
			ss = append(ss, endpoint)
		}
	}
	return ss
}

func (i Identifiers) AnySlice() []any {
	ss := []any{}
	for endpoint, has := range i {
		if has {
			ss = append(ss, endpoint)
		}
	}
	return ss
}
