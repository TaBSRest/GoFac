package samplestructs

import (
	"errors"
	"sync"
)

type SingletonIndependentStruct struct {
	once sync.Once
}

func NewS() IIndependentStruct {
	return &SingletonIndependentStruct{}
}

func NewSWithErr() (IIndependentStruct, error) {
	return &SingletonIndependentStruct{}, nil
}

func NewSReturningError() (IIndependentStruct, error) {
	return nil, errors.New("SingletonIndependentStruct: Error Forming SingletonIndependentStruct!")
}

func (i *SingletonIndependentStruct) ReturnNameIndependentStruct() string {
	message := "SingletonStruct"
	alreadyRan := "Already Ran!"
	returnMessage := alreadyRan
	i.once.Do(func() {
		returnMessage = message
	})
	return returnMessage
}
