package log

import (
	"context"
	"time"

	logpb "github.com/yasinsaee/go-log-service/log-service/log"
	"github.com/yasinsaee/go-log-service/pkg/elastic"
)

type Handler struct {
	logpb.UnimplementedLogServiceServer
	Elastic *elastic.Client
}

func New(es *elastic.Client) *Handler {
	return &Handler{
		Elastic: es,
	}
}

func (h *Handler) WriteLog(ctx context.Context, req *logpb.LogRequest) (*logpb.LogResponse, error) {
	logEntry := elastic.LogEntry{
		Level:     req.Level,
		Message:   req.Message,
		Service:   req.Service,
		Module:    req.Module,
		RequestID: req.RequestId,
		UserID:    req.UserId,
		Host:      req.Host,
		Error:     req.Error,
		Timestamp: time.Now(),
	}

	extra := make(map[string]interface{})
	for k, v := range req.Extra {
		extra[k] = v
	}
	logEntry.Extra = extra

	if err := h.Elastic.IndexDocument("logs", logEntry); err != nil {
		return &logpb.LogResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &logpb.LogResponse{
		Success: true,
		Error:   "",
	}, nil
}
