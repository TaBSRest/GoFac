package Helpers

import "reflect"

func GetName[registrarInterface any]() string {
	reflection := reflect.TypeFor[registrarInterface]()
	return reflection.String()
}
