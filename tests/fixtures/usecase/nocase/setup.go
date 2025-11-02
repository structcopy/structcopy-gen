//go:build structcopygen

package nocase

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/nocase/model"
)

type ModelA struct {
	ID   uint64
	Name string
}

func (a *ModelA) name() string {
	return a.Name
}

type ModelB struct {
	id   uint64
	name string
}

//go:generate structcopy-gen
type StructCopyGen interface {
	// :case:off
	// :getter
	// AtoB demonstrates local to local copy with case-insensitive field matching.
	// It shows that a private getter precedence over its (exported) counterpart field.
	AtoB(*ModelA) *ModelB
	// :case:off
	BtoA(*ModelB) *ModelA
	// :case:off
	// UserToB demonstrates copy an external package type to internal.
	// It skips private fields (and getters) in the former type.
	UserToB(*model.User) *ModelB
	// :case:off
	// BtoUser demonstrates copy an internal to external package type.
	// It skips private fields (and getters) in the latter type.
	BtoUser(*ModelB) *model.User
}
