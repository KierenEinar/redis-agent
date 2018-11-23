package commons

import "strings"

func IsBlank(s *string) bool {

	if s == nil || strings.TrimSpace(*s) == "" {
		return true
	}

	return false
}

func IsNotBlank(s *string) bool  {
	return !IsBlank(s)
}
