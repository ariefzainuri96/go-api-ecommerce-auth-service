package interceptor

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CidKey struct{} // unique type → no collisions

const correlationIDHeader = "x-correlation-id"

func CorrelationIDInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {

        md, ok := metadata.FromIncomingContext(ctx)
        var cid string

        if ok && len(md[correlationIDHeader]) > 0 {
            cid = md[correlationIDHeader][0]
        } else {
            cid = uuid.NewString()
        }

        // ✔ Safe: use cidKey{} instead of string
        ctx = context.WithValue(ctx, CidKey{}, cid)

        logger.Info("Incoming gRPC request",
            zap.String("cid", cid),
            zap.String("method", info.FullMethod),
        )

        return handler(ctx, req)
    }
}

