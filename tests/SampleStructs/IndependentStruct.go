package samplestructs

import "errors"

type IIndependentStruct interface {
	VoidFuncA()
}

type IndependentStruct struct {}

func NewA() IIndependentStruct {
	return &IndependentStruct{}
}

func NewAWithErr() (IIndependentStruct, error) {
	return &IndependentStruct{}, nil
}

func NewAReturningError() (IIndependentStruct, error) {
	return nil, errors.New("IndependentStruct: Error Forming IndependentStruct!")
}

func (i *IndependentStruct) VoidFuncA() { }
