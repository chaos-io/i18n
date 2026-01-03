package i18n

import (
	"context"
)

//go:generate mockgen -destination=mocks/i18n.go -package=mocks . ITranslator
type ITranslator interface {
	Translate(ctx context.Context, key, lang string) (string, error)
	MustTranslate(ctx context.Context, key, lang string) string
}
