package locale

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Locale struct {
	bundle *i18n.Bundle
}

var (
	once     sync.Once
	instance *Locale
)

func Init(rootPath string) *Locale {
	once.Do(func() {
		b := i18n.NewBundle(language.English)
		b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

		entries, err := os.ReadDir(rootPath)
		if err != nil {
			fmt.Printf("Error reading locale root directory %s: %v\n", rootPath, err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() {
					langDir := filepath.Join(rootPath, entry.Name())
					files, err := os.ReadDir(langDir)
					if err != nil {
						fmt.Printf("Error reading locale dir %s: %v\n", langDir, err)
						continue
					}
					for _, file := range files {
						if !file.IsDir() && strings.HasSuffix(file.Name(), ".toml") {
							filePath := filepath.Join(langDir, file.Name())
							fmt.Printf("Loading locale file: %s\n", filePath)
							if _, err := b.LoadMessageFile(filePath); err != nil {
								fmt.Printf("Failed to load %s: %v\n", filePath, err)
							}
						}
					}
				}
			}
		}

		instance = &Locale{bundle: b}
	})
	return instance
}

func (l *Locale) NewLocalizer(lang string) *i18n.Localizer {
	if l == nil || l.bundle == nil {
		b := i18n.NewBundle(language.English)
		b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		return i18n.NewLocalizer(b, "en")
	}
	return i18n.NewLocalizer(l.bundle, lang)
}

func (l *Locale) T(localizer *i18n.Localizer, messageID string, data map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		fmt.Printf("Error translating message %s: %v\n", messageID, err)
		return messageID
	}
	return msg
}
