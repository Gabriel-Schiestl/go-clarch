package exceptions

type TechnicalException struct {
	s string
}

func NewTechnicalException(s string) *TechnicalException {
	return &TechnicalException{s: s}
}

func (e *TechnicalException) Error() string {
	return e.s
}