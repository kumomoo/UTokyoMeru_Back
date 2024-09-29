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

func (crud GoodsCRUD) FindAllOrdered() ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var goods []model.Good
	result := db.Order("updated_at DESC").Find(&goods)
	return goods, result.Error
}

func (crud GoodsCRUD) FindByUserId(userId uint) ([]model.Good, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var goods []model.Good
	result := db.Where("user_id = ?", userId).Find(&goods)
	return goods, result.Error
}