package SampleStructs

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

func (i *StructRelyingOnIndependentStruct) ReturnSubStructName() string {
	return i.dependency.ReturnNameIndependentStruct()
}
