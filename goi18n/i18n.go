package goi18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chaos-io/core/go/chaos/core"
	"github.com/chaos-io/core/go/logs"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	i18n2 "github.com/chaos-io/i18n"
)

type translator struct {
	localizers map[language.Tag]*i18n.Localizer
}

func NewTranslator(langDir string) (i18n2.ITranslator, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	if !core.IsExist(langDir) {
		return nil, fmt.Errorf("langDir %s does not exist", langDir)
	}

	localizers := make(map[language.Tag]*i18n.Localizer)
	if err := filepath.Walk(langDir, func(p string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		name := f.Name()
		if !strings.HasSuffix(name, ".yaml") {
			return nil
		}

		langStr := strings.TrimSuffix(name, ".yaml")
		tag, err := language.Parse(langStr)
		if err != nil {
			logs.Debugw("failed to parse language", "lang", langStr)
			return nil
		}

		langFile := filepath.Join(langDir, name)
		if _, err := bundle.LoadMessageFile(langFile); err != nil {
			return fmt.Errorf("failed to load language file %s: %w", langFile, err)
		}

		localizers[tag] = i18n.NewLocalizer(bundle, tag.String())
		return nil
	}); err != nil {
		return nil, err
	}

	return &translator{localizers: localizers}, nil
}

func (i *translator) Translate(ctx context.Context, key, lang string) (string, error) {
	langTag, err := language.Parse(lang)
	if err != nil {
		return "", fmt.Errorf("invalid language (%s) error: %w", lang, err)
	}

	localizer, ok := i.localizers[langTag]
	if !ok {
		return "", fmt.Errorf("language (%s) is not supported", lang)
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return "", fmt.Errorf("failed to localize, lang: %s, key: %s, error: %w", lang, key, err)
	}

	return msg, nil
}

func (i *translator) MustTranslate(ctx context.Context, key, lang string) string {
	msg, err := i.Translate(ctx, key, lang)
	if err != nil {
		logs.Debugw("failed to translate", "lang", lang, "key", key, "error", err)
		return ""
	}
	return msg
}
