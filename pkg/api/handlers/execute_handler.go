package handlers

import (
	"mentalartsapi/pkg/runner/client"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type ExecuteRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
	Timeout  int32  `json:"timeout,omitempty"`
}

var runnerClient *client.RunnerClient

func InitRunnerClient() error {
	var err error
	runnerClient, err = client.NewRunnerClient("localhost:8001")
	if err != nil {
		return err
	}
	return nil
}

func HandleExecute(c *gin.Context) {
	ctx := c.Request.Context()

	// get span from context
	span := trace.SpanFromContext(ctx)
	defer span.End()

	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default timeout if not specified
	if req.Timeout == 0 {
		req.Timeout = 30 // 30 seconds default
	}

	// Execute code
	response, err := runnerClient.RunCode(ctx, req.Code, req.Language, req.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
