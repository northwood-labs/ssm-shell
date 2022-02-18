package main

import (
	prompt "github.com/c-bata/go-prompt"
)

// CustomOptions are a standardized set of parameters that get passed uniformly to all prompt instances.
func CustomOptions() []prompt.Option {
	options := []prompt.Option{
		// Initial
		prompt.OptionInputTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Cyan),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),

		// Selected
		prompt.OptionSelectedSuggestionTextColor(prompt.Yellow),
		prompt.OptionSelectedSuggestionBGColor(prompt.Cyan),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),

		// Other
		prompt.OptionShowCompletionAtStart(),
		// prompt.OptionCompletionWordSeparator(emptyString),
	}

	return options
}
