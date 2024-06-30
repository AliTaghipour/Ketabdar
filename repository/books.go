package repository

import (
	"Ketab/model"
	"errors"
	"sync"
)

type BookRepository interface {
	AddNewBook(book model.Book) error
	GetBookById(bookId int32) (model.Book, error)
}

type BookRepositoryImpl struct {
	data map[int32]*model.Book
	lock sync.Mutex

	currentId int32
	idLock    sync.Mutex
}

func (u *BookRepositoryImpl) AddNewBook(book model.Book) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	if _, ok := u.data[book.Id]; ok {
		return errors.New("book already exists")
	}
	u.data[book.Id] = &book
	return nil
}

func (u *BookRepositoryImpl) GetBookById(bookId int32) (model.Book, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	book, ok := u.data[bookId]
	if !ok {
		return model.Book{}, errors.New("book not found")
	}
	return *book, nil
}

func (u *BookRepositoryImpl) getNextId() int32 {
	u.idLock.Lock()
	defer u.idLock.Unlock()

	id := u.currentId
	u.currentId++

	return id
}
