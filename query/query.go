package query

import (
	"github.com/asdine/genji/record"
	"github.com/asdine/genji/table"
)

type Query struct {
	selectors []FieldSelector
	matchers  []Matcher
}

func Select(selectors ...FieldSelector) Query {
	return Query{selectors: selectors}
}

type FieldSelector interface {
	Name() string
}

func (q Query) Run(t table.Reader) (table.Reader, error) {
	var rb table.RecordBuffer

	err := table.NewBrowser(t).ForEach(func(r record.Record) error {
		var fb record.FieldBuffer

		for _, s := range q.selectors {
			f, err := r.Field(s.Name())
			if err != nil {
				return err
			}

			fb.Add(f)
		}

		rb.Add(&fb)
		return nil
	}).Err()

	if err != nil {
		return nil, err
	}

	return &rb, nil
}

type Matcher interface {
	Match(record.Record) (bool, error)
}

func (q Query) Where(matchers ...Matcher) Query {
	q.matchers = append(q.matchers, matchers...)
	return q
}