package quote

import (
	"html"
	"regexp"
	"strings"
)

type Parser interface {
	ParseAll(lines []string) []ParsedQuote
}

type textRule struct {
	pattern   *regexp.Regexp
	htmlRepl  string
	plainRepl string
}

type parser struct {
	dialogueLineRegex *regexp.Regexp
	voiceMetaRegex    *regexp.Regexp
	bracketRegex      *regexp.Regexp
	cleanupPatterns   []string
	textRules         []textRule
}

func NewParser() Parser {
	return &parser{
		dialogueLineRegex: regexp.MustCompile(`^d2? \[lv`),
		voiceMetaRegex:    regexp.MustCompile(`\[lv 0\*"(\d+)"\*"(\d+)"\]`),
		bracketRegex:      regexp.MustCompile(`\[[^\]]*\]`),
		cleanupPatterns: []string{
			"`[@]", "`[\\]", "`[|]", "`\"", "\"`",
			"[@]", "[\\]", "[|]",
		},
		textRules: []textRule{
			{regexp.MustCompile(`\{c:([A-Fa-f0-9]+):([^}]+)\}`), `<span style="color:#$1">$2</span>`, "$2"},
			{regexp.MustCompile(`\{f:\d+:([^}]+)\}`), `<span class="quote-name">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:\d{2,}:([^}]+)\}`), `<span class="quote-name">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:1:([^}]+)\}`), `<span class="red-truth">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:2:([^}]+)\}`), `<span class="blue-truth">$1</span>`, "$1"},
			{regexp.MustCompile(`\{ruby:([^:]+):([^}]+)\}`), `<ruby>$2<rp>(</rp><rt>$1</rt><rp>)</rp></ruby>`, "$2 ($1)"},
			{regexp.MustCompile(`\{i:([^}]+)\}`), `<em>$1</em>`, "$1"},
			{regexp.MustCompile(`\{[a-z]+:[^}]*\}`), "", ""},
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

func (p *parser) parseLine(line string) *ParsedQuote {
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

	for _, pattern := range p.cleanupPatterns {
		text = strings.ReplaceAll(text, pattern, "")
	}

	text = p.bracketRegex.ReplaceAllString(text, "")
	text = strings.TrimPrefix(text, "d2 ")
	text = strings.TrimPrefix(text, "d ")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.TrimSpace(text)

	plainText := text
	textHtml := html.EscapeString(text)

	for _, rule := range p.textRules {
		textHtml = rule.pattern.ReplaceAllString(textHtml, rule.htmlRepl)
		plainText = rule.pattern.ReplaceAllString(plainText, rule.plainRepl)
	}

	return plainText, textHtml
}

func (p *parser) ParseAll(lines []string) []ParsedQuote {
	var quotes []ParsedQuote

	for i := 0; i < len(lines); i++ {
		parsed := p.parseLine(lines[i])
		if parsed != nil && len(parsed.Text) > 10 {
			quotes = append(quotes, *parsed)
		}
	}

	return quotes
}
