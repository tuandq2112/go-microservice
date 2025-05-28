package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tuandq2112/go-microservices/shared/locale"
	"golang.org/x/text/language"
)

type contextKey string

const (
	localizerKey  contextKey = "localizer"
	translatorKey contextKey = "translator"
)

func LocaleMiddleware(l *locale.Locale) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptLang := r.Header.Get("Accept-Language")
			if acceptLang == "" {
				acceptLang = "en"
			}
			matcher := language.NewMatcher([]language.Tag{
				language.English,
				language.Vietnamese,
			})

			tag, _ := language.MatchStrings(matcher, acceptLang)
			lang := strings.Split(tag.String(), "-")[0]

			localizer := l.NewLocalizer(lang)

			ctx := context.WithValue(r.Context(), localizerKey, localizer)
			ctx = context.WithValue(ctx, translatorKey, func(messageID string, data map[string]interface{}) string {
				return l.T(localizer, messageID, data)
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLocalizer(r *http.Request) *i18n.Localizer {
	if localizer, ok := r.Context().Value(localizerKey).(*i18n.Localizer); ok {
		return localizer
	}
	return nil
}

func T(r *http.Request, messageID string, data map[string]interface{}) string {
	if translator, ok := r.Context().Value(translatorKey).(func(string, map[string]interface{}) string); ok {
		return translator(messageID, data)
	}
	return messageID
}
