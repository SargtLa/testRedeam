package api

import (
	"go/types"
	"runtime/trace"

	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

var (
	ParamBookID = apis.InParam{
		Name:     "id",
		Desc:     "id of book",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int64),
	}
)

var (
	Routes     = apis.NewMapRoutes()
	RoutesList = apis.ApiRoutes{
		ADD_BOOK_ROUTE: {
			Fnc:      HandleAddBook,
			Desc:     "adds new book to database",
			DTO:      &BookDTO{},
			NeedAuth: false,
			Method:   apis.POST,
		},
		EDIT_BOOK_ROUTE: {
			Fnc:      HandleEditBook,
			Desc:     "edit book",
			DTO:      &BookDTO{},
			NeedAuth: false,
			Method:   apis.POST,
			Params: []apis.InParam{
				ParamBookID,
			},
		},
		DELETE_BOOK_ROUTE: {
			Fnc:      HandleDeleteBook,
			Desc:     "delete book",
			NeedAuth: false,
			Method:   apis.POST,
			Params: []apis.InParam{
				ParamBookID,
			},
		},
		GET_ALL_BOOKS_ROUTE: {
			Fnc:      HandleAllBooks,
			Desc:     "get all books",
			NeedAuth: false,
			Method:   apis.GET,
		},
		GET_BOOK_ROUTE: {
			Fnc:      HandleGetBook,
			Desc:     "get one book",
			NeedAuth: false,
			Method:   apis.GET,
			Params: []apis.InParam{
				ParamBookID,
			},
		},
	}
)

func init() {
	for path, route := range RoutesList {
		if trace.IsEnabled() {
			route.Fnc = func(handler apis.ApiRouteHandler) apis.ApiRouteHandler {
				return func(ctx *fasthttp.RequestCtx) (resp interface{}, err error) {
					taskCtx, task := trace.NewTask(ctx, path)
					defer task.End()
					reg := trace.StartRegion(taskCtx, path)
					defer reg.End()
					logs.DebugLog(reg, task)
					trace.WithRegion(taskCtx, path,
						func() {
							resp, err = handler(ctx)
						})

					return
				}
			}(route.Fnc)
		}
	}
	Routes.AddRoutes(RoutesList)
}
