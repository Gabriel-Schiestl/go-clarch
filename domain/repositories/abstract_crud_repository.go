package repositories

type ID interface {
	~string | ~int
}

type AbstractCrudRepository[E any, I ID] interface {
	Save(entity E) (E, error)
	FindById(id I) (E, error)
	FindAll() ([]E, error)
	Delete(id I) error
}