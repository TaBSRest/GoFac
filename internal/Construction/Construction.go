package Construction

import (
	"reflect"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
)

type Construction struct {
	Info reflect.Type
	Value reflect.Value
}

func NewConstruction(info reflect.Type, value reflect.Value) (Construction, error) {
	if info == nil {
		return Construction{}, h.MakeError("Construtor.New", "info is nil!")
	}

	if value.Interface() == nil {
		return Construction{}, h.MakeError("Construtor.New", "value is nil!")
	}

	return Construction{
		Info: info,
		Value: value,
	}, nil
}
