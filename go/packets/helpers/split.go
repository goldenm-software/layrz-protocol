package helpers

import (
	"regexp"
	"strings"
)

var packetTag = regexp.MustCompile(`<P[A-Za-z]>`)

// Split splits a string into packets, this function exists due to some
// issues on TCP that sends multiple packages at the same time
func Split(data string) []string {
	data = strings.TrimRight(data, "\n\r")
	locs := packetTag.FindAllStringIndex(data, -1)
	if len(locs) == 0 {
		if t := strings.TrimSpace(data); t != "" {
			return []string{t}
		}
		return nil
	}

	result := make([]string, 0, len(locs))
	for i, loc := range locs {
		start := loc[0]
		var end int
		if i+1 < len(locs) {
			end = locs[i+1][0]
		} else {
			end = len(data)
		}
		if p := strings.TrimSpace(data[start:end]); p != "" {
			result = append(result, p)
		}
	}
	return result
}
