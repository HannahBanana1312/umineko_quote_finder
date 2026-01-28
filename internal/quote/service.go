package quote

import (
	"embed"
	"math/rand/v2"
	"strings"

	"github.com/sahilm/fuzzy"
)

//go:embed data/*
var dataFS embed.FS

type Service interface {
	Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse
	GetByCharacter(lang string, characterID string, limit int, offset int, episode int) CharacterResponse
	GetByAudioID(lang string, audioID string) *ParsedQuote
	Random(lang string, characterID string, episode int) *ParsedQuote
	GetCharacters() map[string]string
}

type service struct {
	parser     Parser
	quotes     map[string][]ParsedQuote // "en" → English quotes, "ja" → Japanese quotes
	quoteTexts map[string][]string      // "en" → English texts for fuzzy search
}

func NewService() Service {
	p := NewParser()
	quotes := make(map[string][]ParsedQuote)
	texts := make(map[string][]string)

	langFiles := map[string]string{
		"en": "data/english.txt",
		"ja": "data/japanese.txt",
	}

	for lang, path := range langFiles {
		data, err := dataFS.ReadFile(path)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		parsed := p.ParseAll(lines)
		quotes[lang] = parsed

		quoteTexts := make([]string, len(parsed))
		for i := 0; i < len(parsed); i++ {
			quoteTexts[i] = parsed[i].Text
		}
		texts[lang] = quoteTexts
	}

	return &service{
		parser:     p,
		quotes:     quotes,
		quoteTexts: texts,
	}
}

func (s *service) Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	quoteTexts := s.quoteTexts[lang]
	if quotes == nil {
		return NewSearchResponse(nil, limit, offset)
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

		for i := 0; i < len(quotes); i++ {
			if strings.Contains(strings.ToLower(quoteTexts[i]), queryLower) {
				if matchesFilter(quotes[i]) {
					exactMatches = append(exactMatches, NewSearchResult(quotes[i], 100))
				}
			}
		}

		if len(exactMatches) > 0 {
			return NewSearchResponse(exactMatches, limit, offset)
		}
	}

	matches := fuzzy.Find(query, quoteTexts)
	if len(matches) == 0 {
		return NewSearchResponse(nil, limit, offset)
	}

	topScore := matches[0].Score
	relativeThreshold := topScore / 10
	minFuzzyScore := len(query) * 100

	var fuzzyResults []SearchResult
	for i := 0; i < len(matches); i++ {
		if matches[i].Score >= relativeThreshold && matches[i].Score >= minFuzzyScore {
			if matchesFilter(quotes[matches[i].Index]) {
				fuzzyResults = append(fuzzyResults, NewSearchResult(quotes[matches[i].Index], matches[i].Score))
			}
		}
	}

	return NewSearchResponse(fuzzyResults, limit, offset)
}

func (s *service) GetByCharacter(lang string, characterID string, limit int, offset int, episode int) CharacterResponse {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	var all []ParsedQuote
	for i := 0; i < len(quotes); i++ {
		if quotes[i].CharacterID == characterID {
			if episode <= 0 || quotes[i].Episode == episode {
				all = append(all, quotes[i])
			}
		}
	}

	return NewCharacterResponse(characterID, all, limit, offset)
}

func (s *service) Random(lang string, characterID string, episode int) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil || len(quotes) == 0 {
		return nil
	}

	var filtered []ParsedQuote
	for i := 0; i < len(quotes); i++ {
		if characterID != "" && quotes[i].CharacterID != characterID {
			continue
		}
		if episode > 0 && quotes[i].Episode != episode {
			continue
		}
		filtered = append(filtered, quotes[i])
	}

	if len(filtered) == 0 {
		return nil
	}

	idx := rand.IntN(len(filtered))
	return &filtered[idx]
}

func (s *service) GetByAudioID(lang string, audioID string) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return nil
	}

	for i := range quotes {
		if quotes[i].AudioID == audioID || strings.Contains(quotes[i].AudioID, audioID) {
			return &quotes[i]
		}
	}
	return nil
}

func (s *service) GetCharacters() map[string]string {
	return CharacterNames.GetAllCharacters()
}
