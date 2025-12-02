package interceptor

import (
    "context"
    "time"

    "go.uber.org/zap"
    "google.golang.org/grpc"
)

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {

        start := time.Now()
        resp, err := handler(ctx, req)
        duration := time.Since(start)

        cid, _ := ctx.Value(CidKey{}).(string)
        if cid == "" {
            cid = "unknown"
        }

        fields := []zap.Field{
            zap.String("cid", cid),
            zap.String("method", info.FullMethod),
            zap.Duration("latency", duration),
        }

        if err != nil {
            logger.Error("gRPC error", append(fields, zap.Error(err))...)
        } else {
            logger.Info("gRPC request completed", fields...)
        }

        return resp, err
    }
}

