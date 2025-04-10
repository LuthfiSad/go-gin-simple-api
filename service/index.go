package service

import (
	"go-gin-simple-api/lib"
	"go-gin-simple-api/repository"
)

type Service struct {
	AuthService            AuthService
	BookService            BookService
	MediaService           MediaService
	BookStockService       BookStockService
	CustomerService        CustomerService
	ChargeService          ChargeService
	BookTransactionService BookTransactionService
}

func NewService(repo *repository.Repository, cloudinary *lib.CloudinaryService) *Service {
	return &Service{
		AuthService:            NewAuthService(repo.AuthRepository),
		BookService:            NewBookService(repo.BookRepository, repo.MediaRepository),
		MediaService:           NewMediaService(repo.MediaRepository, repo.BookRepository, cloudinary),
		BookStockService:       NewBookStockService(repo.BookStockRepository, repo.BookRepository),
		CustomerService:        NewCustomerService(repo.CustomerRepository, repo.BookTransactionRepository),
		ChargeService:          NewChargeService(repo.ChargeRepository, repo.BookTransactionRepository, repo.AuthRepository),
		BookTransactionService: NewBookTransactionService(repo.BookTransactionRepository, repo.BookRepository, repo.BookStockRepository, repo.CustomerRepository),
	}
}
