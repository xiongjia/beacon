package strutil

import (
	"unicode"
)

func IsEmpty[T ~string](s T) bool {
	return s == ""
}

func IsEmptyPtr[T ~string](s *T) bool {
	if s == nil {
		return true
	}
	return IsEmpty(*s)
}

func IsBlank[T ~string](s T) bool {
	if IsEmpty(s) {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func IsBlankPtr[T ~string](s *T) bool {
	if s == nil {
		return true
	}
	return IsBlank(*s)
}
