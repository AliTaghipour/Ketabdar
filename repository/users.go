package repository

import (
	"Ketab/model"
	"context"
	"errors"
	"sync"
)

type UsersRepository interface {
	AddUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id int32) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type UserRepositoryImpl struct {
	data map[int32]*model.User
	lock sync.Mutex

	currentId int32
	idLock    sync.Mutex
}

func NewUserRepositoryImpl() UsersRepository {
	data := make(map[int32]*model.User)
	return &UserRepositoryImpl{data: data}
}

func (u *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	u.idLock.Lock()
	defer u.idLock.Unlock()

	for _, v := range u.data {
		if v.Username == username {
			return v, nil
		}
	}
	return nil, errors.New("user not found")
}

func (u *UserRepositoryImpl) AddUser(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	id := u.getNextId()

	user.Id = id
	u.data[id] = user

	return nil
}

func (u *UserRepositoryImpl) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	user, ok := u.data[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserRepositoryImpl) getNextId() int32 {
	u.idLock.Lock()
	defer u.idLock.Unlock()

	id := u.currentId
	u.currentId++

	return id
}
