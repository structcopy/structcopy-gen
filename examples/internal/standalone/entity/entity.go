package entity

type User struct {
	FirstName   string
	LastName    string
	EMail       string
	privateName string
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
