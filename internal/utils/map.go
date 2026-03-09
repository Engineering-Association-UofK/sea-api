package utils

import (
	"errors"
	"fmt"
	"strings"
)

var errKeyExists = errors.New("key already exists in map")
var errMissingKey = errors.New("key does not exist in map")

type Mpp[K comparable, V any] map[K]V

// Constructor

func NewMpp[K comparable, V any]() Mpp[K, V] {
	return Mpp[K, V]{}
}

// Get Value

func (mpp Mpp[K, V]) Value(key K) (V, error) {
	value, ok := mpp[key]
	if ok {
		return value, nil
	}
	var zero V
	return zero, errMissingKey
}

func (mpp Mpp[K, V]) GetOrCreate(key K, create func() V) V {
	if value, ok := mpp[key]; ok {
		return value
	}
	value := create()
	mpp[key] = value
	return value
}

// Add

func (mpp Mpp[K, V]) Add(key K, value V) error {
	if _, exists := mpp[key]; exists {
		return errKeyExists
	}
	mpp[key] = value
	return nil
}

// Change Value

func (mpp Mpp[K, V]) Update(key K, value V) error {
	if _, exists := mpp[key]; exists {
		mpp[key] = value
		return nil
	}
	return errMissingKey
}

func (mpp Mpp[K, V]) ForEach(fn func(K, V)) {
	for k, v := range mpp {
		fn(k, v)
	}
}

func MapValues[K comparable, V any, T any](m *Mpp[K, V], fn func(V) T) Mpp[K, T] {
	result := NewMpp[K, T]()

	for k, v := range *m {
		result[k] = fn(v)
	}

	return result
}

// Delete

func (mpp Mpp[K, V]) Delete(key K) error {
	if _, exists := mpp[key]; exists {
		delete(mpp, key)
		return nil
	}
	return errMissingKey
}

func (mpp Mpp[K, V]) Empty() {
	clear(mpp)
}

// Makes no changes

func (mpp Mpp[K, V]) Len() int {
	return len(mpp)
}

func (mpp Mpp[K, V]) Exists(key K) bool {
	_, ok := mpp[key]
	return ok
}

// Specialized

func (mpp Mpp[K, V]) Filter(fn func(K, V) bool) Mpp[K, V] {
	result := NewMpp[K, V]()

	for k, v := range mpp {
		if fn(k, v) {
			_ = result.Add(k, v)
		}
	}

	return result
}

func (mpp Mpp[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("{\n")
	if mpp.Len() == 0 {
		sb.WriteString("}")
		return sb.String()
	}
	i := 0
	for k, v := range mpp {
		if i > 0 {
			sb.WriteString(",")
		}
		s := fmt.Sprintf("   \"%v\": \"%v\"", k, v)
		sb.WriteString(s)
		i++
	}
	sb.WriteString("\n}")

	return sb.String()
}

func (mpp Mpp[K, V]) Clone() Mpp[K, V] {
	result := NewMpp[K, V]()

	for k, v := range mpp {
		result[k] = v
	}

	return result
}

func (mpp Mpp[K, V]) ToMap() map[K]V {
	result := make(map[K]V, len(mpp))
	for k, v := range mpp {
		result[k] = v
	}
	return result
}

func FromSlice[K comparable, V any](slice []V, getKey func(V) K) Mpp[K, V] {
	m := NewMpp[K, V]()
	for _, item := range slice {
		m[getKey(item)] = item
	}
	return m
}

// Lists

func (mpp Mpp[K, V]) Keys() []K {
	keys := make([]K, mpp.Len())
	i := 0
	for k := range mpp {
		keys[i] = k
		i++
	}
	return keys
}

func (mpp Mpp[K, V]) Values() []V {
	values := make([]V, len(mpp))
	i := 0
	for _, v := range mpp {
		values[i] = v
		i++
	}
	return values
}
