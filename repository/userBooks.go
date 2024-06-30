package repository

import (
	"Ketab/model"
	"sync"
)

type UserBooksRepository interface {
	AddNewUserBook(userBook *model.UserBook, userId int32) error
	GetUserBooks(userId int32) ([]*model.UserBook, error)
}

type UserBookRepoImpl struct {
	data map[int32][]*model.UserBook
	lock sync.Mutex
}

func (u *UserBookRepoImpl) AddNewUserBook(userBook *model.UserBook, userId int32) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	list, ok := u.data[userId]
	if !ok {
		list = []*model.UserBook{}
		u.data[userId] = list
	}
	list = append(list, userBook)
	u.data[userId] = list
	return nil
}

func (u *UserBookRepoImpl) GetUserBooks(userId int32) ([]*model.UserBook, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	list, ok := u.data[userId]
	if !ok {
		return []*model.UserBook{}, nil
	}
	return list, nil
}
