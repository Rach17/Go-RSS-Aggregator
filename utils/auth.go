package utils
import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	var val string = header.Get("Authorization")

	if val == "" {
		return "", errors.New("API key is missing in the request header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || strings.ToLower(vals[0]) != "apikey" {
		return "", errors.New("invalid API key format, expected 'Authorization: apikey <token>'")
	}

	return vals[1], nil
}