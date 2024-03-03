package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type InsuranceRepo struct {
	*postgres.Postgres
}

func NewInsuranceRepo(pg *postgres.Postgres) *InsuranceRepo {
	return &InsuranceRepo{pg}
}

// AddInsurance 添加一条用户的保险记录
func (i *InsuranceRepo) AddInsurance(ctx context.Context, insurance entity.Insurance) error {
	return i.Db.WithContext(ctx).Save(&insurance).Error
}

// GetInsuranceByUserID 获取用户所有的保险记录
func (i *InsuranceRepo) GetInsuranceByUserID(ctx context.Context, userID int) ([]entity.Insurance, error) {
	var insurances []entity.Insurance
	if err := i.Db.WithContext(ctx).Where("user_id = ?", userID).Limit(20).Order("id desc").Find(&insurances).Error; err != nil {
		return nil, err
	}

	return insurances, nil
}
