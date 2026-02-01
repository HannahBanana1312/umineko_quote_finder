package controllers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) getAllOGAPIRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupOGImageRoute,
	}
}

func (s *Service) getAllOGPageRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupOGPageRoute,
	}
}

func (s *Service) setupOGImageRoute(routeGroup fiber.Router) {
	routeGroup.Get("/og/:audioId.png", s.ogImage)
}

func (s *Service) ogImage(ctx *fiber.Ctx) error {
	audioId := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	lang := ctx.Query("lang", "en")

	q := s.QuoteService.GetByAudioID(lang, audioId)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}

	data, err := s.OGImageGenerator.Generate(audioId, lang, q.Text, q.Character, q.Episode)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate image",
		})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func (s *Service) setupOGPageRoute(routeGroup fiber.Router) {
	routeGroup.Get("/", s.ogPage)
}

func (s *Service) ogPage(ctx *fiber.Ctx) error {
	audioId := ctx.Query("quote")
	if audioId == "" {
		return ctx.Next()
	}

	lang := ctx.Query("lang", "en")

	q := s.QuoteService.GetByAudioID(lang, audioId)
	if q == nil {
		return ctx.Next()
	}

	scheme := "https"
	if strings.HasPrefix(ctx.Hostname(), "localhost") || strings.HasPrefix(ctx.Hostname(), "127.0.0.1") {
		scheme = "http"
	}
	proto := ctx.Get("X-Forwarded-Proto")
	if proto != "" {
		scheme = proto
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, ctx.Hostname())

	title := fmt.Sprintf("%s \u2014 Umineko Quote", q.Character)
	description := q.Text
	if len(description) > 200 {
		description = description[:197] + "..."
	}
	imageURL := fmt.Sprintf("%s/api/v1/og/%s.png?lang=%s", baseURL, audioId, lang)

	html := s.HTMLContent
	html = strings.Replace(html, `<meta property="og:title" content="Umineko Quote Search">`, fmt.Sprintf(`<meta property="og:title" content="%s">`, escapeAttr(title)), 1)
	html = strings.Replace(html, `<meta property="og:description" content="Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni. When the seagulls cry, none shall remain.">`, fmt.Sprintf(`<meta property="og:description" content="%s">`, escapeAttr(description)), 1)
	html = strings.Replace(html, `<meta property="og:image" content="https://waifuvault.moe/f/5e9cf90a-8a63-48b3-802d-1bc9be9062ea/clipboard-image-1769601762638.png">`, fmt.Sprintf(`<meta property="og:image" content="%s">`, imageURL), 1)
	html = strings.Replace(html, `<meta name="twitter:title" content="Umineko Quote Search">`, fmt.Sprintf(`<meta name="twitter:title" content="%s">`, escapeAttr(title)), 1)
	html = strings.Replace(html, `<meta name="twitter:description" content="Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni.">`, fmt.Sprintf(`<meta name="twitter:description" content="%s">`, escapeAttr(description)), 1)
	html = strings.Replace(html, `<meta name="twitter:image" content="https://waifuvault.moe/f/5e9cf90a-8a63-48b3-802d-1bc9be9062ea/clipboard-image-1769601762638.png">`, fmt.Sprintf(`<meta name="twitter:image" content="%s">`, imageURL), 1)

	ctx.Set("Content-Type", "text/html; charset=utf-8")
	return ctx.SendString(html)
}

func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
