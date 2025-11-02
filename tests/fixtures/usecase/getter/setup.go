//go:build structcopygen
// +build structcopygen

package getter

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/ddd/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/ddd/model"
)

// :getter:off
//
//go:generate structcopy-gen
type StructCopyGen interface {
	// DomainToModel copies domain.Pet to model.Pet.
	// :skip PhotoUrls
	// :getter
	DomainToModel(pet *domain.Pet) *model.Pet

	// DomainToModelNoGetter copies domain.Pet to model.Pet but not using getters.
	// :getter
	DomainToModelNoGetter(pet *domain.Pet) *model.Pet
}
