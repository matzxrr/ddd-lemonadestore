package interceptors

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// ErrorInterceptor handles panics and ensures proper error responses
// WHY: Prevents server crashes and provides consistent error handling
func ErrorInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    // Recover from panics
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[gRPC] Panic in %s: %v", info.FullMethod, r)
            err = status.Error(codes.Internal, "internal server error")
        }
    }()
    
    // Call handler
    resp, err = handler(ctx, req)
    
    // Ensure errors are proper gRPC status
    if err != nil && status.Code(err) == codes.Unknown {
        err = status.Error(codes.Internal, err.Error())
    }
    
    return resp, err
}
