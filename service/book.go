package service

import (
	"fmt"

	"github.com/train-do/Golang-Restfull-API/model"
	"github.com/train-do/Golang-Restfull-API/repository"
	"go.uber.org/zap"
)

type BookService struct {
	repo *repository.BookRepository
	log  *zap.Logger
}

func NewBookService(repo *repository.BookRepository,
	log *zap.Logger) *BookService {
	return &BookService{repo, log}
}

func (s *BookService) GetAllBook() ([]model.Book, error) {
	return s.repo.GetAll()
}

// func (s *BookService) GetProductByID(id int) (*model.Book, error) {
// 	return s.repo.GetByID(id)
// }

func (s *BookService) CreateBook(book *model.Book) error {
	fmt.Println("MASUK SERVICE CREATE BOOK")
	return s.repo.Create(book)
}

// func (s *BookService) UpdateProduct(product *model.Book) error {
// 	return s.repo.Update(product)
// }

// func (s *BookService) DeleteProduct(id int) error {
// 	return s.repo.Delete(id)
// }
