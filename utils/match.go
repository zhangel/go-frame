package utils

import (
	"fmt"
	"net/http"
	"regexp"
)

func CreateHeaderRegexMap(toMatch map[string]string) (map[string]*regexp.Regexp, map[string]string, error) {
	m := make(map[string]*regexp.Regexp, len(toMatch))
	ms := make(map[string]string, len(toMatch))
	for k, v := range toMatch {
		regex, err := regexp.Compile(v)
		if err != nil {
			return nil, nil, err
		}
		m[http.CanonicalHeaderKey(k)] = regex
		ms[http.CanonicalHeaderKey(k)] = v
	}

	return m, ms, nil
}

func MatchHeadersWithRegex(toCheck map[string]*regexp.Regexp, toMatch map[string][]string) bool {
	if toCheck == nil {
		return true
	}

	for k, v := range toCheck {
		if values := toMatch[http.CanonicalHeaderKey(k)]; values == nil {
			return false
		} else if v != nil {
			valueExists := false
			for _, value := range values {
				if v.MatchString(value) {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}

func CreateHeaderMap(toMatch map[string]string) (map[string]string, error) {
	m := make(map[string]string, len(toMatch))
	for k, v := range toMatch {
		m[http.CanonicalHeaderKey(k)] = v
	}
	return m, nil
}

func MatchHeadersWithString(toCheck map[string]string, toMatch map[string][]string) bool {
	if toCheck == nil {
		return true
	}

	for k, v := range toCheck {
		if values := toMatch[http.CanonicalHeaderKey(k)]; values == nil {
			return false
		} else if v != "" {
			valueExists := false
			for _, value := range values {
				if v == value {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}

func checkPairs(pairs ...string) (int, error) {
	length := len(pairs)
	if length%2 != 0 {
		return length, fmt.Errorf("number of parameters must be multiple of 2, got %v", pairs)
	}
	return length, nil
}

func MapFromPairsToString(pairs ...string) (map[string]string, error) {
	length, err := checkPairs(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[http.CanonicalHeaderKey(pairs[i])] = pairs[i+1]
	}
	return m, nil
}
