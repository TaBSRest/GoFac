package SampleStructs

import "errors"

type IndependentStruct struct {}

func NewIndependentStruct() *IndependentStruct {
	return &IndependentStruct{}
}

func NewA() IIndependentStruct {
	return &IndependentStruct{}
}

func NewAWithErr() (IIndependentStruct, error) {
	return &IndependentStruct{}, nil
}

func NewAReturningError() (IIndependentStruct, error) {
	return nil, errors.New("IndependentStruct: Error Forming IndependentStruct!")
}

func (i *IndependentStruct) ReturnNameIndependentStruct() string {
	return "IndependentStruct"
}
