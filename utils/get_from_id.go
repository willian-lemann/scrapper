package utils

import (
	"regexp"
)

func GetIDFromLink(str string) string {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(str)
	return match
}
