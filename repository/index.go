package repository

import "gorm.io/gorm"

type Repository struct {
	AuthRepository            AuthRepository
	BookRepository            BookRepository
	MediaRepository           MediaRepository
	BookStockRepository       BookStockRepository
	CustomerRepository        CustomerRepository
	ChargeRepository          ChargeRepository
	BookTransactionRepository BookTransactionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepository:            NewAuthRepository(db),
		BookRepository:            NewBookRepository(db),
		MediaRepository:           NewMediaRepository(db),
		BookStockRepository:       NewBookStockRepository(db),
		CustomerRepository:        NewCustomerRepository(db),
		ChargeRepository:          NewChargeRepository(db),
		BookTransactionRepository: NewBookTransactionRepository(db),
	}
}
