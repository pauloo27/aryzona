package youtube

import (
	"errors"

	"github.com/Pauloo27/searchtube"
)

func GetBestResult(searchQuery string) (*searchtube.SearchResult, error) {
	results, err := searchtube.Search(searchQuery, 1)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	return results[0], nil
}
