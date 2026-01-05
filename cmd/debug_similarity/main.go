package main

import (
	"fmt"
	"strings"
)

func main() {
	// Expose the private methods via a temporary public wrapper or just copy-paste for the test.
	// Since jaccardSimilarity is private in RankingService, we'll implement a local version here to match the logic exactly.

	titles := []string{
		"Delcy Rodríguez es investida como presidenta encargada de Venezuela",
		"Delcy Rodríguez jura como presidenta encargada de Venezuela tras la captura de Nicolás Maduro: 'Vengo con dolor'",
		"Delcy Rodríguez jura como presidenta encargada de Venezuela",
		"Delcy Rodríguez asume como Presidenta interina de Venezuela tras la captura de Nicolás Maduro",
	}

	fmt.Println("🔹 Analyzing Similarity Matrix:")

	for i := 0; i < len(titles); i++ {
		for j := i + 1; j < len(titles); j++ {
			s1 := titles[i]
			s2 := titles[j]
			score := jaccardSimilarity(s1, s2)
			fmt.Printf("A: %-30.30s... | B: %-30.30s... | Jaccard: %.4f\n", s1, s2, score)
		}
	}
}

// COPIED FROM internal/fetcher/ranking.go
func jaccardSimilarity(s1, s2 string) float64 {
	set1 := tokenize(s1)
	set2 := tokenize(s2)

	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}

	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	return float64(intersection) / float64(union)
}

func tokenize(text string) map[string]bool {
	tokens := make(map[string]bool)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		// Basic cleanup: remove punctuation
		w = strings.TrimFunc(w, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r > 127) // Keep latin supplement
		})
		if len(w) > 2 {
			tokens[w] = true
		}
	}
	return tokens
}
