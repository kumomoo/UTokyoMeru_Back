package db

import(
	"backend/internal/model"
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
	result := db.Preload("User").Preload("Comments").Find(&goods)
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

func (crud GoodsCRUD) Search(keyword string, orderBy string, order string) ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var goods []model.Good
	result := db.Where("title LIKE ? OR description LIKE ? OR tags LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Order(orderBy+" "+order).Find(&goods)
	return goods, result.Error
}
