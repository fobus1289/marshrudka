package server

import (
	"net/http"

	"github.com/fobus1289/marshrudka/request"
	"github.com/fobus1289/marshrudka/response"
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/validator"
)

type (
	Server = router.IServer
	Router = router.IRouter
	Group  = router.IGroup

	Send = response.ISend

	Request    = request.IRequest
	FormFile   = request.IFormFile
	Param      = request.IParam
	QueryParam = request.IQueryParam
	Header     = request.IHeader

	Validator            = validator.IValidator
	ValidateErrorMessage = validator.MessageMapResult
)

var (
	ValidatorBuild = validator.Build
)

func New() Server {
	return router.NewServer()
}

func OK() Send {
	return response.Response().Send(http.StatusOK)
}

func Created() Send {
	return response.Response().Send(http.StatusCreated)
}

func NoContent() Send {
	return response.Response().Abort(http.StatusNoContent)
}

func BadRequest() Send {
	return response.Response().Abort(http.StatusBadRequest)
}

func Unauthorized() Send {
	return response.Response().Abort(http.StatusUnauthorized)
}

func Forbidden() Send {
	return response.Response().Abort(http.StatusForbidden)
}

func InternalServerError() Send {
	return response.Response().Abort(http.StatusInternalServerError)
}

func BadGateway() Send {
	return response.Response().Abort(http.StatusBadGateway)
}

func ServiceUnavailable() Send {
	return response.Response().Abort(http.StatusServiceUnavailable)
}
