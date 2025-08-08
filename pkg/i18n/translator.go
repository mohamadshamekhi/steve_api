package i18n

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Translator struct {
	bundle   *i18n.Bundle
	fallback language.Tag
}

func New(dirPath string, fallbackLang string) (*Translator, error) {
	fallback := language.MustParse(fallbackLang)
	bundle := i18n.NewBundle(fallback)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not read i18n dir: %w", err)
	}

	for _, file := range files {
		if _, err := bundle.LoadMessageFile(fmt.Sprintf("%s/%s", dirPath, file.Name())); err != nil {
			return nil, fmt.Errorf("failed to load message file %s: %w", file.Name(), err)
		}
	}

	return &Translator{bundle: bundle, fallback: fallback}, nil
}

func (t *Translator) Translate(lang, key string, data map[string]any) string {
	localizer := i18n.NewLocalizer(t.bundle, lang)
	localizedMsg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})

	if err != nil {
		return key
	}

	return localizedMsg
}
