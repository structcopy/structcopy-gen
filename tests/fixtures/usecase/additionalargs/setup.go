package converter

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model"
	"github.com/bookweb/structcopy-gen/tests/fixtures/data/model/abc222"
)

//go:generate structcopy-gen

type StructCopyGen interface {
	//:map $2 List
	DomainToModel(*model.Additional, []abc222.Additional321) *abc222.AdditionalItem123
}
