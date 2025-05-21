package SampleStructs

import "errors"

type IndependentStructB struct {
	something string // To not make it a Zero-Sized Struct
}

func NewB() IIndependentStruct {
	return &IndependentStructB{}
}

func NewBWithErr() (IIndependentStruct, error) {
	return &IndependentStructB{}, nil
}

func NewBReturningError() (IIndependentStruct, error) {
	return nil, errors.New("IndependentStructB: Error Forming IndependentStructB!")
}

func (i *IndependentStructB) ReturnNameIndependentStruct() string {
	return "IndependentStructB"
}
