package server

import (
	"context"
	"mentalartsapi/pkg/runner/pb"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RunnerServer struct {
	pb.UnimplementedRunnerServiceServer
}

func NewRunnerServer() *RunnerServer {
	return &RunnerServer{}
}

func (s *RunnerServer) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	startTime := time.Now()

	// Create temporary directory for code execution
	span.AddEvent("creating_temp_directory")
	tmpDir, err := os.MkdirTemp("", "runner-*")
	if err != nil {
		span.AddEvent("temp_directory_creation_failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Create source file based on language
	span.AddEvent("creating_source_file")
	filename := getSourceFilename(req.Language)
	filepath := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(filepath, []byte(req.Code), 0644); err != nil {
		span.AddEvent("source_file_creation_failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	// Prepare command based on language
	span.AddEvent("preparing_command")
	cmd := prepareCommand(req.Language, filepath)

	span.SetAttributes(
		attribute.String("command", cmd.String()),
		attribute.String("language", req.Language),
		attribute.String("code", req.Code),
	)
	// Set timeout
	if req.Timeout > 0 {
		span.AddEvent("setting_timeout", trace.WithAttributes(
			attribute.Int("timeout_seconds", int(req.Timeout)),
		))
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	}

	// Execute command
	span.AddEvent("executing_command")
	output, err := cmd.CombinedOutput()

	span.AddEvent("creating_response")
	response := &pb.RunResponse{
		ExecutionTime: time.Since(startTime).Milliseconds(),
		Output:        string(output),
		ExitCode:      0,
	}

	if err != nil {
		span.AddEvent("command_execution_failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("output", string(output)),
		))
		if exitErr, ok := err.(*exec.ExitError); ok {
			span.AddEvent("setting_exit_code", trace.WithAttributes(
				attribute.Int("exit_code", exitErr.ExitCode()),
			))
			response.ExitCode = int32(exitErr.ExitCode())
		}
		response.Error = err.Error()
	} else {
		span.AddEvent("command_execution_completed", trace.WithAttributes(
			attribute.String("output", string(output)),
		))
	}

	span.AddEvent("setting_final_attributes")
	span.SetAttributes(
		attribute.String("output", response.Output),
		attribute.Int("exit_code", int(response.ExitCode)),
		attribute.Int("execution_time", int(response.ExecutionTime)),
	)

	span.AddEvent("returning_response")
	return response, nil
}

func getSourceFilename(language string) string {
	switch language {
	case "go":
		return "main.go"
	case "python":
		return "main.py"
	case "bash":
		return "script.sh"
	default:
		return "main.txt"
	}
}

func prepareCommand(language, filepath string) *exec.Cmd {
	switch language {
	case "go":
		return exec.Command("go", "run", filepath)
	case "python":
		return exec.Command("python3", filepath)
	case "bash":
		return exec.Command("bash", filepath)
	default:
		return exec.Command("cat", filepath) // Just output the content for unknown languages
	}
}
