package exceptions

type BusinessException struct {
	s string
}

func NewBusinessException(s string) *BusinessException {
	return &BusinessException{s: s}
}

func (e *BusinessException) Error() string {
	return e.s
}