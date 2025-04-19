package Helpers

import "reflect"

func GetName[registrarInterface any]() string {
	reflection := reflect.TypeFor[registrarInterface]()
	return GetNameFor(reflection)
}

func GetNameFor(reflection reflect.Type) string {
	return reflection.PkgPath() + "/" + reflection.Name()
}
