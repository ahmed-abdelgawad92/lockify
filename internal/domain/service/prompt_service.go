package service

type PromptService interface {
	GetUserInputForKeyAndValue(isSecret bool) (key, value string)
	GetPassphraseInput(message string) string
}
