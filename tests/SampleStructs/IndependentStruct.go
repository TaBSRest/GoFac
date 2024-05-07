package samplestructs

import "errors"

type IIndependentStruct interface {
	VoidFuncA()
}

type IndependentStruct struct {}

func NewA() *IndependentStruct {
	return &IndependentStruct{}
}

func NewAWithErr() (*IndependentStruct, error) {
	return &IndependentStruct{}, nil
}

func NewAReturningError() (*IndependentStruct, error) {
	return nil, errors.New("IndependentStruct: Error Forming IndependentStruct!")
}

func (i *IndependentStruct) VoidFuncA() { }
