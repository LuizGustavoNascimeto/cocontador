package user

import "fmt"

type IUserHandler interface {
	Create(entity User) (User, error)
	ListAll() (map[string]User, error)
	CreateOrGet(entity User) (User, error)
}

type UserHandler struct {
	repo UserRepository
}

var _ IUserHandler = (*UserHandler)(nil)

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
func (h *UserHandler) CreateOrGet(entity User) (User, error) {
	existing, err := h.repo.GetByID(entity.ID)
	if err == nil {
		return existing, nil
	}
	if err.Error() != fmt.Sprintf("user %s nao encontrado: sql: no rows in result set", entity.ID) {
		return User{}, err
	}
	return h.repo.Create(entity)
}
