package simple

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	DomainToModel(*domain.Pet) *model.Pet
	ModelToDomain(*model.Pet) *domain.Pet
}
