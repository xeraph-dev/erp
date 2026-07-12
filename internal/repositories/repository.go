package repositories

type Repository interface {
	__internal()
}

type Repositories struct {
	User UserRepository
}
