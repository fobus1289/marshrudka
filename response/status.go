package response

import "net/http"

func Response() IResponse {
	return &response{
		Content: TEXT,
		Status:  http.StatusOK,
	}
}

func OK() ISend {
	return Response().Send(http.StatusOK)
}

func Created() ISend {
	return Response().Send(http.StatusCreated)
}

func NoContent() ISend {
	return Response().Abort(http.StatusNoContent)
}

func BadRequest() ISend {
	return Response().Abort(http.StatusBadRequest)
}

func Unauthorized() ISend {
	return Response().Abort(http.StatusUnauthorized)
}

func Forbidden() ISend {
	return Response().Abort(http.StatusForbidden)
}

func InternalServerError() ISend {
	return Response().Abort(http.StatusInternalServerError)
}

func BadGateway() ISend {
	return Response().Abort(http.StatusBadGateway)
}

func ServiceUnavailable() ISend {
	return Response().Abort(http.StatusServiceUnavailable)
}
