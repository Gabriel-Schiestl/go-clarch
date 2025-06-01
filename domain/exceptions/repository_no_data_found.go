package exceptions

type RepositoryNoDataFoundException struct {
	s string
}

func NewRepositoryNoDataFoundException(s string) *RepositoryNoDataFoundException {
	return &RepositoryNoDataFoundException{s: s}
}

func (e *RepositoryNoDataFoundException) Error() string {
	return e.s
}