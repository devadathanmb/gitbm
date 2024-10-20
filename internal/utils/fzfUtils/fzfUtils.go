package fzfutils

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
)

// ErrSelectionCancelled is returned when the user cancels the fuzzy selection.
var ErrSelectionCancelled = fmt.Errorf("selection cancelled")

// FuzzyFind presents a list of items to the user for fuzzy selection.
// It returns the selected item and any error encountered.
func FuzzyFind[T any](items []T, displayFunc func(T) string, promptString string) (T, error) {
	var zero T
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return displayFunc(items[i])
		},
		fuzzyfinder.WithPromptString(promptString),
	)
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return zero, ErrSelectionCancelled
		}
		return zero, err
	}
	return items[idx], nil
}
