package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	routing "github.com/qiangxue/fasthttp-routing"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const uuidStringLength = 36

type Server struct {
	log      *zap.Logger
	cfg      *Config
	router   *routing.Router
	apiGroup *routing.RouteGroup
}

type RequestContext struct {
	srv            *Server
	routingContext *routing.Context
}

type ResponseError struct {
	statusCode int
	message    string
}

func (e *ResponseError) Error() string {
	return e.message
}

type ResponseErrorM struct {
	Error string `json:"error"`
}

func NewResponseError(statusCode int, message string) error {
	return &ResponseError{statusCode: statusCode, message: message}
}

// Creates empty HTTP server, not runs it.
func NewServer(log *zap.Logger, cfg *Config) *Server {
	router := routing.New()
	apiGroup := router.Group("/api")
	return &Server{
		log:      log,
		cfg:      cfg,
		router:   router,
		apiGroup: apiGroup,
	}
}

func (srv *Server) Get(path string, handler func(c RequestContext) error) {
	srv.apiGroup.Get(path, func(routingContext *routing.Context) error {
		rc := RequestContext{routingContext: routingContext, srv: srv}
		var err error
		if err = handler(rc); err != nil {
			// if handler returned our 'status code' error, we can handle it nicely
			if responseError, ok := err.(*ResponseError); ok {
				routingContext.SetStatusCode(responseError.statusCode)
				return rc.AnswerJson(ResponseErrorM{Error: responseError.message})
			}
		}
		return err
	})
}

// Run HTTP server.
func (srv *Server) Run() {
	srv.log.Info("listening http", zap.String("address", srv.cfg.Api.Listen))
	httpServer := &fasthttp.Server{
		Handler: srv.router.HandleRequest,
		Name:    srv.cfg.AppName,
	}
	if err := httpServer.ListenAndServe(srv.cfg.Api.Listen); err != nil {
		log.Fatal("error on server listening", zap.Error(err))
	}
}

// UuidParam returns parsed id in UUID format from request uri.
// name - uri context variable name.
func (ctx *RequestContext) UuidParam(name string) (res uuid.UUID, err error) {
	str := ctx.routingContext.Param(name)
	if len(str) != uuidStringLength {
		message := fmt.Sprintf("uri parameter %s should be uuid, but value length was incorrect", name)
		ctx.srv.log.Debug(message)
		return uuid.Nil, NewResponseError(http.StatusBadRequest, message)
	}
	if res, err = uuid.FromString(str); err != nil {
		message := fmt.Sprintf("uri parameter %s should be uuid, but value was incorrect", name)
		ctx.srv.log.Debug(message)
		return uuid.Nil, NewResponseError(http.StatusBadRequest, message)
	}
	return res, nil
}

// IntQueryParam returns parsed int from request query params.
// name - query variable name.
func (ctx *RequestContext) IntQueryParam(name string) (res int, err error) {
	bytes := ctx.routingContext.QueryArgs().Peek(name)
	if len(bytes) == 0 {
		message := fmt.Sprintf("query parameter %s should be int, but value was empty", name)
		ctx.srv.log.Debug(message)
		return res, NewResponseError(http.StatusBadRequest, message)
	}
	if res, err = strconv.Atoi(string(bytes)); err != nil {
		message := fmt.Sprintf("query parameter %s should be int, but value was incorrect", name)
		ctx.srv.log.Debug(message)
		return res, NewResponseError(http.StatusBadRequest, message)
	}
	return res, nil
}

func (ctx *RequestContext) AnswerJson(value interface{}) error {
	ctx.routingContext.SetContentType("application/json")
	if err := json.NewEncoder(ctx.routingContext).Encode(value); err != nil {
		// it is debug level, because there is no sence to see these strange messages
		// like "client disconected before finish sending" in production environment
		ctx.srv.log.Debug("error while answering to client", zap.Error(err))
		return err
	}
	return nil
}

func (ctx *RequestContext) AnswerBadRequest(msg string) error {
	return NewResponseError(http.StatusBadRequest, msg)
}

func (ctx *RequestContext) AnswerNotFound(msg string) error {
	return NewResponseError(http.StatusNotFound, msg)
}

func (ctx *RequestContext) AnswerInternalError(msg string) error {
	return NewResponseError(http.StatusInternalServerError, msg)
}

func (ctx *RequestContext) WrapActionsError(err *ActionError) error {
	if err.IsNotFound() {
		return ctx.AnswerNotFound(err.Error())
	}
	return ctx.AnswerInternalError(err.Error())
}
