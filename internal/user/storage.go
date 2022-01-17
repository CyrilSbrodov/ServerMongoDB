package user

type Storage interface {
	Create(user *User) (string, error)
	GetAll() (string, error)
	Get(id string) (string, error)
	Update(id string, age int) (string, error)
	Delete(u *User) (string, error)
	MakeFriends(id *ID) (string, error)
}
