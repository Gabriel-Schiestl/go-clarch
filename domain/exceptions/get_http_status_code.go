package exceptions

func GetHTTPStatusCode(err error) int {
	if err == nil {
		return 0
	}

	switch err.(type) {
	case *BusinessException:
		return 400
	case *RepositoryNoDataFoundException:
		return 404
	case *ServiceException:
		return 400
	case *TechnicalException:
		return 500
	default:
		return 500
	}
}