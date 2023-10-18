package localization

import (
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var bundle *i18n.Bundle
var once sync.Once

var supportedLangs = []string{"en", "fr", "de", "ru", "es", "id", "it", "ja", "ko", "pt"}

func SupportedLangs(lang string) string {
	for _, supportedLang := range supportedLangs {
		if strings.EqualFold(lang, supportedLang) {
			return lang
		}
	}
	return "en"
}

func InitLocalization() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		path, err := filepath.Abs("./locales")
		if err != nil {
			log.Fatal(err)
		}

		localeFiles, err := filepath.Glob(path + "/*")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range localeFiles {
			_, err = bundle.LoadMessageFile(file)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func GetBundle() *i18n.Bundle {
	if bundle == nil {
		InitLocalization()
	}
	return bundle
}

func Localize(lang string, messageID string) string {
	lang = SupportedLangs(lang)

	bundle := GetBundle()

	localizer := i18n.NewLocalizer(bundle, lang)

	translatedMessage, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageID,
		},
	})
	if err != nil {
		log.Printf("Unable to localize string with ID %v: %v", messageID, err)
		return ""
	}

	return translatedMessage
}
