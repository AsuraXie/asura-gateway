package util

import (
	"errors"
	"strings"
)

func GetHostFromUrl(url string) string {
	var urlArray = strings.Split(url, "?")
	return urlArray[0]
}

func CheckProtocol(url string) error {
	if !strings.Contains(url, "http") {
		return errors.New("miss protocl http")
	}
	return nil
}
