package controllers

import "umineko_quote/internal/quote"

type Service struct {
	QuoteService quote.Service
}

func NewService(quoteService quote.Service) Service {
	return Service{
		QuoteService: quoteService,
	}
}

func (s *Service) GetAllRoutes() []FSetupRoute {
	all := []FSetupRoute{}
	all = append(all, s.getAllSystemRoutes()...)
	all = append(all, s.getAllQuoteRoutes()...)
	return all
}
