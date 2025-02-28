package validation

import (
	"fmt"
	"strconv"
)

func FetchPathID(pathsID map[string]string, specified ...string) (map[string]int, error) {
	if len(pathsID) == 0 {
		return nil, fmt.Errorf("nothing to parse")
	}
	parsedIDs := make(map[string]int, len(pathsID))
	for _, IDType := range specified {
		idString, ok := pathsID[IDType]
		if !ok {
			return nil, fmt.Errorf("%s didn't provide", IDType)
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		parsedIDs[IDType] = id
	}
	return parsedIDs, nil
}
