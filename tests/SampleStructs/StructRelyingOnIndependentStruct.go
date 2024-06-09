package samplestructs

type StructRelyingOnIndependentStruct struct {
	dependency IIndependentStruct
}

func NewStructRelyingOnIndependentStruct(dependency IIndependentStruct) IStructRelyingOnIndependentStruct {
	return &StructRelyingOnIndependentStruct{
		dependency: dependency,
	}
}

func (i *StructRelyingOnIndependentStruct) ReturnNameStructRelyingOnIndependentStruct() string {
	return "StructRelyingOnIndependentStruct"
}
