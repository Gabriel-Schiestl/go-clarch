package exceptions

type ServiceException struct {
	s string
}

func NewServiceException(s string) *ServiceException {
	return &ServiceException{s: s}
}

func (e *ServiceException) Error() string {
	return e.s
}