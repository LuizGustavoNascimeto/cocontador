package repository

// PageRequest define os parâmetros de paginação e ordenação.
type PageRequest struct {
	Limit   int
	Offset  int
	OrderBy string
	Desc    bool
}

// ReadRepository define operações de leitura fortemente tipadas.
type ReadRepository[T any, ID comparable] interface {
	GetByID(id ID) (T, error)
	List(page PageRequest) ([]T, error)
}

// WriteRepository define operações de escrita fortemente tipadas.
type WriteRepository[T any, ID comparable] interface {
	Create(entity T) (T, error)
	Update(entity T) (T, error)
	DeleteByID(id ID) error
}

// Repository compõe leitura e escrita em um único contrato.
type Repository[T any, ID comparable] interface {
	ReadRepository[T, ID]
	WriteRepository[T, ID]
}
