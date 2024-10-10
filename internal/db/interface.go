package db

type CRUD [T any] interface  {
	CreateByObject(T) error
	FindAll() ([]T , error)
	FindAllOrdered(fieldName string, order string) ([]T, error)
	FindById(id uint) (*T, error)
	UpdateByObject(T) error
	DeleteById(id uint) error
	FindAllByField(fieldName string, value interface{}, orderBy string, order string) ([]T, error)
	FindOneByUniqueField(fieldName string, value interface{}) (*T, error)
	Search(keyword string, orderBy string, order string) ([]T, error)
}

const (
	OrderAsc = "ASC"
	OrderDesc = "DESC"
)
