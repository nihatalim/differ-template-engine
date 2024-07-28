package repository

import (
	"context"
	"differ-template-engine/log"
	"differ-template-engine/pkg/client/nodiffer"
	"differ-template-engine/pkg/config"
	"github.com/jackc/pgx/v4"
)

type ExecutionResultRepository interface {
	SaveExecutionResult(operationId string, userId string, req nodiffer.HasDiffRequest, hasDiff bool) error
}

type executionResultRepository struct {
	cfg    *config.ApplicationConfig
	conn   *pgx.Conn
	logger log.Logger
}

var InsertExecutionResultSQL = "INSERT INTO execution_results (operation_id, user_id, request, diff) VALUES ($1, $2, $3, $4) ON CONFLICT (operation_id) DO UPDATE SET request = EXCLUDED.request, diff = EXCLUDED.diff, created_date = now() RETURNING operation_id"

func NewExecutionResultRepository(cfg *config.ApplicationConfig, conn *pgx.Conn, logger log.Logger) ExecutionResultRepository {
	return &executionResultRepository{
		cfg:    cfg,
		conn:   conn,
		logger: logger,
	}
}

func (e *executionResultRepository) SaveExecutionResult(operationId string, userId string, req nodiffer.HasDiffRequest, hasDiff bool) error {
	var opId string
	err := e.conn.QueryRow(context.Background(), InsertExecutionResultSQL, operationId, userId, req, hasDiff).Scan(&opId)
	if err != nil {
		return err
	}

	return nil
}
