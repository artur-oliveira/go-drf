package i18n

import (
	"embed"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

type Service struct {
	Bundle *i18n.Bundle
}

func NewI18nService() *Service {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files, err := localeFS.ReadDir("locales")
	if err != nil {
		log.Fatalf("Falha ao ler o diretório de locales (embed): %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := "locales/" + file.Name()
		if _, err := bundle.LoadMessageFileFS(localeFS, path); err != nil {
			log.Printf("Aviso: Falha ao carregar arquivo de locale '%s': %v", path, err)
		}
	}

	return &Service{Bundle: bundle}
}
func (s *Service) GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(s.Bundle, lang)
}

func GetLocalizer(c *fiber.Ctx) *i18n.Localizer {
	loc := c.Locals("localizer")
	if loc == nil {
		return i18n.NewLocalizer(nil, "pt-BR")
	}
	return loc.(*i18n.Localizer)
}
