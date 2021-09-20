package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/internal"
	stdhttp "net/http"
)

var _ api.Service = (*svc)(nil)

type (
	svc    struct{}
	Server struct {
		Port    uint64
		server  *stdhttp.Server
		Debug   bool
		service api.Service
	}
)

func NewServer(port uint64, debug bool) *Server {
	sv := &Server{
		Port:    port,
		Debug:   debug,
		service: &svc{},
	}

	h := sv.handler()

	sv.server = &stdhttp.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h,
	}

	return sv
}
func sendErrMessage(writer stdhttp.ResponseWriter, code int, err api.ErrResponse) {
	response := internal.NewResponse(code, err)
	internal.Reply(response, writer)
}
func (server *Server) addHandler(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
	var req api.Request
	err := internal.ReceivePayload(request, &req)
	if err != nil {
		errMessage := api.ErrResponse{
			Err:     err.Error(),
			Message: "failed to obtain request body",
		}
		sendErrMessage(writer, stdhttp.StatusInternalServerError, errMessage)
		return
	}

	rs, err := server.service.Add(req.A, req.B)
	if err != nil {
		errMessage := api.ErrResponse{
			Err:     err.Error(),
			Message: "failed to perform addition",
		}

		sendErrMessage(writer, stdhttp.StatusInternalServerError, errMessage)

		return
	}
	result := api.Response{
		Answer: rs,
	}
	response := internal.NewResponse(stdhttp.StatusOK, result)
	internal.Reply(response, writer)
}

func (server *Server) divisionHandler(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
	defer func() {
		if r := recover(); r != nil {
			errMessage := api.ErrResponse{
				Err:     "division by zero",
				Message: "division by zero is not good dont do it",
			}
			sendErrMessage(writer, stdhttp.StatusBadRequest, errMessage)
			return
		}
	}()
	var req api.Request
	err := internal.ReceivePayload(request, &req)
	if err != nil {
		errMessage := api.ErrResponse{
			Err:     err.Error(),
			Message: "failed to obtain request body",
		}
		sendErrMessage(writer, stdhttp.StatusInternalServerError, errMessage)
		return
	}

	rs, err := server.service.Divide(req.A, req.B)
	if err != nil {
		errMessage := api.DivResponse{
			Err:     err.Error(),
			Message: "failed to perform division",
		}

		r := internal.NewResponse(stdhttp.StatusInternalServerError, errMessage)
		internal.Reply(r, writer)
		return
	}
	result := api.DivResponse{
		Answer: rs,
	}
	response := internal.NewResponse(stdhttp.StatusOK, result)
	internal.Reply(response, writer)
}

func (server *Server) handler() stdhttp.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/add", server.addHandler).Methods(stdhttp.MethodGet)
	r.HandleFunc("/div", server.divisionHandler).Methods(stdhttp.MethodGet)
	return r
}

func (server *Server) ListenAndServe() error {

	//s := stdhttp.Server{
	//	Addr:    fmt.Sprintf(":%d", server.Port),
	//	Handler: server.handler(),
	//}

	return server.server.ListenAndServe()

}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

func (s *svc) Add(a, b int64) (int64, error) {
	return a + b, nil
}

func (s *svc) Divide(a, b int64) (int64, error) {
	answer := a / b
	return answer, nil
}
