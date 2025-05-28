package SampleStructs

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

type ISpecificPerContextStruct interface {
	IIndependentStruct
	GetSpecificID() string
}

type ISingletonWithPerContextDependencyStruct interface {
	IIndependentStruct
	GetInitialPerContextDepName() string
	GetInitialPerContextDepSpecificID() string
}
