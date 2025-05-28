package interceptors

import (
	"context"

	"github.com/tuandq2112/go-microservices/shared/locale"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	defaultLang = "en"
)

func getLanguage(md metadata.MD) string {
	if lang := md.Get("locale"); len(lang) > 0 {
		return lang[0]
	}
	if lang := md.Get("grpcgateway-accept-language"); len(lang) > 0 {
		return lang[0]
	}
	return defaultLang
}

func UnaryLocaleInterceptor(l *locale.Locale) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return nil, err
			}
			md, _ := metadata.FromIncomingContext(ctx)
			lang := getLanguage(md)

			loc := l.CreateLocalizer(lang)
			msg := l.Translate(loc, st.Message(), nil)

			return nil, status.Errorf(st.Code(), msg)
		}
		return resp, nil
	}
}

func StreamLocaleInterceptor(l *locale.Locale) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)
		if err != nil {
			md, _ := metadata.FromIncomingContext(ss.Context())
			lang := getLanguage(md)

			loc := l.CreateLocalizer(lang)
			msg := l.Translate(loc, err.Error(), nil)

			return status.Errorf(codes.InvalidArgument, msg)
		}
		return nil
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
