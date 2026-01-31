package quote

type CharacterResponse struct {
	CharacterID string        `json:"characterId"`
	Character   string        `json:"character"`
	Quotes      []ParsedQuote `json:"quotes"`
	Total       int           `json:"total"`
	Limit       int           `json:"limit"`
	Offset      int           `json:"offset"`
}

func NewCharacterResponse(characterID string, quotes []ParsedQuote, limit int, offset int) CharacterResponse {
	if quotes == nil {
		quotes = []ParsedQuote{}
	}

	var characterName string
	if characterID != "" {
		characterName = CharacterNames.GetCharacterName(characterID)
	}

	total := len(quotes)

	if offset >= total {
		return CharacterResponse{
			CharacterID: characterID,
			Character:   characterName,
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
		Character:   characterName,
		Quotes:      quotes[offset:end],
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}
}
