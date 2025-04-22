package SampleStructs

type StructRelyingOnIndependentStructs struct {
	dependency []IIndependentStruct
}

func NewStructRelyingOnIndependentStructs(dependency []IIndependentStruct) IStructRelyingOnIndependentStructs {
	return &StructRelyingOnIndependentStructs{
		dependency: dependency,
	}
}

func (*StructRelyingOnIndependentStructs) ReturnNameStructRelyingOnIndependentStruct() string {
	return "StructRelyingOnIndependentStructs"
}

func (s *StructRelyingOnIndependentStructs) ReturnSubStructNames() []string {
	names := make([]string, len(s.dependency))
	for i, dependency := range s.dependency {
		names[i] = dependency.ReturnNameIndependentStruct()
	}
	return names
}
