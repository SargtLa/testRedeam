package api

// routes
const (
	ADD_BOOK_ROUTE      = "/api/addBook"
	GET_BOOK_ROUTE      = "/api/main/getBook/"
	EDIT_BOOK_ROUTE     = "/api/editBook/"
	DELETE_BOOK_ROUTE   = "/api/deleteBook/"
	GET_ALL_BOOKS_ROUTE = "/api/getAllBooks"
)

// names of system environment variables
const (
	CFG_PATH    = "configPath"
	SYSTEM_PATH = "systemPath"
)

const DELETE_BY_ID_SQL = "DELETE FROM %s WHERE id = $1"
