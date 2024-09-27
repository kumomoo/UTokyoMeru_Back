package db

type CRUD [T any] interface  {
	CreateByObject(T) error
	FindAll() ([]T , error)
	FindById(id uint) (*T, error)
	UpdateByObject(T) error
	DeleteById(id uint) error
}
