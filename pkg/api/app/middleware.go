package app

type MiddlewareStore interface {
	Check() error
}

// Middleware Resource implements account management handler.
type MiddlewareResource struct {
	Store MiddlewareStore
}

func NewMiddlewareResource(store MiddlewareStore) *MiddlewareResource {
	return &MiddlewareResource{
		Store: store,
	}
}
