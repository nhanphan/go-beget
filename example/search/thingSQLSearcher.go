package search

/*
	GENERATED CODE, YOUR EDITS ARE FUTILE

	generated from Thing
	using http://github.com/brianstarke/go-beget/generator/searcher
*/

import (
	"fmt"
	"strings"

	"github.com/brianstarke/go-beget/searcher"
	"github.com/jmoiron/sqlx"

	types "github.com/brianstarke/go-beget/example/types"
)

// ThingSearcher is the interface
type ThingSearcher interface {
	Search(searchRequest ThingSearchRequest) ([]types.Thing, error)
	Count(searchRequest ThingSearchRequest) (int32, error)
	GetByField(field ThingField, value interface{}) (*types.Thing, error)
}

// SQLThingSearcher implements a SQL based searcher
type SQLThingSearcher struct {
	db *sqlx.DB
}

// NewSQLThingSearcher returns a configured repo for you
func NewSQLThingSearcher(db *sqlx.DB) ThingSearcher {
	return &SQLThingSearcher{db: db}
}

// Search converts a ThingSearchRequest in to SQL and executes it
func (r *SQLThingSearcher) Search(searchRequest ThingSearchRequest) ([]types.Thing, error) {
	results := []types.Thing{}

	sqlStr, values, err := searcher.GenerateSelectSQL(&searchRequest, false)

	if err != nil {
		return nil, fmt.Errorf("Error generating SQL for ThingSearchRequest : %s", err.Error())
	}

	err = r.db.Select(&results, sqlStr, values.([]interface{})...)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return results, nil
		}
	}

	return results, err
}

// Count converts a ThingSearchRequest in to SQL and executes it
func (r *SQLThingSearcher) Count(searchRequest ThingSearchRequest) (int32, error) {
	var results []struct {
		Count int32 `db:"cnt"`
	}

	sqlStr, values, err := searcher.GenerateSelectSQL(&searchRequest, true)

	if err != nil {
		return 0, fmt.Errorf("Error generating SQL for ThingSearchRequest : %s", err.Error())
	}

	err = r.db.Select(&results, sqlStr, values.([]interface{})...)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return 0, fmt.Errorf("No results returned")
		}
	}

	return results[0].Count, err
}

// GetByField is a convenience method to return just one result by matching on a single field
func (r *SQLThingSearcher) GetByField(field ThingField, value interface{}) (*types.Thing, error) {
	var searchRequest ThingSearchRequest

	searchRequest.AddFilter(
		field,
		value,
		searcher.EQ,
		searcher.AND,
	)
	searchRequest.Limit = 1

	sqlStr, values, err := searcher.GenerateSelectSQL(&searchRequest, false)

	if err != nil {
		return nil, fmt.Errorf("Error generating GetByField SQL for Thing [db field:%s]: %s", field.DbFieldName(), err.Error())
	}

	var result types.Thing
	err = r.db.Get(&result, sqlStr, values.([]interface{})...)

	if err != nil {
		return nil, err
	}

	return &result, err
}
