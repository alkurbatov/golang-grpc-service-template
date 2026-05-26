package middleware

import (
	"context"

	"google.golang.org/grpc"
)

type ServerStreamMock struct {
	grpc.ServerStream

	ctx context.Context
}

func NewServerStreamMock(ctx context.Context) ServerStreamMock {
	return ServerStreamMock{ctx: ctx}
}

func (f ServerStreamMock) Context() context.Context {
	return f.ctx
}
