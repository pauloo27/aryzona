package arkw

type User struct {
	id string
}

func (u User) ID() string {
	return u.id
}

func buildUser(id string) User {
	return User{
		id: id,
	}
}
