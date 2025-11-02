package stringer

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/stringer/local"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :stringer
	// :getter
	LocalToModel(pet *local.Pet) *model.Pet
}
