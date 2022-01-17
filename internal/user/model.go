package user

import "fmt"

type User struct {
	Id string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age int `json:"age" bson:"age"`
	Friends []string `json:"friends" bson:"friends,omitempty"`
}

type CreateUserDTO struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Friends []string `json:"friends"`
}

type ID struct {
	Source_id string `json:"source_id"`
	Target_id string `json:"target_id"`
}

type UpdateUser struct {
	Id int `json:"new id"`
	Name string `json:"new name"`
	New_age int `json:"new_age"`
}

func (u *User) ToString()string  {
	return fmt.Sprintf("Имя: %s. Возраст: %d. ID пользователя: %s. ID друзей: %s \n",
		u.Name, u.Age, u.Id, u.Friends)
}
