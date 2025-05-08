package Group

import (
	gi "github.com/TaBSRest/GoFac/internal/RegistrationOption/GroupInfo"
	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type Group struct {
	*gi.GroupInfo
	Registrations []*r.Registration
}
