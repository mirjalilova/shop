package repo

import (
	"context"
	"fmt"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"strconv"
	"strings"
	"time"
)

type DebtLogsRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewDebtLogsRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *DebtLogsRepo {
	return &DebtLogsRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *DebtLogsRepo) Create(ctx context.Context, req *entity.DebtLogCreate) error {
	query := `
		INSERT INTO debt_logs (user_id, amount, reason)
		VALUES($1, $2, $3)
	`
	_, err := r.pg.Pool.Exec(ctx, query,
		req.UserID,
		req.Amount,
		req.Reason,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *DebtLogsRepo) GetDebtLogs(ctx context.Context, user_id string, status string) (*entity.DebtLogGetAllRes, error) {
	var res entity.DebtLogGetAllRes
	query := `SELECT 
				COUNT(id) OVER () AS total_count,
				d.id,
				d.user_id,
				u.name,
				d.amount,
				d.reason,
				d.debt_type,
				d.time_taken,
				d.given_time
			FROM debt_logs d
			JOIN users u ON d.user_id = u.id
			WHERE d.deleted_at = 0`

	if status != "" {
		query += fmt.Sprintf(" AND d.debt_type = '%s'", status)
	}

	if user_id != "" {
		query += fmt.Sprintf(" AND d.user_id = '%s'", user_id)
	}

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debtLogs []entity.DebtLogRes
	for rows.Next() {
		var debtLog entity.DebtLogRes
		var timeTaken, givenTime time.Time
		var count int

		if err := rows.Scan(
			&count,
			&debtLog.Id,
			&debtLog.UserID,
			&debtLog.UserName,
			&debtLog.Amount,
			&debtLog.Reason,
			&debtLog.Status,
			&timeTaken,
			&givenTime,
		); err != nil {
			return nil, err
		}

		debtLog.TakenTime = timeTaken.Format(time.RFC3339)
		debtLog.GivenTime = givenTime.Format(time.RFC3339)
		res.Count = count
		debtLogs = append(debtLogs, debtLog)
	}

	res.DebtLogs = debtLogs

	return &res, nil
}

func (r *DebtLogsRepo) Update(ctx context.Context, req entity.DebtLogUpdate) error {
	query := `
		UPDATE debt_logs
		SET 
	`
	var conditions []string
	var args []interface{}

	if req.Amount != 0 {
		conditions = append(conditions, " amount = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Amount)
	}

	if req.Reason != "" && req.Reason != "string" {
		conditions = append(conditions, " reason = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Reason)
	}

	if req.Status != "" && req.Status != "string" {
		conditions = append(conditions, " debt_type = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Status)
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0"

	args = append(args, req.ID)

	_, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
