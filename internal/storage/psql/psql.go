package psql

// describe statement types.
const (
	StatementSELECT = "SELECT"
	StatementINSERT = "INSERT"
	StatementUPDATE = "UPDATE"
	StatementDELETE = "DELETE"
)

// describes sorting directions
const (
	SortASC  = "asc"
	SortDESC = "desc"
)

// PgStatStatement model for pg_stat_statements table.
type PgStatStatement struct {
	QueryID       int64   `gorm:"column:queryid"`
	Query         string  `gorm:"column:query"`
	Calls         uint64  `gorm:"column:calls"`
	MinExecTime   float64 `gorm:"column:min_exec_time"`
	MaxExecTime   float64 `gorm:"column:max_exec_time"`
	MeanExecTime  float64 `gorm:"column:mean_exec_time"`
	TotalExecTime float64 `gorm:"column:total_exec_time"`
	Rows          uint64  `gorm:"column:rows"`
}
