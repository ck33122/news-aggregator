package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	router   *routing.Router
	apiGroup *routing.RouteGroup
}

type RequestContext struct {
	routingContext *routing.Context
}

type IActionError interface {
	error
	IsNotFound() bool
}

type ResponseError struct {
	statusCode int
	message    string
}

func (e *ResponseError) Error() string {
	return e.message
}

func NewResponseError(statusCode int, message string) error {
	return &ResponseError{statusCode: statusCode, message: message}
}

// Creates empty HTTP server, not runs it.
func NewServer() *Server {
	router := routing.New()
	apiGroup := router.Group("/api")
	return &Server{
		router:   router,
		apiGroup: apiGroup,
	}
}

func (server *Server) GetHandler(path string, handler func(c RequestContext) error) {
	server.apiGroup.Get(path, func(routingContext *routing.Context) error {
		rc := RequestContext{routingContext: routingContext}
		var err error
		if err = handler(rc); err != nil {
			// if handler returned our 'status code' error, we can handle it nicely
			if responseError, ok := err.(*ResponseError); ok {
				routingContext.SetStatusCode(responseError.statusCode)
				return rc.AnswerJson(struct{ Error string }{Error: responseError.message})
			}
		}
		return err
	})
}

// Run HTTP server on listenAddress.
func (server *Server) Run(listenAddress string) {
	log := GetLog()
	log.Info("listening http", zap.String("address", listenAddress))
	httpServer := &fasthttp.Server{
		Handler: server.router.HandleRequest,
		Name:    appName,
	}
	if err := httpServer.ListenAndServe(listenAddress); err != nil {
		log.Fatal("error on server listening", zap.Error(err))
	}
}

// UuidParam returns parsed id in UUID format from request uri.
// name - uri context variable name.
func (ctx RequestContext) UuidParam(name string) (res uuid.UUID, err error) {
	str := ctx.routingContext.Param(name)
	const uuidStringLength = 36
	if len(str) != uuidStringLength {
		return uuid.Nil, NewResponseError(
			http.StatusBadRequest,
			fmt.Sprintf("argument %s should be uuid, but value is incorrect", name),
		)
	}
	if res, err = uuid.FromString(str); err != nil {
		return uuid.Nil, NewResponseError(
			http.StatusBadRequest,
			fmt.Sprintf("argument %s should be uuid, but value is incorrect", name),
		)
	}
	return res, nil
}

func (ctx RequestContext) AnswerJson(value interface{}) error {
	if err := json.NewEncoder(ctx.routingContext).Encode(value); err != nil {
		// it is debug level, because there is no sence to see these strange messages
		// like "client disconected before finish sending" in production environment
		GetLog().Debug("error while answering to client", zap.Error(err))
		return err
	}
	return nil
}

func (ctx RequestContext) AnswerBadRequest(msg string) error {
	return NewResponseError(http.StatusBadRequest, msg)
}

func (ctx RequestContext) AnswerNotFound(msg string) error {
	return NewResponseError(http.StatusNotFound, msg)
}

func (ctx RequestContext) AnswerInternalError(msg string) error {
	return NewResponseError(http.StatusInternalServerError, msg)
}

func (ctx RequestContext) WrapActionsError(err IActionError) error {
	if err.IsNotFound() {
		return ctx.AnswerNotFound(err.Error())
	}
	return ctx.AnswerInternalError(err.Error())
}
