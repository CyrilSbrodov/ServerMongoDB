package repository

//
//import (
//	"errors"
//	"fmt"
//	"sync"
//
//	"server/internal/user"
//)
//
//type storage struct {
//	sync.RWMutex
//	store map[int]*user.User
//}
//
//func NewStorage() user.Storage {
//	data := &storage{}
//	data.store = make(map[int]*user.User)
//	return &storage{
//		store: map[int]*user.User{},
//	}
//}
//
//func (s *storage) Create(user *user.User) (string, error) {
//	maxID := 0
//	for _, user := range s.store {
//		if user.Id > maxID {
//			maxID = user.Id
//		}
//	}
//	maxID++
//	user.Id = maxID
//	s.store[user.Id] = user
//	response := fmt.Sprintf("Пользователь %s создан\n", user.Name)
//	response += fmt.Sprintf("ID пользователя %s = %v", user.Name, user.Id)
//	return response, nil
//}
//
//func (s *storage) GetAll() (string, error) {
//	response := ""
//	for _, user := range s.store{
//		response += user.ToString()
//	}
//	return response, nil
//}
//
//func (s *storage) Get(id int) string {
//	response := ""
//	for _, user := range s.store {
//		if user.Id == id {
//			response += fmt.Sprintf("Друзья пользователя %s, ID пользователя %d \n", user.Name, user.Id)
//			for i := 0; i < len(user.Friends); i++ {
//				for _, j := range s.store {
//					if user.Friends[i] == j.Id {
//						response += fmt.Sprintf("Имя %s, возраст %d, ID пользователя %d, ID друзей %d \n", j.Name, j.Age, j.Id, j.Friends)
//					}
//				}
//			}
//		}
//	}
//	return response
//}
//
//func (s *storage) Update(id int, age int) (string, error) {
//
//	if _, exist := s.store[id]; !exist {
//		return "", errors.New(fmt.Sprintf("Пользователь с ID %v не найден.", id))
//	} else {
//		user := s.store[id]
//		user.Age = age
//
//		return fmt.Sprintf("Возраст пользователя %s изменен на %d.", user.Name, user.Age), nil
//	}
//}
//
//func (s *storage) Delete(id int) (string, error) {
//	//удаление пользователя из списка друзей
//	for _, user := range s.store {
//		for i := 0; i < len(user.Friends); i++ {
//			if user.Friends[i] == id {
//				user.Friends[i] = user.Friends[len(user.Friends)-1]
//				user.Friends = user.Friends[:len(user.Friends)-1]
//			}
//		}
//	}
//
//	//удаление пользователя
//	if _, exist := s.store[id]; !exist {
//		return "", errors.New(fmt.Sprintf("Пользователь с ID %v не найден.", id))
//	} else  {
//		delete(s.store, id)
//
//		return fmt.Sprintf("Пользователь с ID %d удален.", id), nil
//	}
//}
//
//func (s *storage) MakeFriends(id *user.ID) (string, error) {
//	if _, exits := s.store[id.Source_id]; !exits {
//		return "", errors.New(fmt.Sprintf("Пользователь с ID %v не найден.", id.Source_id))
//	} else if _, ok := s.store[id.Target_id]; !ok {
//		return "", errors.New(fmt.Sprintf("Пользователь с ID %v не найден.", id.Target_id))
//	} else {
//		userSource := s.store[id.Source_id]
//		userTarget := s.store[id.Target_id]
//		userSource.Friends = append(userSource.Friends, userTarget.Id)
//		userTarget.Friends = append(userTarget.Friends, userSource.Id)
//		return fmt.Sprintf("Пользователи %s и %s теперь друзья.", userSource.Name, userTarget.Name), nil
//	}
//}
//
