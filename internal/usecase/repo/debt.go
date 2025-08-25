package repo

import (
	"context"
	"fmt"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"shop/pkg/telegram"
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
				COUNT(d.id) OVER () AS total_count,
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
		var timeTaken, givenTime *time.Time
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
		if givenTime != nil {
			formatted := givenTime.Format(time.RFC3339)
			debtLog.GivenTime = &formatted
		}

		res.Count = count
		debtLogs = append(debtLogs, debtLog)
	}

	res.DebtLogs = debtLogs

	return &res, nil
}

func (r *DebtLogsRepo) Update(ctx context.Context, req entity.DebtLogUpdate) error {

	query := "SELECT amount FROM debt_logs WHERE id = $1 AND deleted_at = 0"
	var currentAmount int
	err := r.pg.Pool.QueryRow(ctx, query, req.ID).Scan(&currentAmount)
	if err != nil {
		return err
	}

	query = `
		UPDATE debt_logs
		SET 
	`
	var conditions []string
	var args []interface{}

	if req.Amount != 0 {
		conditions = append(conditions, " amount = amount - $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Amount)
	}

	if req.Reason != "" && req.Reason != "string" {
		conditions = append(conditions, " reason = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Reason)
	}

	if req.Status != "" && req.Status != "string" {
		if req.Amount == currentAmount {
			conditions = append(conditions, " debt_type = $"+strconv.Itoa(len(args)+1))
			args = append(args, req.Status)
		}

		if req.Status == "gave" {
			conditions = append(conditions, " given_time = now()")
		}
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0"

	args = append(args, req.ID)

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *DebtLogsRepo) Report(ctx context.Context, report *entity.Report) error {
	query := `
		SELECT
			name,
			phone_number,
			debt
		FROM users
		WHERE id = $1 AND deleted_at = 0
	`
	var name, phoneNumber string
	var debt int
	err := r.pg.Pool.QueryRow(ctx, query, report.UserId).Scan(&name, &phoneNumber, &debt)
	if err != nil {
		return err
	}

	go r.sendDebtToTelegram(context.Background(), name, phoneNumber, report.Chek, debt)

	return nil
}
func (r *DebtLogsRepo) sendDebtToTelegram(ctx context.Context, name, phoneNumber, chek string, debt int) {

	message := fmt.Sprintf(
		"🆕 Qarz tolanganlik haqida xabar\n\n"+
			"👤 Mijoz: %s\n"+
			"📞 Telefon: %s\n"+
			"💳 Umumiy qarzi: %d\n"+
			"🕒 Sana: %s\n"+
			"📄 Chek: %s\n",
		name,
		phoneNumber,
		debt,
		time.Now().Format("2006-01-02 15:04:05"),
		chek,
	)

	// encodedURL := url.PathEscape(chek)

	tg := telegram.NewClient(r.config.Telegram.Token, r.config.Telegram.ChatID)
	if err := tg.SendMessage(message); err != nil {
		r.logger.Error("telegram send error: ", err)
	}
}
