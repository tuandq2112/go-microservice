package locale

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Locale struct {
	bundle *i18n.Bundle
}

// New creates a new Locale instance and loads all TOML files from the given path
func New(rootPath string) *Locale {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if err := loadLocaleFiles(b, rootPath); err != nil {
		fmt.Printf("Error loading locale files: %v\n", err)
	}

	return &Locale{bundle: b}
}

func loadLocaleFiles(b *i18n.Bundle, rootPath string) error {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("error reading locale root directory %s: %v", rootPath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		langDir := filepath.Join(rootPath, entry.Name())
		files, err := os.ReadDir(langDir)
		if err != nil {
			fmt.Printf("Error reading locale dir %s: %v\n", langDir, err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".toml") {
				filePath := filepath.Join(langDir, file.Name())
				if _, err := b.LoadMessageFile(filePath); err != nil {
					fmt.Printf("Failed to load %s: %v\n", filePath, err)
				}
			}
		}
	}
	return nil
}

// CreateLocalizer creates a new localizer for the given language
func (l *Locale) CreateLocalizer(lang string) *i18n.Localizer {
	if l == nil || l.bundle == nil {
		return nil
	}
	return i18n.NewLocalizer(l.bundle, lang)
}

func (l *Locale) Translate(localizer *i18n.Localizer, messageID string, data map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return msg
}
