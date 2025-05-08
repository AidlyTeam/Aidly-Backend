package domains

type User struct {
	Name    string
	Surname string
}

func ToUser(name, surname string) *User {
	return &User{
		Name:    name,
		Surname: surname,
	}
}
