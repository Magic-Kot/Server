package postgresql

import (
	"3.Server/pkg/postg"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type repository struct {
	client postg.Client
	log    *slog.Logger
}

func NewRepository(client postg.Client, logger *slog.Logger) *repository {
	return &repository{
		client: client,
		log:    logger,
	}
}

// Create - создание новой записи в таблице url
func (r *repository) Create(ctx context.Context, nameurl string, url string) (int, error) {
	q := `
		INSERT INTO url 
		    (alias, url) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`

	var id int

	if err := r.client.QueryRow(ctx, q, nameurl, url).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create url: %w", err)
	}

	r.log.Info("a new url has been saved")

	return id, nil
}

// Get - получение url по alias
func (r *repository) Get(ctx context.Context, alias string) (string, error) {
	q := `
		SELECT url
		FROM url
		WHERE alias = $1
	`

	var resURL string

	err := r.client.QueryRow(ctx, q, alias).Scan(&resURL)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	r.log.Info("url has been found")

	return resURL, nil
}

// GetAll - получение всех записей из таблицы url
func (r *repository) GetAll(ctx context.Context) ([]string, error) {
	q := `
		SELECT alias
		FROM url
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	alias := make([]string, 0)

	for rows.Next() {
		var al string

		err = rows.Scan(&al)
		if err != nil {
			return nil, err
		}

		alias = append(alias, al)
	}

	return alias, nil
}

// Delete - удаление записи из таблицы url
func (r *repository) Delete(ctx context.Context, alias string) error {
	q := `
		DELETE FROM url
		WHERE alias = $1
	`

	commandTag, err := r.client.Exec(ctx, q, alias)
	if err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	r.log.Info("url has been deleted")

	return nil
}

// Update - обновление записи в таблице url
/*func (r *repository) Update(ctx context.Context, alias string, url string) error {
q := `
	UPDATE url
	SET url = $1
	WHERE alias = $2
`*/
