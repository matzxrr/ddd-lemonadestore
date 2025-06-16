package interceptors

import (
    "context"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/status"
)

// LoggingInterceptor logs all gRPC requests
// WHY: Provides visibility into API usage and performance
// WHERE: Applied to all gRPC methods
func LoggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    // Log request start
    log.Printf("[gRPC] Started %s", info.FullMethod)
    
    // Call the handler
    resp, err := handler(ctx, req)
    
    // Log completion
    duration := time.Since(start)
    if err != nil {
        st, _ := status.FromError(err)
        log.Printf("[gRPC] Failed %s - Code: %s, Message: %s, Duration: %v",
            info.FullMethod, st.Code(), st.Message(), duration)
    } else {
        log.Printf("[gRPC] Completed %s - Duration: %v", info.FullMethod, duration)
    }
    
    return resp, err
}
