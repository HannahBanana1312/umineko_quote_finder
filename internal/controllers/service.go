package controllers

import (
	"umineko_quote/internal/og"
	"umineko_quote/internal/quote"
)

type Service struct {
	QuoteService     quote.Service
	OGImageGenerator *og.ImageGenerator
	HTMLContent      string
}

func NewService(quoteService quote.Service, ogGen *og.ImageGenerator, htmlContent string) Service {
	return Service{
		QuoteService:     quoteService,
		OGImageGenerator: ogGen,
		HTMLContent:      htmlContent,
	}
}

func (s *Service) GetAPIRoutes() []FSetupRoute {
	all := []FSetupRoute{}
	all = append(all, s.getAllSystemRoutes()...)
	all = append(all, s.getAllQuoteRoutes()...)
	all = append(all, s.getAllOGAPIRoutes()...)
	return all
}

func (s *Service) GetPageRoutes() []FSetupRoute {
	return s.getAllOGPageRoutes()
}
