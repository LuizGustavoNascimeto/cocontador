package barrigade

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cocontador/internal/db"
	"cocontador/internal/repository"
)

// Barrigade representa a entidade de domínio usada no repositório.
type Barrigade struct {
	ID         int64
	User_id    string
	Created_at time.Time
}

// BarrigadeRepository especializa o contrato genérico para Barrigade/int64.
type BarrigadeRepository interface {
	repository.Repository[Barrigade, int64]
}

// PGRepository implementa BarrigadeRepository usando PostgreSQL.
type PGRepository struct {
	db  *sql.DB
	ctx context.Context
}

// Garante em tempo de compilação que PGRepository implementa BarrigadeRepository.
var _ BarrigadeRepository = (*PGRepository)(nil)

func NewRepository() (*PGRepository, error) {
	conn, err := db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexao com banco: %w", err)
	}

	return &PGRepository{
		db:  conn,
		ctx: context.Background(),
	}, nil
}

func (r *PGRepository) GetByID(id int64) (Barrigade, error) {
	const q = `SELECT id, user_id, created_at FROM barrigades WHERE id = $1`

	var b Barrigade
	err := r.db.QueryRowContext(r.ctx, q, id).Scan(&b.ID, &b.User_id, &b.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return Barrigade{}, fmt.Errorf("barrigade %d nao encontrado: %w", id, err)
		}
		return Barrigade{}, err
	}

	return b, nil
}

func (r *PGRepository) List(page repository.PageRequest) ([]Barrigade, error) {
	if page.Limit <= 0 {
		page.Limit = 20
	}

	if page.OrderBy == "" {
		page.OrderBy = "id"
	}

	query := fmt.Sprintf(
		"SELECT id, user_id, created_at FROM barrigades ORDER BY %s %s LIMIT $1 OFFSET $2",
		page.OrderBy,
		map[bool]string{true: "DESC", false: "ASC"}[page.Desc],
	)

	rows, err := r.db.QueryContext(r.ctx, query, page.Limit, page.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	barrigades := make([]Barrigade, 0, page.Limit)
	for rows.Next() {
		var b Barrigade
		if err := rows.Scan(&b.ID, &b.User_id, &b.Created_at); err != nil {
			return nil, err
		}
		barrigades = append(barrigades, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return barrigades, nil
}

func (r *PGRepository) Create(entity Barrigade) (Barrigade, error) {
	const q = `INSERT INTO barrigades (user_id, created_at) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(r.ctx, q, entity.User_id, entity.Created_at).Scan(&entity.ID)
	if err != nil {
		return Barrigade{}, err
	}

	return entity, nil
}

func (r *PGRepository) Update(entity Barrigade) (Barrigade, error) {
	const q = `UPDATE barrigades SET user_id = $1, created_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(r.ctx, q, entity.User_id, entity.Created_at, entity.ID)
	if err != nil {
		return Barrigade{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Barrigade{}, err
	}
	if affected == 0 {
		return Barrigade{}, fmt.Errorf("barrigade %d nao encontrado", entity.ID)
	}

	return entity, nil
}

func (r *PGRepository) DeleteByID(id int64) error {
	const q = `DELETE FROM barrigades WHERE id = $1`

	result, err := r.db.ExecContext(r.ctx, q, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("barrigade %d nao encontrado", id)
	}

	return nil
}
