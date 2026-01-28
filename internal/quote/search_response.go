package quote

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
}

func NewSearchResponse(results []SearchResult, limit int, offset int) SearchResponse {
	if results == nil {
		results = []SearchResult{}
	}

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
