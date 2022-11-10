package rest

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	FnGetter func(key string, defaultValue ...string) string

	QueryGetter interface {
		Query(key string, defaultValue ...string) string
	}

	Parser struct {
		source FnGetter
	}
)

const (
	PaginationParamPage  = "page"
	PaginationParamLimit = "limit"
	DefaultPage          = 1
)

type SimplePaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func QueryParam(s QueryGetter) *Parser {
	return &Parser{source: s.Query}
}

// StringFromEnum get value by key from source and compare with enum values.
// if value in source is empty returns empty string.
func (p *Parser) StringFromEnum(name string, enum []string) (string, error) {
	val := p.source(name)
	if val == "" {
		return val, nil
	}

	for _, s := range enum {
		if val == s {
			return s, nil
		}
	}

	return "", ErrBadRequestInvalidParameter(name, fmt.Sprintf(
		"Available values: %s",
		strings.Join(enum, ", "),
	))
}

// SimplePagination returns SimplePaginationParams from source or parsing error.
func (p *Parser) SimplePagination() (SimplePaginationParams, error) {
	var (
		err error
		pag = SimplePaginationParams{Page: DefaultPage}
	)

	if page := p.source(PaginationParamPage); page != "" {
		pag.Page, err = strconv.Atoi(page)
		if err != nil {
			return pag, ErrBadRequestInvalidParameter(PaginationParamPage).WithErr(err)
		}
	}

	if limit := p.source(PaginationParamLimit); limit != "" {
		pag.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return pag, ErrBadRequestInvalidParameter(PaginationParamLimit).WithErr(err)
		}
	}

	return pag, nil
}
