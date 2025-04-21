package handler

import (
	"go-gin-simple-api/service"
	"sync"
)

type Handler struct {
	service *service.Service

	authHandler     *AuthHandler
	authHandlerOnce sync.Once

	bookHandler     *BookHandler
	bookHandlerOnce sync.Once

	mediaHandler     *MediaHandler
	mediaHandlerOnce sync.Once

	bookStockHandler     *BookStockHandler
	bookStockHandlerOnce sync.Once

	customerHandler     *CustomerHandler
	customerHandlerOnce sync.Once

	chargeHandler     *ChargeHandler
	chargeHandlerOnce sync.Once

	bookTransactionHandler     *BookTransactionHandler
	bookTransactionHandlerOnce sync.Once
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAuthHandler() *AuthHandler {
	h.authHandlerOnce.Do(func() {
		h.authHandler = NewAuthHandler(h.service.AuthService)
	})
	return h.authHandler
}

func (h *Handler) GetBookHandler() *BookHandler {
	h.bookHandlerOnce.Do(func() {
		h.bookHandler = NewBookHandler(h.service.BookService)
	})
	return h.bookHandler
}

func (h *Handler) GetMediaHandler() *MediaHandler {
	h.mediaHandlerOnce.Do(func() {
		h.mediaHandler = NewMediaHandler(h.service.MediaService)
	})
	return h.mediaHandler
}

func (h *Handler) GetBookStockHandler() *BookStockHandler {
	h.bookStockHandlerOnce.Do(func() {
		h.bookStockHandler = NewBookStockHandler(h.service.BookStockService)
	})
	return h.bookStockHandler
}

func (h *Handler) GetCustomerHandler() *CustomerHandler {
	h.customerHandlerOnce.Do(func() {
		h.customerHandler = NewCustomerHandler(h.service.CustomerService)
	})
	return h.customerHandler
}

func (h *Handler) GetChargeHandler() *ChargeHandler {
	h.chargeHandlerOnce.Do(func() {
		h.chargeHandler = NewChargeHandler(h.service.ChargeService)
	})
	return h.chargeHandler
}

func (h *Handler) GetBookTransactionHandler() *BookTransactionHandler {
	h.bookTransactionHandlerOnce.Do(func() {
		h.bookTransactionHandler = NewBookTransactionHandler(h.service.BookTransactionService)
	})
	return h.bookTransactionHandler
}

// package handler

// import "go-gin-simple-api/service"

// type Handler struct {
// 	AuthHandler            AuthHandler
// 	BookHandler            BookHandler
// 	MediaHandler           MediaHandler
// 	BookStockHandler       BookStockHandler
// 	CustomerHandler        CustomerHandler
// 	ChargeHandler          ChargeHandler
// 	BookTransactionHandler BookTransactionHandler
// }

// func NewHandler(service *service.Service) *Handler {
// 	return &Handler{
// 		AuthHandler:            *NewAuthHandler(service.AuthService),
// 		BookHandler:            *NewBookHandler(service.BookService),
// 		MediaHandler:           *NewMediaHandler(service.MediaService),
// 		BookStockHandler:       *NewBookStockHandler(service.BookStockService),
// 		CustomerHandler:        *NewCustomerHandler(service.CustomerService),
// 		ChargeHandler:          *NewChargeHandler(service.ChargeService),
// 		BookTransactionHandler: *NewBookTransactionHandler(service.BookTransactionService),
// 	}
// }
