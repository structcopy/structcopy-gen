package example

import (
	"github.com/bookweb/structcopy-gen/examples/internal/example/dto"
	"github.com/bookweb/structcopy-gen/examples/internal/example/entity"
)

//go:generate structcopy-gen structcopy-gen.go
type StructCopyGen interface {
	// :match_field Email EMail
	// :match_method FullName FullName()
	// :conv LastName TestConvert
	// :conv Email TestConvert
	UserToUserDTO(src *entity.User) (dst *dto.UserDTO)

	// :match_field Email EMail
	UserToUserDTORaw(src entity.User) (dst dto.UserDTO)

	// :struct_conv UserToUserDTO
	UserSliceToUserDTOSlice(src []*entity.User) (dst []*dto.UserDTO)

	// :struct_conv UserToUserDTORaw
	UserSliceToUserDTOSliceRaw(src []entity.User) (dst []dto.UserDTO)
}
