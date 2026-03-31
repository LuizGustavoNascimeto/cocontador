package barrigade

import (
	"database/sql"
	"errors"
	"fmt"

	"cocontador/internal/repository"
)

type IBarrigadeHandler interface {
	Create(entity Barrigade) (Barrigade, error)
	GetByID(id int64) (Barrigade, error)
	List(page repository.PageRequest) ([]Barrigade, error)
	Update(entity Barrigade) (Barrigade, error)
	DeleteByID(id int64) error
	CreateOrGet(entity Barrigade) (Barrigade, error)
}

type BarrigadeHandler struct {
	repo BarrigadeRepository
}

var _ IBarrigadeHandler = (*BarrigadeHandler)(nil)

func NewHandler() (*BarrigadeHandler, error) {
	repo, err := NewRepository()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar barrigade handler: %w", err)
	}

	return &BarrigadeHandler{repo: repo}, nil
}

func (h *BarrigadeHandler) Create(entity Barrigade) (Barrigade, error) {
	return h.repo.Create(entity)
}

func (h *BarrigadeHandler) GetByID(id int64) (Barrigade, error) {
	return h.repo.GetByID(id)
}

func (h *BarrigadeHandler) List(page repository.PageRequest) ([]Barrigade, error) {
	return h.repo.List(page)
}

func (h *BarrigadeHandler) Update(entity Barrigade) (Barrigade, error) {
	return h.repo.Update(entity)
}

func (h *BarrigadeHandler) DeleteByID(id int64) error {
	return h.repo.DeleteByID(id)
}

func (h *BarrigadeHandler) CreateOrGet(entity Barrigade) (Barrigade, error) {
	existing, err := h.repo.GetByID(entity.ID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Barrigade{}, err
	}

	return h.repo.Create(entity)
}
