package jsonutils

import (
	"encoding/json"
	"fmt"
)

func FlattenStringOrSetOrMap(m map[string]any, fields []FlattenField) error {
	for _, field := range fields {
		k := field.Name
		raw, ok := m[k]
		if !ok {
			if !field.Optional {
				return fmt.Errorf("missing field: %s", k)
			}
			continue
		}
		if a, ok := raw.([]any); ok {
			if len(a) == 1 {
				m[k] = a[0]
			}
			continue
		}
		if s, ok := raw.(map[string]any); ok {
			if len(s) == 1 {
				for n, v := range s {
					if v == nil || v == "" {
						m[k] = n
					}
				}
			}
			continue
		}
	}
	return nil
}

func UnmarshalFields(m map[string]json.RawMessage, fields []UnmarshalField) error {
	for _, field := range fields {
		raw, ok := m[field.Key]
		if !ok && !field.Optional {
			return fmt.Errorf("missing field: %s", field.Key)
		}
		if len(raw) == 0 {
			continue // Ignore Optional empty fields.
		}
		if err := json.Unmarshal(raw, field.Target); err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalTOrSetOrMap(raw json.RawMessage, target *map[string]string) error {
	if err := json.Unmarshal(raw, target); err == nil {
		return nil
	}
	var a []string
	if err := json.Unmarshal(raw, &a); err == nil {
		*target = make(map[string]string, len(a))
		for _, a := range a {
			(*target)[a] = ""
		}
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return err
	}
	*target = map[string]string{s: ""}
	return nil
}

func UnmarshalTOrSets[T any](m map[string]json.RawMessage, sets []UnmarshalTOrSet[T]) error {
	for _, set := range sets {
		raw, ok := m[set.Key]
		if !ok && !set.Optional {
			return fmt.Errorf("missing field: %s", set.Key)
		}
		if len(raw) == 0 {
			continue // Ignore Optional empty fields.
		}
		var a []T
		if err := json.Unmarshal(raw, &a); err == nil {
			*set.Target = a
			return nil
		}
		var s T
		if err := json.Unmarshal(raw, &s); err != nil {
			return err
		}
		*set.Target = []T{s}
	}
	return nil
}

type FlattenField struct {
	Name     string
	Optional bool
}

type UnmarshalField struct {
	Key      string
	Target   any
	Optional bool
}

type UnmarshalTOrSet[T any] struct {
	Key      string
	Target   *[]T
	Optional bool
}
