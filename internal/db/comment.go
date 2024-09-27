package db

import "backend/internal/model"

type CommentsCRUD struct{}

func (crud CommentsCRUD) CreateByObject(c model.Comment) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Create(&c).Error
}

func (crud CommentsCRUD) FindAll() ([]model.Comment, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var comments []model.Comment
	result := db.Find(&comments)
	return comments, result.Error
}

func (crud CommentsCRUD) FindById(id uint) (*model.Comment, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var comment model.Comment
	result := db.First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (crud CommentsCRUD) UpdateByObject(c model.Comment) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Save(&c).Error
}

func (crud CommentsCRUD) DeleteById(id uint) error {
	obj, err:=crud.FindById(id)
	if err != nil {
		return err
	}
	obj.IsDeleted = true
	return crud.UpdateByObject(*obj)
}