package utils

import (
	"backend/internal/model"
)

type GoodTransform struct{}
type UserTransform struct{}
type CommentTransform struct{}

func (gt *GoodTransform) Db2ResponseModel(dbModel model.Good) model.GetGoodsResponse {
	return model.GetGoodsResponse{
		GoodID:      dbModel.ID,
		CreatedTime: dbModel.CreatedAt,
		UpdatedTime: dbModel.UpdatedAt,
		Title:       dbModel.Title,
		Description: dbModel.Description,
		Images:      dbModel.Images,
		Price:       dbModel.Price,
		Views:       dbModel.Views,
		IsInvisible: dbModel.IsInvisible,
		IsDeleted:   dbModel.IsDeleted,
		IsBought:    dbModel.IsBought,
		Tags:        dbModel.Tags,
		SellerID:    dbModel.SellerID,
	}
}

func (gt *GoodTransform) Post2DbModel(apiModel model.PostGoodsReceive) model.Good {
	return model.Good{
		Title:       apiModel.Title,
		Description: apiModel.Description,
		Images:      apiModel.Images,
		Price:       apiModel.Price,
		Tags:        apiModel.Tags,
		SellerID:    apiModel.SellerID,
	}
}
