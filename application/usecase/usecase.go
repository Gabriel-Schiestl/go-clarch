package usecase

import "context"

type UseCase[R any] interface {
	Execute(ctx context.Context) (R, error)
}

type UseCaseWithProps[P, R any] interface {
	Execute(ctx context.Context, props P) (R, error)
}