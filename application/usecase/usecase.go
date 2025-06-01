package usecase

type UseCase[R any] interface {
	Execute() (R, error)
}

type UseCaseWithProps[P, R any] interface {
	Execute(props P) (R, error)
}