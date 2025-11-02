package standalone

import (
	"github.com/bookweb/structcopy-gen/examples/internal/standalone/dto"
	"github.com/bookweb/structcopy-gen/examples/internal/standalone/entity"
)

// :structcopygen
//
//go:generate structcopy-gen structcopy-gen.go
type StructCopyGen interface {

	// :match_field Email EMail
	// :match_method FullName FullName()
	// :conv LastName TestConvert
	// :conv Email TestConvert
	UserToUserDTO(src *entity.User) (dst *dto.UserDTO)

	UserToUserDTORaw(src entity.User) (dst dto.UserDTO)

	// :struct_conv UserToUserDTO
	UserSliceToUserDTOSlice(src []*entity.User) (dst []*dto.UserDTO)

	// :struct_conv UserToUserDTORaw
	UserSliceToUserDTOSliceRaw(src []entity.User) (dst []dto.UserDTO)

	TestToTestDTO(src *Test) (dst *TestDTO)

	TestToTestDTORaw(src Test) (dst TestDTO)
}
