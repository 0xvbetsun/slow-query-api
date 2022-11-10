package dto

type Query struct {
	ID           int64   `json:"id"`
	Statement    string  `json:"statement"`
	MaxExecTime  float64 `json:"max_exec_time"`
	MeanExecTime float64 `json:"mean_exec_time"`
}
