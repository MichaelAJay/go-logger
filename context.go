package logger

import "context"

type requestIDKey struct{}
type userIDKey struct{}
type sesisonIDKey struct{}

var (
	RequestIDKey = requestIDKey{}
	UserIDKey    = userIDKey{}
	SessionIDKey = sesisonIDKey{}
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(RequestIDKey).(string)
	return value, ok
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(UserIDKey).(string)
	return value, ok
}

func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

func GetSessionID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(SessionIDKey).(string)
	return value, ok
}
