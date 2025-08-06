package main

import (
	"strings"
)

func profanityHandler(params *respParams) string {
	profanityCheck := false
	profanityWords := []string{"kerfuffle", "sharbert", "fornax"}
	profanitySet := make(map[string]struct{}, len(profanityWords))
	for _, word := range profanityWords {
		profanitySet[word] = struct{}{}
	}

	wordsList := strings.Split(params.Body, " ")
	for i, word := range wordsList {
		if _, ok := profanitySet[strings.ToLower(word)]; ok {
			wordsList[i] = "****"
			profanityCheck = true
		}
	}

	if profanityCheck {
		cleanText := strings.Join(wordsList, " ")
		return cleanText
	} else {
		return params.Body
	}
}
