package quote

import (
	_ "embed"
	"math/rand/v2"
	"strings"

	"github.com/sahilm/fuzzy"
)

//go:embed data.txt
var dataFile string

type Service interface {
	Search(query string, limit int) []SearchResult
	GetByCharacter(characterID string, limit int) []ParsedQuote
	Random(characterID string) *ParsedQuote
	GetCharacters() map[string]string
}

type service struct {
	parser     Parser
	quotes     []ParsedQuote
	quoteTexts []string
}

type SearchResult struct {
	Quote ParsedQuote `json:"quote"`
	Score int         `json:"score"`
}

func NewService() Service {
	p := NewParser()
	lines := strings.Split(dataFile, "\n")
	quotes := p.ParseAll(lines)

	quoteTexts := make([]string, len(quotes))
	for i := 0; i < len(quotes); i++ {
		quoteTexts[i] = quotes[i].Text
	}

	return &service{
		parser:     p,
		quotes:     quotes,
		quoteTexts: quoteTexts,
	}
}

func (s *service) Search(query string, limit int) []SearchResult {
	matches := fuzzy.Find(query, s.quoteTexts)

	if limit <= 0 || limit > len(matches) {
		limit = len(matches)
	}

	results := make([]SearchResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = SearchResult{
			Quote: s.quotes[matches[i].Index],
			Score: matches[i].Score,
		}
	}

	return results
}

func (s *service) GetByCharacter(characterID string, limit int) []ParsedQuote {
	var results []ParsedQuote

	for i := 0; i < len(s.quotes); i++ {
		if s.quotes[i].CharacterID == characterID {
			results = append(results, s.quotes[i])
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results
}

func (s *service) Random(characterID string) *ParsedQuote {
	if len(s.quotes) == 0 {
		return nil
	}

	if characterID == "" {
		idx := rand.IntN(len(s.quotes))
		return &s.quotes[idx]
	}

	var filtered []ParsedQuote
	for i := 0; i < len(s.quotes); i++ {
		if s.quotes[i].CharacterID == characterID {
			filtered = append(filtered, s.quotes[i])
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	idx := rand.IntN(len(filtered))
	return &filtered[idx]
}

func (s *service) GetCharacters() map[string]string {
	return GetAllCharacters()
}
