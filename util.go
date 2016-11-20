package main

import (
	"errors"
	"regexp"
)

func ValidateUuid(v string) bool {
	match, err := regexp.MatchString(
		"[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}", v,
	)
	if err != nil {
		return false
	}
	return match
}

func ToInt(n interface{}) (int, error) {
	if n == nil {
		return 0, nil
	}
	switch n.(type) {
	default:
		return 0, errors.New("Unknown type")
	case float32:
		return int(n.(float32)), nil
	case float64:
		return int(n.(float64)), nil
	case uint:
		return int(n.(uint)), nil
	case uint8:
		return int(n.(uint8)), nil
	case uint16:
		return int(n.(uint16)), nil
	case uint32:
		return int(n.(uint32)), nil
	case int:
		return n.(int), nil
	case int8:
		return int(n.(int8)), nil
	case int16:
		return int(n.(int16)), nil
	case int32:
		return int(n.(int32)), nil
	}
}
