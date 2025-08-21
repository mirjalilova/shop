package repo

import (
	"context"

	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"shop/pkg/ws"

	"github.com/gorilla/websocket"
)

type KassaRepo struct {
	pg      *postgres.Postgres
	config  *config.Config
	logger  *logger.Logger
	clients map[*websocket.Conn]bool
	// mu      sync.Mutex
	hub *ws.Hub
}

// New -.
func NewKassaRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger, hub *ws.Hub) *KassaRepo {
	return &KassaRepo{
		pg:      pg,
		config:  config,
		logger:  logger,
		clients: make(map[*websocket.Conn]bool),
		hub:     hub,
	}
}

func (r *KassaRepo) AddItem(ctx context.Context, productID string) error {
	query := `
		SELECT id, name, price, img_url, count, size
		FROM products 
		WHERE id = $1 AND deleted_at = 0
	`

	res := entity.KassaEvent{}
	err := r.pg.Pool.QueryRow(ctx, query, productID).
		Scan(&res.ProductID, &res.Name, &res.Price, &res.ImgURL, &res.Count, &res.Size)
	if err != nil {
		return err
	}

	// eventni broadcast qilish
	r.broadcastEvent(res)

	return nil
}

func (r *KassaRepo) Formalize(ctx context.Context, formalize []entity.Formalize) error {
	query := `
		UPDATE products set
			count = count - $1,
			sales_count = sales_count + $1
		WHERE id = $2
	`

	tr, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)

	for _, item := range formalize {
		_, err := tr.Exec(ctx, query, item.Count, item.ProductID)
		if err != nil {
			tr.Rollback(ctx)
			return err
		}
	}
	if err := tr.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *KassaRepo) broadcastEvent(event entity.KassaEvent) {
	if r.hub != nil {
		r.hub.Broadcast(event)
	}
}
