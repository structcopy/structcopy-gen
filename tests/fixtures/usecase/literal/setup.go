//go:build structcopygen

package literal

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :literal  Name   "abc  def"
	DomainToModel(*domain.Pet) *model.Pet
	ModelToDomain(*model.Pet) *domain.Pet
}
