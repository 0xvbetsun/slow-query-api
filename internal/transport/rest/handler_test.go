package rest

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/vbetsun/slow-query-api/internal/storage/psql"
	"gorm.io/gorm"
)

type StatementsRepoMock struct {
	GetByTypeMock func(minExecTime uint64, qType string, timeSpentSort string, limit, offset int) ([]psql.PgStatStatement, error)
}

func (s StatementsRepoMock) GetByType(minExecTime uint64, qType string, timeSpentSort string, limit, offset int) ([]psql.PgStatStatement, error) {
	return s.GetByTypeMock(minExecTime, qType, timeSpentSort, limit, offset)
}

func TestHandler_GetAll(t *testing.T) {
	tests := map[string]struct {
		repo        StatementsRepo
		queryParams map[string]string
		exp         ResponseWithPayload
		expErr      error
	}{
		"valid": {
			repo: StatementsRepoMock{
				GetByTypeMock: func(uint64, string, string, int, int) ([]psql.PgStatStatement, error) {
					return []psql.PgStatStatement{
						{QueryID: 12, Query: "some query", Calls: 2, MinExecTime: 10, MaxExecTime: 20, MeanExecTime: 15, TotalExecTime: 35, Rows: 2},
					}, nil
				},
			},
			queryParams: map[string]string{
				"limit":      "2",
				"page":       "1",
				"time_spent": "asc",
				"type":       "INSERT",
			},
			exp: ResponseWithPayload{
				Pagination: Pagination{Page: 1, PageSize: 1},
				Payload: []interface{}{map[string]interface{}{
					"id":             float64(12),
					"max_exec_time":  float64(20),
					"mean_exec_time": float64(15),
					"statement":      "some query",
				}},
			},
		},
		"wrong_param_time_spent_sort": {
			queryParams: map[string]string{
				"time_spent": "invalid-sort",
			},
			expErr: &ErrWithHint{
				Code:    400,
				Message: "Invalid parameter 'time_spent'. Available values: asc, desc",
				Field:   "time_spent_sort",
			},
		},
		"wrong_param_type": {
			queryParams: map[string]string{
				"type": "invalid_type",
			},
			expErr: &ErrWithHint{
				Code:    400,
				Message: "Invalid parameter 'type'. Available values: INSERT, UPDATE, SELECT, DELETE",
				Field:   "type",
			},
		},
		"wrong_param_page": {
			queryParams: map[string]string{
				"page": "asd",
			},
			expErr: &ErrWithHint{
				Code:    400,
				Message: "Invalid parameter 'page'. ",
				Field:   "page",
			},
		},
		"wrong_param_limit": {
			queryParams: map[string]string{
				"page":  "1",
				"limit": "asd",
			},
			expErr: &ErrWithHint{
				Code:    400,
				Message: "Invalid parameter 'limit'. ",
				Field:   "limit",
			},
		},
		"DB_error": {
			repo: StatementsRepoMock{
				GetByTypeMock: func(uint64, string, string, int, int) ([]psql.PgStatStatement, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expErr: errors.New("record not found"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewHandler(tt.repo, 2000)

			q := make(url.Values)
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}

			var ctx fasthttp.RequestCtx
			var req fasthttp.Request
			req.SetRequestURI("http://example.com/queries?" + q.Encode())
			ctx.Init(&req, nil, nil)
			err := c.GetAll(fiber.New().AcquireCtx(&ctx))
			if err != nil {
				if tt.expErr.Error() != err.Error() {
					t.Fatalf("expected: %v, got: %v", tt.expErr, err)
				}
				return
			}

			var result ResponseWithPayload
			err = json.Unmarshal(ctx.Response.Body(), &result)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.exp, result) {
				t.Fatalf("expected: %v, got: %v", tt.exp, result)
			}
		})
	}
}
