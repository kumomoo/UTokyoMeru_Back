package db

import(
	"backend/internal/model"
	"gorm.io/gorm"

	"errors"
)

type GoodsCRUD struct{}

func (crud GoodsCRUD) CreateByObject(g model.Good) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Create(&g).Error
}

func (crud GoodsCRUD) FindAll() ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var goods []model.Good
	result := db.Preload("User").Find(&goods)
	return goods, result.Error
}

func (crud GoodsCRUD) FindById(id uint) (*model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var good model.Good
	result := db.First(&good, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &good, nil
}

func (crud GoodsCRUD) UpdateByObject(g model.Good) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Save(&g).Error
}

func (crud GoodsCRUD) DeleteById(id uint) error {
	obj, err:=crud.FindById(id)
	if err != nil {
		return err
	}
	obj.IsDeleted = true
	return crud.UpdateByObject(*obj)
}

func (crud GoodsCRUD) FindAllOrdered(fieldName string, order string) ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var goods []model.Good
	result := db.Order(fieldName + " " + order).Find(&goods)
	return goods, result.Error
}


func (crud GoodsCRUD) FindAllByField(fieldName string, value interface{}, orderBy string, order string) ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var goods []model.Good
	result := db.Where(fieldName+" = ?", value).Order(orderBy+" "+order).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	return goods, nil
}

func (crud GoodsCRUD) FindOneByUniqueField(fieldName string, value interface{}) (*model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var good model.Good
	result := db.Where(fieldName+" = ?", value).First(&good)
	if result.Error != nil {
		return nil, result.Error
	}

	return &good, nil
}

func (crud GoodsCRUD) Search(ops ...searchOption) ([]model.Good, error) {
	params := &SearchParams{}
	for _, op := range ops {
		op.f(params)
	}
	if params.Keyword == "" {
		return nil, errors.New("keyword is required")
	}
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var goods []model.Good
	var result *gorm.DB
	if params.OrderBy != "" && params.Order != "" {
		result = db.Where("title LIKE ? OR description LIKE ? OR tags LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%").Order(params.OrderBy+" "+params.Order).Find(&goods)
	}else {
		result = db.Where("title LIKE ? OR description LIKE ? OR tags LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%").Find(&goods)
	}
	return goods, result.Error
}
