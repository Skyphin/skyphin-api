package utils

import (
	"strings"
)

type ProfanityFilter struct {
	bannedWords map[string]bool
}

func NewProfanityFilter(bannedWords []string) *ProfanityFilter {
	filter := &ProfanityFilter{
		bannedWords: make(map[string]bool),
	}
	for _, word := range bannedWords {
		filter.bannedWords[strings.ToLower(word)] = true
	}
	return filter
}

func (f *ProfanityFilter) HasProfanity(text string) bool {
	words := strings.Fields(strings.ToLower(text))
	for _, word := range words {
		if f.bannedWords[word] {
			return true
		}
	}
	return false
}
