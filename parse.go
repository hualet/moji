package moji

import (
	"strconv"
	"strings"
)

func parseWindLevel(windLevel string) (levelLow, levelHigh int) {
	tokens := strings.Split(windLevel, "-")

	if val, err := strconv.Atoi(tokens[0]); err == nil {
		levelLow = val
		levelHigh = levelLow
	}

	if len(tokens) == 2 {
		if val, err := strconv.Atoi(tokens[1]); err == nil {
			levelHigh = val
		}
	}

	return
}
