package interceptors

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tuandq2112/go-microservices/shared/locale"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	localizerKey  contextKey = "localizer"
	translatorKey contextKey = "translator"
)

func UnaryLocaleInterceptor(l *locale.Locale) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		acceptLang := "en"
		if lang := md.Get("locale"); len(lang) > 0 {
			acceptLang = lang[0]
		} else if lang := md.Get("grpcgateway-accept-language"); len(lang) > 0 {
			acceptLang = lang[0]
		}

		localizer := l.NewLocalizer(acceptLang)
		newCtx := context.WithValue(ctx, localizerKey, localizer)

		translatorFunc := func(messageID string, data map[string]interface{}) string {
			return l.T(localizer, messageID, data)
		}
		newCtx = context.WithValue(newCtx, translatorKey, translatorFunc)

		return handler(newCtx, req)
	}
}

func StreamLocaleInterceptor(l *locale.Locale) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}

		acceptLang := "en" // default to English
		if lang := md.Get("locale"); len(lang) > 0 {
			acceptLang = lang[0]
		} else if lang := md.Get("accept-language"); len(lang) > 0 {
			acceptLang = lang[0]
		}

		localizer := l.NewLocalizer(acceptLang)

		translatorFunc := func(messageID string, data map[string]interface{}) string {
			return l.T(localizer, messageID, data)
		}

		newCtx := context.WithValue(ss.Context(), localizerKey, localizer)
		newCtx = context.WithValue(newCtx, translatorKey, translatorFunc)

		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrappedStream)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// GetLocalizer retrieves the localizer from the context
func GetLocalizer(ctx context.Context) *i18n.Localizer {
	if localizer, ok := ctx.Value(localizerKey).(*i18n.Localizer); ok {
		return localizer
	}
	return nil
}

// T retrieves the translator function from the context and uses it to translate a message
func T(ctx context.Context, messageID string, data map[string]interface{}) string {
	if translator, ok := ctx.Value(translatorKey).(func(string, map[string]interface{}) string); ok {
		return translator(messageID, data)
	}
	return messageID
}
