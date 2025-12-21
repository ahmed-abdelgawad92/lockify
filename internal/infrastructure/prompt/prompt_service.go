package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// Service implements PromptService for user input prompts.
type Service struct{}

// NewService creates a new Service instance.
func NewService() service.PromptService {
	return &Service{}
}

// GetUserInputForKeyAndValue prompts the user for a key and value, optionally hiding the value input.
func (p *Service) GetUserInputForKeyAndValue(isSecret bool) (key, value string, err error) {
	prompt := &survey.Input{Message: "Enter key:"}
	err = survey.AskOne(prompt, &key)
	if err != nil {
		return "", "", fmt.Errorf("failed to get key input: %w", err)
	}

	if isSecret {
		prompt := &survey.Password{Message: "Enter secret:"}
		err = survey.AskOne(prompt, &value)
		if err != nil {
			return "", "", fmt.Errorf("failed to get secret input: %w", err)
		}
	} else {
		prompt = &survey.Input{Message: "Enter value:"}
		err = survey.AskOne(prompt, &value)
		if err != nil {
			return "", "", fmt.Errorf("failed to get value input: %w", err)
		}
	}

	return key, value, nil
}

// GetPassphraseInput prompts the user for a passphrase with the given message.
func (p *Service) GetPassphraseInput(message string) (string, error) {
	var passphrase string
	prompt := &survey.Password{Message: message}
	err := survey.AskOne(prompt, &passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase input: %w", err)
	}
	return passphrase, nil
}
