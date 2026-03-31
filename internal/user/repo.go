package user

import (
	"context"
	"database/sql"
	"fmt"

	"cocontador/internal/db"
	"cocontador/internal/repository"
)

// User representa a entidade de domínio usada no repositório.
type User struct {
	ID   string
	Name string
}

// UserRepository especializa o contrato genérico para User/int64.
type UserRepository interface {
	repository.Repository[User, string]
	ListAll() (map[string]User, error)
}

// PGRepository implementa UserRepository usando PostgreSQL.
type PGRepository struct {
	db  *sql.DB
	ctx context.Context
}

// Garante em tempo de compilação que PGRepository implementa UserRepository.
var _ UserRepository = (*PGRepository)(nil)

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

func (r *PGRepository) GetByID(id string) (User, error) {
	const q = `SELECT id, name FROM users WHERE id = $1`

	var u User
	err := r.db.QueryRowContext(r.ctx, q, id).Scan(&u.ID, &u.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user %s nao encontrado: %w", id, err)
		}
		return User{}, err
	}

	return u, nil
}

func (r *PGRepository) List(page repository.PageRequest) ([]User, error) {
	if page.Limit <= 0 {
		page.Limit = 20
	}

	if page.OrderBy == "" {
		page.OrderBy = "id"
	}

	query := fmt.Sprintf(
		"SELECT id, name FROM users ORDER BY %s %s LIMIT $1 OFFSET $2",
		page.OrderBy,
		map[bool]string{true: "DESC", false: "ASC"}[page.Desc],
	)

	rows, err := r.db.QueryContext(r.ctx, query, page.Limit, page.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0, page.Limit)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PGRepository) Create(entity User) (User, error) {
	const q = `INSERT INTO users (id,name) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(r.ctx, q, entity.ID, entity.Name).Scan(&entity.ID)
	if err != nil {
		return User{}, err
	}

	return entity, nil
}

func (r *PGRepository) Update(entity User) (User, error) {
	const q = `UPDATE users SET name = $1, WHERE id = $3`

	result, err := r.db.ExecContext(r.ctx, q, entity.Name, entity, entity.ID)
	if err != nil {
		return User{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	}
	if affected == 0 {
		return User{}, fmt.Errorf("user %s nao encontrado", entity.ID)
	}

	return entity, nil
}

func (r *PGRepository) DeleteByID(id string) error {
	const q = `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(r.ctx, q, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user %s nao encontrado", id)
	}

	return nil
}

func (r *PGRepository) ListAll() (map[string]User, error) {
	const q = `SELECT id, name FROM users`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]User)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users[u.ID] = u

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
