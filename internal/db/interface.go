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

//用于Search方法的参数
type SearchParams struct {
	Keyword string
	Value interface{}
	OrderBy string
	Order string
}

func NewSearchParams(keyword string, value interface{}, orderBy string, order string) *SearchParams {
	return &SearchParams{
		Keyword: keyword,
		Value: value,
		OrderBy: orderBy,
		Order: order,
	}
}


