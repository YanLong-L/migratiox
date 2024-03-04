package fixer

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OverrideFixer struct {
	// 因为本身其实这个不涉及什么领域对象，
	// 这里操作的不是 migrator 本身的领域对象
	base    *gorm.DB
	target  *gorm.DB
	columns []string
	Table   string
}

func NewOverrideFixer(base *gorm.DB,
	target *gorm.DB) (*OverrideFixer, error) {
	// 在这里需要查询一下数据库中究竟有哪些列
	rows, err := base.Table("").Limit(1).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &OverrideFixer{
		base:    base,
		target:  target,
		columns: columns,
		Table:   "",
	}, nil
}

func (o *OverrideFixer) Fix(ctx context.Context, id int64) error {
	var src map[string]interface{}
	// 找出数据
	err := o.base.WithContext(ctx).Table("").Where("id = ?", id).
		First(&src).Error
	switch err {
	// 找到了数据
	case nil:
		return o.target.Clauses(&clause.OnConflict{
			// 我们需要 Entity 告诉我们，修复哪些数据
			DoUpdates: clause.AssignmentColumns(o.columns),
		}).Create(&src).Error
	case gorm.ErrRecordNotFound:
		return o.target.Delete("id = ?", id).Error
	default:
		return err
	}
}
