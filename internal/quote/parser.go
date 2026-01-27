package quote

import (
	"html"
	"regexp"
	"strings"
)

type Parser interface {
	ParseLine(line string) *ParsedQuote
	ParseAll(lines []string) []ParsedQuote
}

type parser struct {
	dialogueLineRegex  *regexp.Regexp
	voiceMetaRegex     *regexp.Regexp
	characterNameRegex *regexp.Regexp
	colourTextRegex    *regexp.Regexp
	bracketRegex       *regexp.Regexp
	cleanupPatterns    []string
}

func NewParser() Parser {
	return &parser{
		dialogueLineRegex:  regexp.MustCompile(`^d2? \[lv`),
		voiceMetaRegex:     regexp.MustCompile(`\[lv 0\*"(\d+)"\*"(\d+)"\]`),
		characterNameRegex: regexp.MustCompile(`\{f:\d+:([^}]+)\}`),
		colourTextRegex:    regexp.MustCompile(`\{c:([A-Fa-f0-9]+):([^}]+)\}`),
		bracketRegex:       regexp.MustCompile(`\[[^\]]*\]`),
		cleanupPatterns: []string{
			"`[@]", "`[\\]", "`[|]", "`\"", "\"`",
			"[@]", "[\\]", "[|]",
		},
	}
}

type ParsedQuote struct {
	Text        string `json:"text"`
	TextHtml    string `json:"textHtml"`
	CharacterID string `json:"characterId"`
	Character   string `json:"character"`
	AudioID     string `json:"audioId"`
	Episode     int    `json:"episode"`
}

func (p *parser) ParseLine(line string) *ParsedQuote {
	if !p.dialogueLineRegex.MatchString(line) {
		return nil
	}

	matches := p.voiceMetaRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}

	characterID := matches[1]
	audioID := matches[2]
	episode := p.parseEpisodeFromAudioID(audioID)

	text, textHtml := p.extractText(line)
	if text == "" {
		return nil
	}

	return &ParsedQuote{
		Text:        text,
		TextHtml:    textHtml,
		CharacterID: characterID,
		Character:   GetCharacterName(characterID),
		AudioID:     audioID,
		Episode:     episode,
	}
}

func (p *parser) parseEpisodeFromAudioID(audioID string) int {
	if len(audioID) < 1 {
		return 0
	}
	ep := int(audioID[0] - '0')
	if ep >= 1 && ep <= 8 {
		return ep
	}
	return 0
}

func (p *parser) extractText(line string) (string, string) {
	text := line

	for i := 0; i < len(p.cleanupPatterns); i++ {
		text = strings.ReplaceAll(text, p.cleanupPatterns[i], "")
	}

	text = p.bracketRegex.ReplaceAllString(text, "")
	text = strings.TrimPrefix(text, "d2 ")
	text = strings.TrimPrefix(text, "d ")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.TrimSpace(text)

	escapedText := html.EscapeString(text)

	textHtml := p.colourTextRegex.ReplaceAllString(escapedText, `<span style="color:#$1">$2</span>`)
	textHtml = p.characterNameRegex.ReplaceAllString(textHtml, `<span class="quote-name">$1</span>`)

	plainText := p.colourTextRegex.ReplaceAllString(text, "$2")
	plainText = p.characterNameRegex.ReplaceAllString(plainText, "$1")

	return plainText, textHtml
}

func (p *parser) ParseAll(lines []string) []ParsedQuote {
	var quotes []ParsedQuote

	for i := 0; i < len(lines); i++ {
		parsed := p.ParseLine(lines[i])
		if parsed != nil && len(parsed.Text) > 10 {
			quotes = append(quotes, *parsed)
		}
	}

	return quotes
}
