package samplestructs

type IIndependentStruct interface {
	ReturnNameIndependentStruct() string
}

type IStructRelyingOnIndependentStruct interface {
	ReturnNameStructRelyingOnIndependentStruct() string
}

type IStructRelyingOnIndependentStructs interface {
	IStructRelyingOnIndependentStruct
	ReturnSubStructNames() []string
}
