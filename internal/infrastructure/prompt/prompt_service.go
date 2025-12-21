package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type PromptService struct{}

func NewPromptService() service.PromptService {
	return &PromptService{}
}

func (p *PromptService) GetUserInputForKeyAndValue(isSecret bool) (key, value string, err error) {
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

func (p *PromptService) GetPassphraseInput(message string) (string, error) {
	var passphrase string
	prompt := &survey.Password{Message: message}
	err := survey.AskOne(prompt, &passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase input: %w", err)
	}
	return passphrase, nil
}
