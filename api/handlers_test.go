package api

import (
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

var flagTestPort = flag.String("test_port", ":8080", "Test server listening port")

type testBookStruct struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int32     `json:"rating"`
	Status      bool      `json:"status"`
}

// We may no be sure about the tests order, so each test func should have its own test values
// to avoid conflicts between test.
var (
	testAddValues = testBookStruct{
		Title:     "Harry Potter and the Philosopher's Stone",
		Author:    "J. K. Rowling",
		Publisher: "Bloomsbury Children's",
		PublishDate: time.Date(2014, 1, 1,
			0, 0, 0, 0, time.UTC),
		Rating: 3,
		Status: true,
	}

	testDeleteValues = testBookStruct{
		Title:     "Harry Potter and the Chamber of Secrets",
		Author:    "J. K. Rowling",
		Publisher: "Bloomsbury Children's",
		PublishDate: time.Date(2014, 1, 1,
			0, 0, 0, 0, time.UTC),
		Rating: 3,
		Status: true,
	}

	testEditValues = testBookStruct{
		Title:     "Harry Potter and the Prisoner of Azkaban",
		Author:    "J. K. Rowling",
		Publisher: "Bloomsbury Children's",
		PublishDate: time.Date(2014, 1, 1,
			0, 0, 0, 0, time.UTC),
		Rating: 3,
		Status: true,
	}

	testGetValues = testBookStruct{
		Title:     "Harry Potter and the Goblet of Fire",
		Author:    "J. K. Rowling",
		Publisher: "Bloomsbury Children's",
		PublishDate: time.Date(2014, 1, 1,
			0, 0, 0, 0, time.UTC),
		Rating: 3,
		Status: true,
	}
)

func TestHandleAddBook(t *testing.T) {
	as := assert.New(t)

	testValuesJSON, err := jsoniter.Marshal(testAddValues)

	as.Nil(err)

	req := &fasthttp.Request{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s", *flagTestPort, ADD_BOOK_ROUTE))
	req.SetBody(testValuesJSON)

	resp, err := DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusCreated, resp.StatusCode())
	bookId := resp.Body()

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusConflict, resp.StatusCode())

	// to keep database the same as it was we must delete values
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, DELETE_BOOK_ROUTE, bookId))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())
}

func TestHandleDeleteBook(t *testing.T) {
	as := assert.New(t)

	testDeleteValuesJSON, err := jsoniter.Marshal(testDeleteValues)

	as.Nil(err)

	req := &fasthttp.Request{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s", *flagTestPort, ADD_BOOK_ROUTE))
	req.SetBody(testDeleteValuesJSON)

	resp, err := DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusCreated, resp.StatusCode())

	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, DELETE_BOOK_ROUTE, resp.Body()))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())

}

func TestHandleEditBook(t *testing.T) {
	as := assert.New(t)

	testEditValuesJSON, err := jsoniter.Marshal(testEditValues)

	as.Nil(err)

	req := &fasthttp.Request{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s", *flagTestPort, ADD_BOOK_ROUTE))
	req.SetBody(testEditValuesJSON)

	resp, err := DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusCreated, resp.StatusCode())
	bookId := resp.Body()

	testNewEditValues := testEditValues
	testNewEditValues.Status = false
	testNewEditValuesJSON, err := jsoniter.Marshal(testNewEditValues)
	as.Nil(err)
	req.SetBody(testNewEditValuesJSON)

	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, EDIT_BOOK_ROUTE, bookId))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusAccepted, resp.StatusCode())

	// to keep database the same as it was we must delete values
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, DELETE_BOOK_ROUTE, bookId))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())

}

func TestHandleGetBook(t *testing.T) {
	as := assert.New(t)

	testGetValuesJSON, err := jsoniter.Marshal(testGetValues)

	as.Nil(err)

	req := &fasthttp.Request{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s", *flagTestPort, ADD_BOOK_ROUTE))
	req.SetBody(testGetValuesJSON)

	resp, err := DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusCreated, resp.StatusCode())
	bookId := resp.Body()

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, GET_BOOK_ROUTE, bookId))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())

	checkReturnValues := testBookStruct{}
	err = jsoniter.Unmarshal(resp.Body(), &checkReturnValues)
	as.Nil(err)

	as.EqualValues(testGetValues, checkReturnValues)

	// to keep database the same as it was we must delete values
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s/%s", *flagTestPort, DELETE_BOOK_ROUTE, bookId))

	resp, err = DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())
}

func TestGetAllBooks(t *testing.T) {
	as := assert.New(t)

	req := &fasthttp.Request{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(fmt.Sprintf("http://localhost%s%s", *flagTestPort, GET_ALL_BOOKS_ROUTE))

	resp, err := DoRequestForTest(req)
	as.Nil(err)

	as.EqualValues(fasthttp.StatusOK, resp.StatusCode())
}

func DoRequestForTest(req *fasthttp.Request) (resp fasthttp.Response, err error) {
	c := fasthttp.Client{
		MaxIdleConnDuration: time.Minute * 3,
		MaxConnDuration:     time.Minute * 3,
	}

	err = c.Do(req, &resp)

	return
}
