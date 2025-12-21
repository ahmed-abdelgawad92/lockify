package service

type PromptService interface {
	GetUserInputForKeyAndValue(isSecret bool) (key, value string, err error)
	GetPassphraseInput(message string) (string, error)
}
