package domain

import (
	"github.com/bookweb/structcopy-gen/tests/fixtures/usecase/typecast/enums"
)

type User struct {
	ID     int
	Name   string
	Status enums.Status
	Origin Origin
}

type Origin struct {
	Region string
}
