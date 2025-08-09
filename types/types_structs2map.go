package types

import "reflect"

func Structs2SliceMap[T any](stSlice []T) []map[string]any {
	length := len(stSlice)
	arrs := make([]map[string]any, length)
	if length > 0 {
		for i := 0; i < length; i++ {
			arrs[i] = Structs2Map(stSlice[i])
		}
	}
	return arrs
}

// Structs2Map must have json tag
func Structs2Map(st any) map[string]any {
	vof := reflect.ValueOf(st)
	tof := reflect.TypeOf(st)
	m := make(map[string]any)
	structReflects(tof, vof, &m)
	return m
}

func structReflects(tof reflect.Type, vof reflect.Value, w *map[string]any) {
	m := *w
	for i := 0; i < tof.NumField(); i++ {
		if tof.Field(i).Anonymous {
			structReflects(tof.Field(i).Type, vof.Field(i), &m)
			continue
		}
		key := tof.Field(i).Tag.Get("json")
		if key != "" && key != "-" {
			value := vof.Field(i).Interface()
			m[key] = value
		}
	}
	w = &m
}
