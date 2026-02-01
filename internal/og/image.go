package og

import (
	"bytes"
	"embed"
	"image/color"
	"image/png"
	"sync"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed fonts/NotoSansJP-Regular.ttf
var notoSansJPData embed.FS

const (
	imgWidth  = 1200
	imgHeight = 630
)

var (
	bgColor        = color.RGBA{R: 10, G: 6, B: 18, A: 255}
	goldColor      = color.RGBA{R: 212, G: 168, B: 75, A: 255}
	creamColor     = color.RGBA{R: 232, G: 224, B: 240, A: 255}
	mutedColor     = color.RGBA{R: 168, G: 155, B: 184, A: 255}
	accentBarColor = color.RGBA{R: 212, G: 168, B: 75, A: 255}
)

type ImageGenerator struct {
	regularFont *sfnt.Font
	boldFont    *sfnt.Font
	jpFont      *sfnt.Font
	cache       sync.Map
}

func NewImageGenerator() *ImageGenerator {
	regular, _ := opentype.Parse(goregular.TTF)
	bold, _ := opentype.Parse(gobold.TTF)

	jpData, _ := notoSansJPData.ReadFile("fonts/NotoSansJP-Regular.ttf")
	jp, _ := opentype.Parse(jpData)

	return &ImageGenerator{
		regularFont: regular,
		boldFont:    bold,
		jpFont:      jp,
	}
}

func (g *ImageGenerator) textFont(lang string) *sfnt.Font {
	if lang == "ja" {
		return g.jpFont
	}
	return g.regularFont
}

func (g *ImageGenerator) boldOrFallback(lang string) *sfnt.Font {
	if lang == "ja" {
		return g.jpFont
	}
	return g.boldFont
}

func (g *ImageGenerator) Generate(audioId, lang, text, character string, episode int) ([]byte, error) {
	cacheKey := audioId + ":" + lang
	if cached, ok := g.cache.Load(cacheKey); ok {
		return cached.([]byte), nil
	}

	dc := gg.NewContext(imgWidth, imgHeight)

	dc.SetColor(bgColor)
	dc.Clear()

	dc.SetColor(accentBarColor)
	dc.DrawRectangle(0, 0, 8, float64(imgHeight))
	dc.Fill()

	quoteFace, err := opentype.NewFace(g.boldFont, &opentype.FaceOptions{Size: 72, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(quoteFace)
	dc.SetColor(color.RGBA{R: 212, G: 168, B: 75, A: 60})
	dc.DrawString("\u201C", 40, 100)

	textFace, err := opentype.NewFace(g.textFont(lang), &opentype.FaceOptions{Size: 28, DPI: 72})
	if err != nil {
		return nil, err
	}

	displayText := truncateText(text, 300)

	dc.SetFontFace(textFace)
	dc.SetColor(creamColor)
	dc.DrawStringWrapped(displayText, 60, 120, 0, 0, float64(imgWidth)-120, 1.5, gg.AlignLeft)

	charFace, err := opentype.NewFace(g.boldOrFallback(lang), &opentype.FaceOptions{Size: 24, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(charFace)
	dc.SetColor(goldColor)
	dc.DrawString("\u2014 "+character, 60, float64(imgHeight)-120)

	if episode > 0 {
		epFace, err := opentype.NewFace(g.regularFont, &opentype.FaceOptions{Size: 18, DPI: 72})
		if err != nil {
			return nil, err
		}
		dc.SetFontFace(epFace)
		dc.SetColor(mutedColor)
		dc.DrawString(g.episodeName(episode), 60, float64(imgHeight)-88)
	}

	brandFace, err := opentype.NewFace(g.regularFont, &opentype.FaceOptions{Size: 16, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(brandFace)
	dc.SetColor(mutedColor)
	dc.DrawStringAnchored("Umineko Quote Search", float64(imgWidth)-40, float64(imgHeight)-30, 1, 0)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	g.cache.Store(cacheKey, data)
	return data, nil
}

func truncateText(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes-3]) + "..."
	}
	return s
}

func (*ImageGenerator) episodeName(ep int) string {
	names := map[int]string{
		1: "Episode 1 \u2014 Legend",
		2: "Episode 2 \u2014 Turn",
		3: "Episode 3 \u2014 Banquet",
		4: "Episode 4 \u2014 Alliance",
		5: "Episode 5 \u2014 End",
		6: "Episode 6 \u2014 Dawn",
		7: "Episode 7 \u2014 Requiem",
		8: "Episode 8 \u2014 Twilight",
	}
	if name, ok := names[ep]; ok {
		return name
	}
	return ""
}
