// Package interceptor contains gRPC interceptors.
package interceptor

import (
	"context"

	"github.com/devldavydov/gophkeeper/internal/common/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const MetaUserID = "USER_ID"

// AuthTokenInterceptor represents interceptor to validate requests against token.
type AuthTokenInterceptor struct {
	protectedMethods map[string]bool
	serverSecret     []byte
}

// NewAuthTokenInterceptor creates new AuthTokenInterceptor object.
func NewAuthTokenInterceptor(protectedMethods []string, serverSecret []byte) *AuthTokenInterceptor {
	pm := make(map[string]bool, len(protectedMethods))
	for _, p := range protectedMethods {
		pm[p] = true
	}
	return &AuthTokenInterceptor{protectedMethods: pm, serverSecret: serverSecret}
}

// Handle implements authentication for request.
//
// If incoming request does not containt metadata with token
// or token invalid, returns PersmissionDenied gRPC code.
//
// Else retrieves user id from token, adds token to request context and calls request handler.
func (a *AuthTokenInterceptor) Handle(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if !a.protectedMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	tokens := md.Get(token.HeaderName)
	if len(tokens) != 1 {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	userID, err := token.GetUserFromJWT(tokens[0], a.serverSecret)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	newMD := md.Copy()
	newMD.Append(MetaUserID, userID)

	return handler(metadata.NewIncomingContext(ctx, newMD), req)
}
