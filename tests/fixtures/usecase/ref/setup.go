package ref

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :conv CatDomainToModel Category
	DomainToModel(*domain.Pet) *model.Pet

	// :map ID CategoryID
	// :typecast
	CatDomainToModel(*domain.Category) model.Category
}
