package handler

import "go-gin-simple-api/service"

type Handler struct {
	AuthHandler            AuthHandler
	BookHandler            BookHandler
	MediaHandler           MediaHandler
	BookStockHandler       BookStockHandler
	CustomerHandler        CustomerHandler
	ChargeHandler          ChargeHandler
	BookTransactionHandler BookTransactionHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		AuthHandler:            *NewAuthHandler(service.AuthService),
		BookHandler:            *NewBookHandler(service.BookService),
		MediaHandler:           *NewMediaHandler(service.MediaService),
		BookStockHandler:       *NewBookStockHandler(service.BookStockService),
		CustomerHandler:        *NewCustomerHandler(service.CustomerService),
		ChargeHandler:          *NewChargeHandler(service.ChargeService),
		BookTransactionHandler: *NewBookTransactionHandler(service.BookTransactionService),
	}
}
