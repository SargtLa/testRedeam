package api

import (
	"fmt"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
	"strings"

	"github.com/SargtLa/testRedeam/db"
)

type BookDTO struct {
	*db.BookFields
}

func (c *BookDTO) GetValue() interface{} {
	return c
}

func (c *BookDTO) NewValue() interface{} {
	return &BookDTO{BookFields: &db.BookFields{}}
}

func HandleAddBook(ctx *fasthttp.RequestCtx) (interface{}, error) {
	book, ok := ctx.UserValue(apis.JSONParams).(*BookDTO)
	if !ok {
		return "wrong BookDTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, err := db.NewBooksTable(DB)
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)

	for _, col := range table.Columns() {
		name := col.Name()
		newValue := book.ColValue(name)
		if !EmptyValue(newValue) {
			columns = append(columns, name)
			args = append(args, newValue)
		}
	}

	id, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return nil, nil
		}
		return nil, err

	}

	book.Id = id
	ctx.SetStatusCode(fasthttp.StatusCreated)

	return id, nil
}

func HandleEditBook(ctx *fasthttp.RequestCtx) (interface{}, error) {
	book, ok := ctx.UserValue(apis.JSONParams).(*BookDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamBookID.Name).(int64)
	if !ok {
		return map[string]string{
			ParamBookID.Name: fmt.Sprintf("wrong type %T, expect int64 ",
				ctx.UserValue(ParamBookID.Name)),
		}, apis.ErrWrongParamsList
	}

	book.Id = id

	table, err := db.NewBooksTable(DB)
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)

	for _, col := range table.Columns() {
		name := col.Name()
		newValue := book.ColValue(name)
		if !EmptyValue(newValue) {
			columns = append(columns, name)
			args = append(args, newValue)
		}
	}

	if len(columns) == 0 {
		return "nothing to update", apis.ErrWrongParamsList
	}

	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(append(args, book.Id)...),
	)
	if err != nil {
		return nil, err
	}

	if i > 0 {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	}

	return nil, nil

}

func HandleDeleteBook(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamBookID.Name).(int64)
	if !ok {
		return map[string]string{
			ParamBookID.Name: fmt.Sprintf("wrong type %T, expect int64 ",
				ctx.UserValue(ParamBookID.Name)),
		}, apis.ErrWrongParamsList
	}

	err := DB.Conn.ExecDDL(ctx, fmt.Sprintf(DELETE_BY_ID_SQL, db.TABLE_BOOKS), id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func HandleAllBooks(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, err := db.NewBooksTable(DB)
	if err != nil {
		return nil, err
	}

	var allBooks []db.BookFields

	err = table.SelectSelfScanEach(ctx,
		func(record *db.BookFields) error {
			allBooks = append(allBooks, *record)
			return nil
		},
	)

	if err != nil {
		logs.ErrorLog(err, "while reading all books")
	}

	return allBooks, nil
}

func HandleGetBook(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamBookID.Name).(int64)
	if !ok {
		return map[string]string{
			ParamBookID.Name: fmt.Sprintf("wrong type %T, expect int64 ",
				ctx.UserValue(ParamBookID.Name)),
		}, apis.ErrWrongParamsList
	}

	table, _ := db.NewBooksTable(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return nil, err
	}

	return *table.Record, nil
}
