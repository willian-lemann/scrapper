package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func CheckURLImage(url string) bool {
	if strings.Trim(url, " ") == "" {
		fmt.Println("URL is empty")
		return false
	}

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true

}
