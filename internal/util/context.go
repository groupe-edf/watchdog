package util

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string

// ContextValues context values
type ContextValues struct {
	ValuesMap map[ContextKey]interface{}
}

// Get return value for given key
func (values ContextValues) Get(key ContextKey) interface{} {
	return values.ValuesMap[key]
}

const (
	// ContextKeyRequestID is the ContextKey for RequestID
	ContextKeyRequestID ContextKey = "requestID"
	// ContextKeyStartTime is the ContextKey for UserID
	ContextKeyStartTime ContextKey = "startTime"
	// ContextKeyUserID is the ContextKey for UserID
	ContextKeyUserID ContextKey = "userId"
	// ContextKeyValues is the ContextKey for all values
	ContextKeyValues ContextKey = "contextValues"
)

// GetRequestID will get requestID from a git push
func GetRequestID(ctx context.Context) string {
	contextValues := ctx.Value(ContextKeyValues)
	if values, ok := contextValues.(ContextValues); ok {
		return values.Get(ContextKeyRequestID).(string)
	}
	return ""
}

// GetStartTime get start time
func GetStartTime(ctx context.Context) time.Time {
	contextValues := ctx.Value(ContextKeyValues)
	if values, ok := contextValues.(ContextValues); ok {
		return values.Get(ContextKeyStartTime).(time.Time)
	}
	return time.Now()
}

// GetUserID will get userID from a git push
func GetUserID(ctx context.Context) string {
	contextValues := ctx.Value(ContextKeyValues)
	if values, ok := contextValues.(ContextValues); ok {
		return values.Get(ContextKeyUserID).(string)
	}
	return ""
}

// InitializeContext initialize context
func InitializeContext() context.Context {
	ctx := context.Background()
	userID := os.Getenv("GL_USERNAME")
	if userID == "" {
		userID = "UNKNOWN"
	}
	values := ContextValues{
		ValuesMap: map[ContextKey]interface{}{
			ContextKeyRequestID: uuid.New().String(),
			ContextKeyStartTime: time.Now(),
			ContextKeyUserID:    userID,
		}}
	return context.WithValue(ctx, ContextKeyValues, values)
}
