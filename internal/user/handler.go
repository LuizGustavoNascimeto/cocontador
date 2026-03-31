package user

import "fmt"

type UserHanlder interface {
	Create(entity User) (User, error)
}

type UserHandler struct {
	repo UserRepository
}

var _ UserHanlder = (*UserHandler)(nil)

func NewHandler() (*UserHandler, error) {
	repo, err := NewRepository()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar user handler: %w", err)
	}

	return &UserHandler{repo: repo}, nil
}

func (h *UserHandler) Create(entity User) (User, error) {
	return h.repo.Create(entity)
}

func (h *UserHandler) ListAll() (map[string]User, error) {
	return h.repo.ListAll()
}
