package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"ecochitas/internal/config"
)

type Api_server struct {
	http_server        *http.Server
	application_logger *slog.Logger
}

func New_api_server(
	application_config config.App_config,
	api_handler *Api_handler,
	application_logger *slog.Logger,
) *Api_server {
	http_multiplexer := http.NewServeMux()
	api_handler.Register_routes(http_multiplexer)

	server_address := fmt.Sprintf("%s:%s", application_config.Http_host, application_config.Http_port)
	http_server := &http.Server{
		Addr:         server_address,
		Handler:      with_cors_headers(http_multiplexer),
		ReadTimeout:  time.Duration(application_config.Request_read_timeout_seconds) * time.Second,
		WriteTimeout: time.Duration(application_config.Request_write_timeout_seconds) * time.Second,
		IdleTimeout:  time.Duration(application_config.Request_idle_timeout_seconds) * time.Second,
	}

	return &Api_server{
		http_server:        http_server,
		application_logger: application_logger,
	}
}

func (api_server *Api_server) Start() error {
	api_server.application_logger.Info("api_server_starting", "address", api_server.http_server.Addr)

	listen_server_error := api_server.http_server.ListenAndServe()
	if listen_server_error != nil && !errors.Is(listen_server_error, http.ErrServerClosed) {
		return fmt.Errorf("failed_to_listen_http_server: %w", listen_server_error)
	}

	return nil
}

func (api_server *Api_server) Shutdown(shutdown_context context.Context) error {
	api_server.application_logger.Info("api_server_shutting_down")
	return api_server.http_server.Shutdown(shutdown_context)
}

func with_cors_headers(next_handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response_writer http.ResponseWriter, request *http.Request) {
		response_writer.Header().Set("Access-Control-Allow-Origin", "*")
		response_writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		response_writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
		response_writer.Header().Set("Access-Control-Max-Age", "600")

		if request.Method == http.MethodOptions {
			response_writer.WriteHeader(http.StatusNoContent)
			return
		}

		next_handler.ServeHTTP(response_writer, request)
	})
}
