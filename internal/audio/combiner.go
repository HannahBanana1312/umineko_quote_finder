package audio

import (
	"fmt"
	"os"
)

type FilePathResolver func(charId, audioId string) string

func CombineOgg(charId string, ids []string, resolve FilePathResolver) ([]byte, error) {
	var combined []byte
	for i := 0; i < len(ids); i++ {
		filePath := resolve(charId, ids[i])
		if filePath == "" {
			return nil, fmt.Errorf("audio file not found: %s", ids[i])
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %s", ids[i])
		}
		combined = append(combined, data...)
	}
	return combined, nil
}
