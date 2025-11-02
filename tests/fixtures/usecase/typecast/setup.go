package typecast

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/typecast/domain"
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/typecast/model"
)

//go:generate structcopy-gen
type StructCopyGen interface {
	// :typecast
	// DomainToModel converts domain.User to model.User.
	// typecast works:
	// - int64 -> int
	// - enums.Status -> string
	DomainToModel(*domain.User) *model.User

	// :typecast
	// ModelToDomain converts model.User to domain.User.
	// typecast works:
	// - int -> int64
	// - string -> enums.Status
	//   "enums" package will be imported automatically in the generated code!
	ModelToDomain(*model.User) *domain.User
}
