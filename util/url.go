package util

import "github.com/google/go-querystring/query"

func UrlEncode(data any) (string, error) {
	v, err := query.Values(data)
	if err != nil {
		return "", err
	}
	return v.Encode(), nil
}
