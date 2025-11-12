package service

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/zalando/go-keyring"
)

type PassphraseService struct {
	Env              string
	EnvPassphraseKey string
}

func NewPassphraseService(env string, envPassphraseKeys ...string) *PassphraseService {
	envPassphraseKey := "LOCKIFY_PASSPHRASE"
	if len(envPassphraseKeys) > 0 {
		envPassphraseKey = envPassphraseKeys[0]
	}

	return &PassphraseService{Env: env, EnvPassphraseKey: envPassphraseKey}
}

func ClearAllPassphrases() error {
	return keyring.DeleteAll("lockify")
}

func (service *PassphraseService) GetPassphrase() string {
	pass := os.Getenv(service.EnvPassphraseKey)
	if pass != "" {
		return pass
	}

	pass, _ = service.getCachedPassphrase()
	if pass != "" {
		return pass
	}

	return service.getPassphraseFromUser()
}

func (service *PassphraseService) ClearPassphrase() error {
	return keyring.Delete("lockify", service.Env)
}

func (service *PassphraseService) getPassphraseFromUser() string {
	var pass string
	prompt := &survey.Password{Message: "Enter passphrase:"}
	// use bubble tea
	survey.AskOne(prompt, &pass)
	service.cachePassphrase(pass)

	return pass
}

func (service *PassphraseService) cachePassphrase(pass string) error {
	return keyring.Set("lockify", service.Env, pass)
}

func (service *PassphraseService) getCachedPassphrase() (string, error) {
	return keyring.Get("lockify", service.Env)
}
