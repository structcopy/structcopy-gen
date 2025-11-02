//go:build structcopygen
// +build structcopygen

package converter

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/embedded/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/embedded/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :getter
	// :typecast
	DomainToModel(s *domain.Concrete) (d *model.Concrete)
	// :getter
	// :typecast
	ModelToDomain(*model.Concrete) (*domain.Concrete, error)
}
