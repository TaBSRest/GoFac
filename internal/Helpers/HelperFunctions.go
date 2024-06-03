package helpers

import "reflect"

func GetName[registrarInterface interface{}]() string {
	reflection := reflect.TypeFor[registrarInterface]()
	return GetNameFromType(reflection)
}

func GetNameFromType(reflection reflect.Type) string {
	return reflection.PkgPath() + "/" + reflection.Name()
}
