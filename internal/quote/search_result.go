package quote

type SearchResult struct {
	Quote ParsedQuote `json:"quote"`
	Score int         `json:"score"`
}

func NewSearchResult(quote ParsedQuote, score int) SearchResult {
	return SearchResult{
		Quote: quote,
		Score: score,
	}
}
