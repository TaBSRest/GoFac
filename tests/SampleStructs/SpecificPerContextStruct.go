package SampleStructs

import "github.com/google/uuid"

type SpecificPerContextStruct struct {
	IIndependentStruct
	specificID string
}

func NewSpecificPerContextStruct(instance IIndependentStruct) ISpecificPerContextStruct {
	return &SpecificPerContextStruct{
		IIndependentStruct: instance,
		specificID:         uuid.NewString(),
	}
}

func (s *SpecificPerContextStruct) GetSpecificID() string {
	return s.specificID
}
