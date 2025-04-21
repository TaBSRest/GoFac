package samplestructs

type IIndependentStruct interface {
	ReturnNameIndependentStruct() string
}

type IIndependentStruct2 interface {
	IIndependentStruct
}

type IStructRelyingOnIndependentStructBase interface {
	ReturnNameStructRelyingOnIndependentStruct() string
}

type IStructRelyingOnIndependentStruct interface {
	IStructRelyingOnIndependentStructBase
	ReturnSubStructName() string
}

type IStructRelyingOnIndependentStructs interface {
	IStructRelyingOnIndependentStructBase
	ReturnSubStructNames() []string
}
