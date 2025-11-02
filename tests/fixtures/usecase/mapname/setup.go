//go:build structcopygen
// +build structcopygen

package mapname

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :map Category.ID Category.CategoryID
	// :map Status.String() Status
	// :typecast
	DomainToModel(*domain.Pet) *model.Pet
}
