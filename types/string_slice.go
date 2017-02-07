package types

import (
	"sort"
)

type StringSlice []string

func (slice StringSlice) Uniq() (o StringSlice, err error) {
	if len(slice) <= 0 {
		return
	}

	cache := make(map[string]string)
	for _, t := range slice {
		cache[t] = t
	}

	for k, _ := range cache {
		o = append(o, k)
	}
	return
}

func (slice StringSlice) Sort() (o StringSlice, err error) {
	(sort.StringSlice)(slice).Sort()
	return slice, nil
}
