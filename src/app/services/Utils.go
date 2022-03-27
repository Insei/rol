package services

import (
	"regexp"
	"strings"
)

func toSnakeCase(entityStringFieldsNames *[]string) *[]string {
	snakeNames := &[]string{}
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	for i := 0; i < len(*entityStringFieldsNames); i++ {
		containPass := strings.Contains(strings.ToLower((*entityStringFieldsNames)[i]), "pass")
		containKey := strings.Contains(strings.ToLower((*entityStringFieldsNames)[i]), "key")
		if containPass || containKey {
			continue
		}
		snakeName := matchAllCap.ReplaceAllString((*entityStringFieldsNames)[i], "${1}_${2}")
		*snakeNames = append(*snakeNames, strings.ToLower(snakeName))
	}
	return snakeNames
}
