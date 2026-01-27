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
	Search(query string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse
	GetByCharacter(characterID string, limit int, offset int, episode int) CharacterResponse
	Random(characterID string, episode int) *ParsedQuote
	GetCharacters() map[string]string
}

type CharacterResponse struct {
	CharacterID string        `json:"characterId"`
	Character   string        `json:"character"`
	Quotes      []ParsedQuote `json:"quotes"`
	Total       int           `json:"total"`
	Limit       int           `json:"limit"`
	Offset      int           `json:"offset"`
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

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
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

func (s *service) Search(query string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}

	matchesFilter := func(q ParsedQuote) bool {
		if characterID != "" && q.CharacterID != characterID {
			return false
		}
		if episode > 0 && q.Episode != episode {
			return false
		}
		return true
	}

	if !forceFuzzy {
		queryLower := strings.ToLower(query)
		var exactMatches []SearchResult

		for i := 0; i < len(s.quotes); i++ {
			if strings.Contains(strings.ToLower(s.quoteTexts[i]), queryLower) {
				if matchesFilter(s.quotes[i]) {
					exactMatches = append(exactMatches, SearchResult{
						Quote: s.quotes[i],
						Score: 100,
					})
				}
			}
		}

		if len(exactMatches) > 0 {
			return paginateResults(exactMatches, limit, offset)
		}
	}

	matches := fuzzy.Find(query, s.quoteTexts)
	if len(matches) == 0 {
		return SearchResponse{
			Results: []SearchResult{},
			Total:   0,
			Limit:   limit,
			Offset:  offset,
		}
	}

	topScore := matches[0].Score
	relativeThreshold := topScore / 10
	minFuzzyScore := len(query) * 100

	var fuzzyResults []SearchResult
	for i := 0; i < len(matches); i++ {
		if matches[i].Score >= relativeThreshold && matches[i].Score >= minFuzzyScore {
			if matchesFilter(s.quotes[matches[i].Index]) {
				fuzzyResults = append(fuzzyResults, SearchResult{
					Quote: s.quotes[matches[i].Index],
					Score: matches[i].Score,
				})
			}
		}
	}

	return paginateResults(fuzzyResults, limit, offset)
}

func paginateResults(results []SearchResult, limit int, offset int) SearchResponse {
	total := len(results)

	if offset >= total {
		return SearchResponse{
			Results: []SearchResult{},
			Total:   total,
			Limit:   limit,
			Offset:  offset,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return SearchResponse{
		Results: results[offset:end],
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}
}

func (s *service) GetByCharacter(characterID string, limit int, offset int, episode int) CharacterResponse {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var all []ParsedQuote
	for i := 0; i < len(s.quotes); i++ {
		if s.quotes[i].CharacterID == characterID {
			if episode <= 0 || s.quotes[i].Episode == episode {
				all = append(all, s.quotes[i])
			}
		}
	}

	total := len(all)
	if offset >= total {
		return CharacterResponse{
			CharacterID: characterID,
			Character:   GetCharacterName(characterID),
			Quotes:      []ParsedQuote{},
			Total:       total,
			Limit:       limit,
			Offset:      offset,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return CharacterResponse{
		CharacterID: characterID,
		Character:   GetCharacterName(characterID),
		Quotes:      all[offset:end],
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}
}

func (s *service) Random(characterID string, episode int) *ParsedQuote {
	if len(s.quotes) == 0 {
		return nil
	}

	var filtered []ParsedQuote
	for i := 0; i < len(s.quotes); i++ {
		if characterID != "" && s.quotes[i].CharacterID != characterID {
			continue
		}
		if episode > 0 && s.quotes[i].Episode != episode {
			continue
		}
		filtered = append(filtered, s.quotes[i])
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
