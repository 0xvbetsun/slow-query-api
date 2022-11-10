package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vbetsun/slow-query-api/internal/storage/psql"
	"github.com/vbetsun/slow-query-api/internal/transport/rest/dto"
)

const (
	FilterParamType = "type"
	SortParam       = "time_spent"
)

var (
	sortEnum           = []string{psql.SortASC, psql.SortDESC}
	statementTypesEnum = []string{
		psql.StatementINSERT,
		psql.StatementUPDATE,
		psql.StatementSELECT,
		psql.StatementDELETE,
	}
)

type (
	StatementsRepo interface {
		GetByType(maxExecTime uint64, qType string, timeSpent string, limit, offset int) ([]psql.PgStatStatement, error)
	}

	Handler struct {
		repo             StatementsRepo
		maxQueryDuration uint64
	}
)

func NewHandler(repo StatementsRepo, maxQueryDuration uint64) *Handler {
	return &Handler{repo, maxQueryDuration}
}

// GetAll returns all slow query in database
func (h *Handler) GetAll(ctx *fiber.Ctx) error {
	p := QueryParam(ctx)

	qType, err := p.StringFromEnum(FilterParamType, statementTypesEnum)
	if err != nil {
		return err
	}

	pagination, err := p.SimplePagination()
	if err != nil {
		return err
	}

	sort, err := p.StringFromEnum(SortParam, sortEnum)
	if err != nil {
		return err
	}

	statements, err := h.repo.GetByType(
		h.maxQueryDuration,
		qType, sort,
		pagination.Limit,
		pagination.Limit*(pagination.Page-1),
	)
	if err != nil {
		return err
	}

	queries := make([]dto.Query, len(statements))
	for i, statement := range statements {
		queries[i] = dto.Query{
			ID:           statement.QueryID,
			Statement:    statement.Query,
			MaxExecTime:  statement.MaxExecTime,
			MeanExecTime: statement.MeanExecTime,
		}
	}

	return ctx.JSON(ResponseWithPayload{
		Pagination: Pagination{
			Page:     pagination.Page,
			PageSize: len(queries),
		},
		Payload: queries,
	})
}
