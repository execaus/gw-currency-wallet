package handler

import "gw-currency-wallet/internal/service"

type Handler struct {
	s *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		s: srv,
	}
}
