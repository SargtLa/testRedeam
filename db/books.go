package db

import (
	"database/sql"
	"golang.org/x/net/context"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
)

type BooksTable struct {
	dbEngine.Table
	Record *BookFields
	rows   sql.Rows
}

type BookFields struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int32     `json:"rating"`
	Status      bool      `json:"status"`
}

func (r *BookFields) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		return []interface{}{
			&r.Id,
			&r.Title,
			&r.Author,
			&r.Publisher,
			&r.PublishDate,
			&r.Rating,
			&r.Status,
		}
	}

	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = r.RefColValue(col.Name())
	}

	return v
}

func (r *BookFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "title":
		return &r.Title

	case "author":
		return &r.Author

	case "publisher":
		return &r.Publisher

	case "publish_date":
		return &r.PublishDate

	case "rating":
		return &r.Rating

	case "status":
		return &r.Status

	default:
		return nil
	}
}

func (r *BookFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "title":
		return r.Title

	case "author":
		return r.Author

	case "publisher":
		return r.Publisher

	case "publish_date":
		return r.PublishDate

	case "rating":
		return r.Rating

	case "status":
		return r.Status

	default:
		return nil
	}
}

func NewBooksTable(db *dbEngine.DB) (*BooksTable, error) {
	table, ok := db.Tables[TABLE_BOOKS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_BOOKS}
	}

	return &BooksTable{
		Table: table,
	}, nil
}

func (t *BooksTable) NewRecord() *BookFields {
	t.Record = &BookFields{}
	return t.Record
}

func (t *BooksTable) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		columns = t.Columns()
	}

	t.NewRecord()
	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = t.Record.RefColValue(col.Name())
	}

	return v
}

func (t *BooksTable) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		//v := make([]interface{}, len(t.Columns()))
		//columns := make([]string, len(t.Columns()))
		//for i, col := range t.Columns() {
		//	columns[i] = col.Name()
		//	v[i] = t.Record.ColValue(col.Name())
		//}
		//Options = append(Options,
		//	dbEngine.ColumnsForSelect(columns...),
		//	dbEngine.ArgsForSelect(v...))

		return 0, nil
	}

	return t.Table.Insert(ctx, Options...)
}

func (t *BooksTable) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		//	v := make([]interface{}, len(t.Columns()))
		//	priV := make([]interface{}, 0)
		//	columns := make([]string, 0, len(t.Columns()))
		//	priColumns := make([]string, 0, len(t.Columns()))
		//	for _, col := range t.Columns() {
		//		if col.Primary() {
		//			priColumns = append(priColumns, col.Name())
		//			priV[len(priColumns)-1] = t.Record.ColValue(col.Name())
		//			continue
		//		}
		//
		//		columns = append(columns, col.Name())
		//		v[len(columns)-1] = t.Record.ColValue(col.Name())
		//	}
		//
		//	Options = append(
		//		Options,
		//		dbEngine.ColumnsForSelect(columns...),
		//		dbEngine.WhereForSelect(priColumns...),
		//		dbEngine.ArgsForSelect(append(v, priV...)...),
		//	)
		return 0, nil
	}

	return t.Table.Update(ctx, Options...)
}

func (t *BooksTable) SelectSelfScanEach(ctx context.Context, each func(record *BookFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}
