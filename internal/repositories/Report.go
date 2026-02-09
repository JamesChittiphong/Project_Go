package repositories

import (
	"gorm.io/gorm"
)

type ReportRepository struct{ DB *gorm.DB }

func (r *ReportRepository) Create(report interface{}) error {
	return r.DB.Create(report).Error
}

func (r *ReportRepository) FindAll(reports interface{}) error {
	return r.DB.Find(reports).Error
}
