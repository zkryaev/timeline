package validation

import (
	"fmt"
	"strconv"
)

func FetchSpecifiedID(pathsID map[string]string, specified ...string) (map[string]int, error) {
	if len(pathsID) == 0 {
		return nil, fmt.Errorf("nothing to parse")
	}
	parsedIDs := make(map[string]int, len(pathsID))
	for _, IdType := range specified {
		idString, ok := pathsID[IdType]
		if !ok {
			return nil, fmt.Errorf("%s didn't provide", IdType)
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		parsedIDs[IdType] = id
	}
	return parsedIDs, nil
}
