package example

import (
	"github.com/structcopy/structcopy-gen/examples/internal/example/dto"
	"github.com/structcopy/structcopy-gen/examples/internal/example/entity"
)

//go:generate structcopy-gen structcopy-gen.go
type StructCopyGen interface {
	// :match_field Email EMail
	// :match_method FullName FullName()
	// :conv LastName TestConvert
	// :conv Email TestConvert
	// :skip_field SkipField
	UserToUserDTO(src *entity.User) (dst *dto.UserDTO)

	// :match_field Email EMail
	// :skip_field SkipField
	UserToUserDTORaw(src entity.User) (dst dto.UserDTO)

	// :struct_conv UserToUserDTO
	UserSliceToUserDTOSlice(src []*entity.User) (dst []*dto.UserDTO)

	// :struct_conv UserToUserDTORaw
	UserSliceToUserDTOSliceRaw(src []entity.User) (dst []dto.UserDTO)
}
