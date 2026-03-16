package ContainerBuilder

import (
	"fmt"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

// analyzeDependencies ensures that singleton registrations do not depend on
// registrations with PerScope lifetime. It returns an error if such a
// dependency is found.
func (cb *ContainerBuilder) analyzeDependencies() error {
	// Track visited registrations to avoid processing the same registration
	// multiple times during the root iteration.
	visited := make(map[*r.Registration]bool)
	for _, regs := range cb.cache {
		for _, reg := range regs {
			if visited[reg] {
				continue
			}
			visited[reg] = true
			if reg.Options.Scope == s.Singleton {
				if err := cb.checkSingletonDependencies(reg, make(map[*r.Registration]bool)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// checkSingletonDependencies recursively checks that the provided singleton
// registration does not depend (directly or indirectly) on any PerScope
// registration.
func (cb *ContainerBuilder) checkSingletonDependencies(reg *r.Registration, stack map[*r.Registration]bool) error {
	if stack[reg] {
		return nil
	}
	stack[reg] = true

	ctorInfo := reg.Construction.Info
	for i := 0; i < ctorInfo.NumIn(); i++ {
		depType := ctorInfo.In(i)

		// If dependency is an array or slice, analyze each registration of the element type.
		if h.IsArrayOrSlice(depType) {
			elemType := depType.Elem()
			depRegs, found := cb.cache[elemType]
			if !found {
				continue
			}
			for _, dr := range depRegs {
				if dr.Options.Scope == s.PerScope {
					return te.New(fmt.Sprintf("Singleton %s depends on PerScope registration %s", ctorInfo, elemType))
				}
				if err := cb.checkSingletonDependencies(dr, stack); err != nil {
					return err
				}
			}
			continue
		}

		depRegs, found := cb.cache[depType]
		if !found || len(depRegs) == 0 {
			continue
		}
		depReg := depRegs[len(depRegs)-1]
		if depReg.Options.Scope == s.PerScope {
			return te.New(fmt.Sprintf("Singleton %s depends on PerScope registration %s", ctorInfo, depType))
		}
		if err := cb.checkSingletonDependencies(depReg, stack); err != nil {
			return err
		}
	}
	return nil
}
