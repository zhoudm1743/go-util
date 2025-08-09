package types

import "reflect"

func ArrayColumn[T, V any](array []T, k string) []V {
	l := len(array)
	if l == 0 {
		return nil
	}
	values := make([]V, len(array))
	switch reflect.TypeOf(array).Elem().Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < l; i++ {
			values[i] = reflect.ValueOf(array[i]).Index(int(reflect.ValueOf(k).Int())).Interface().(V)
		}
		break
	case reflect.Map:
		for i := 0; i < l; i++ {
			values[i] = reflect.ValueOf(array[i]).MapIndex(reflect.ValueOf(k)).Interface().(V)
		}
		break
	case reflect.Struct:
		for i := 0; i < l; i++ {
			values[i] = reflect.ValueOf(array[i]).FieldByName(reflect.ValueOf(k).String()).Interface().(V)
		}
		break
	default:
		return nil
	}
	return values
}
