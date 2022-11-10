package psql

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type PgStatStatementRepo struct {
	db *gorm.DB
}

func NewPgStatStatementRepo(db *gorm.DB) *PgStatStatementRepo {
	return &PgStatStatementRepo{db}
}

func (r *PgStatStatementRepo) GetByType(minExecTime uint64, qType string, timeSpentSort string, limit, offset int) ([]PgStatStatement, error) {
	var statements []PgStatStatement

	q := r.db.Where("max_exec_time >= ?", minExecTime).
		Limit(limit).Offset(offset)

	if qType != "" {
		q = q.Where("query like ? OR query like ?", strings.ToLower(qType)+"%", strings.ToUpper(qType)+"%")
	}

	if timeSpentSort != "" {
		q = q.Order(fmt.Sprintf("max_exec_time %s", timeSpentSort))
	}

	return statements, q.Debug().Find(&statements).Error
}
