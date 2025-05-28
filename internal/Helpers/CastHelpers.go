package Helpers

import "reflect"

func DereferencePointedArr(pointedArr []*reflect.Value) []reflect.Value {
	arr := make([]reflect.Value, len(pointedArr))
	for i, val := range pointedArr {
		arr[i] = *val
	}
	return arr
}

func CastInput(uncastedInputs []reflect.Value, target []reflect.Type) ([]reflect.Value, error) {
	if len(uncastedInputs) != len(target) {
		return nil, MakeError("GoFac.Resolve", "Something went horribly wrong! The number of inputs to a constructor does not match with the dependencies retrieved!")
	}

	var castedInputs []reflect.Value = make([]reflect.Value, len(uncastedInputs))
	for i, uncastedInput := range uncastedInputs {
		if IsValueArrayOrSlice(uncastedInput) {
			elementaryType := target[i].Elem()
			castedInput := reflect.MakeSlice(reflect.SliceOf(elementaryType), 0, 10)
			for _, input := range uncastedInput.Interface().([]reflect.Value) {
				if !input.CanConvert(elementaryType) {
					return nil, MakeError("GoFac.Resolve", "Cannot convert "+input.Type().Name()+" to "+elementaryType.Name())
				}
				castedInput = reflect.Append(castedInput, input.Convert(elementaryType))
				castedInputs[i] = castedInput
			}
		} else {
			if !uncastedInput.CanConvert(target[i]) {
				return nil, MakeError("GoFac.Resolve", "Cannot convert "+uncastedInput.Type().Name()+" to "+target[i].Name())
			}
			castedInputs[i] = uncastedInput.Convert(target[i])
		}
	}

	return castedInputs, nil
}
