package errors

import "context"

type contextKeyPrintStackTrace struct{}

func WithStackTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyPrintStackTrace{}, true)
}

func shouldPrintStackTrace(ctx context.Context) bool {
	if printStackTrace, exists := ctx.Value(contextKeyPrintStackTrace{}).(bool); exists && printStackTrace {
		return true
	}

	return false
}
