package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Localizer struct {
	defaultLocale string
	messages      map[string]map[string]string
}

func New(defaultLocale string) (*Localizer, error) {
	messages := map[string]map[string]string{}

	for _, locale := range []string{"en", "ka"} {
		entries, err := loadLocale(locale)
		if err != nil {
			return nil, err
		}
		messages[locale] = entries
	}

	if defaultLocale != "ka" {
		defaultLocale = "en"
	}

	return &Localizer{
		defaultLocale: defaultLocale,
		messages:      messages,
	}, nil
}

func (l *Localizer) ResolveLocale(acceptLanguage string) string {
	header := strings.ToLower(strings.TrimSpace(acceptLanguage))
	if strings.HasPrefix(header, "ka") {
		return "ka"
	}
	if strings.HasPrefix(header, "en") {
		return "en"
	}

	return l.defaultLocale
}

func (l *Localizer) Msg(locale string, key string) string {
	if localeMessages, ok := l.messages[locale]; ok {
		if message, exists := localeMessages[key]; exists {
			return message
		}
	}

	if enMessages, ok := l.messages["en"]; ok {
		if message, exists := enMessages[key]; exists {
			return message
		}
	}

	return key
}

func loadLocale(locale string) (map[string]string, error) {
	path := filepath.Join("locales", locale+".json")

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(file, &result); err != nil {
		return nil, err
	}

	return result, nil
}
