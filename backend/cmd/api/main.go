package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ecochitas/internal/app"
)

func main() {
	application_context, stop_signal_handler := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop_signal_handler()

	backend_application, create_application_error := app.New_backend_app()
	if create_application_error != nil {
		slog.Error("failed_to_create_backend_application", "error", create_application_error)
		os.Exit(1)
	}

	run_application_error := backend_application.Run(application_context)
	if run_application_error != nil {
		slog.Error("backend_application_stopped_with_error", "error", run_application_error)
		os.Exit(1)
	}
}
