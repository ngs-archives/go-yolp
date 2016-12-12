package yolp

import (
	"encoding/xml"
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

type Test struct {
	expected interface{}
	actual   interface{}
}

func (test Test) Compare(t *testing.T) {
	if test.expected != test.actual {
		t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
	}
}

func (test Test) DeepEqual(t *testing.T) {
	if !reflect.DeepEqual(test.expected, test.actual) {
		t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
	}
}

func TestNew(t *testing.T) {
	client, _ := New("APPID", "SECRET")
	for _, test := range []Test{
		{"APPID", client.AppID},
		{"SECRET", client.Secret},
	} {
		test.Compare(t)
	}
}

func TestNewEmptyAppID(t *testing.T) {
	client, err := New("", "SECRET")
	Test{"AppID is not specified", err.Error()}.Compare(t)
	if client != nil {
		t.Errorf(`Expected nil but got "%v"`, client)
	}
}

func TestNewEmptySecret(t *testing.T) {
	client, err := New("APPID", "")
	Test{"Secret is not specified", err.Error()}.Compare(t)
	if client != nil {
		t.Errorf(`Expected nil but got "%v"`, client)
	}
}

func TestNewFromEnvionment(t *testing.T) {
	os.Setenv("YDN_APP_ID", "APPID")
	os.Setenv("YDN_SECRET", "SECRET")
	client, _ := NewFromEnvionment()
	for _, test := range []Test{
		{"APPID", client.AppID},
		{"SECRET", client.Secret},
	} {
		test.Compare(t)
	}
}

type mockRequest struct {
	method string
}

func (mop mockRequest) Query() url.Values {
	return url.Values{
		"foo": []string{"bar"},
	}
}

func (mop mockRequest) HTTPMethod() string {
	if mop.method == "" {
		return "GET"
	}
	return mop.method
}

func (mop mockRequest) Endpoint() string {
	return "https://foo.com/bar/baz"
}

type mockResponse struct {
	XMLName xml.Name `xml:"mock"`
	Result  string   `xml:"result"`
}

func TestClientURL(t *testing.T) {
	client, _ := New("APP", "SEC")
	op := mockRequest{}
	Test{
		"https://foo.com/bar/baz?appid=APP&foo=bar&output=xml",
		client.URL(op).String(),
	}.Compare(t)
}

func TestDoGetRequest(t *testing.T) {
	defer gock.Off()
	gock.DisableNetworking()
	gock.New("https://foo.com/bar/baz?appid=APP&foo=bar&output=xml").
		Reply(200).
		BodyString("<mock><result>OK</result></mock>")
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	Test{200, res.StatusCode}.Compare(t)
	Test{"OK", mockResp.Result}.Compare(t)
}

func TestDoPostRequest(t *testing.T) {
	defer gock.Off()
	gock.DisableNetworking()
	gock.New("https://foo.com").
		Post("/bar/baz").
		BodyString("appid=APP&foo=bar&output=xml").
		Reply(200).
		BodyString("<mock><result>OK</result></mock>")
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{
		method: "POST",
	}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	Test{200, res.StatusCode}.Compare(t)
	Test{"OK", mockResp.Result}.Compare(t)
}

func TestDoInvalidMethodRequest(t *testing.T) {
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{
		method: "DELETE",
	}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	Test{"Unsupported HTTP method: DELETE", err.Error()}.Compare(t)
	if res != nil {
		t.Errorf("Expected nil but got %v", res)
	}
}

func TestDoHTTPError(t *testing.T) {
	gock.New("https://foo.com/bar/baz?appid=APP&foo=bar&output=xml").
		ReplyError(errors.New("oops"))
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	Test{
		"Get https://foo.com/bar/baz?appid=APP&foo=bar&output=xml: oops",
		err.Error()}.Compare(t)
	if res != nil {
		t.Errorf("Expected nil but got %v", res)
	}
}

func TestDoInvalidXML(t *testing.T) {
	gock.New("https://foo.com/bar/baz?foo=bar&output=xml").
		Reply(200).
		BodyString("<invalidmock><result>OK</result></invalidmock>")
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	Test{
		"expected element type <mock> but have <invalidmock>",
		err.Error()}.Compare(t)
	if res != nil {
		t.Errorf("Expected nil but got %v", res)
	}
}

func TestDoErrorResponse(t *testing.T) {
	gock.New("https://foo.com/bar/baz?foo=bar&output=xml").
		Reply(200).
		File("_fixtures/error.xml")
	client, _ := New("APP", "SEC")
	mockOp := &mockRequest{}
	mockResp := mockResponse{}
	res, err := client.DoRequest(mockOp, &mockResp)
	Test{
		"Your Request was Forbidden",
		err.Error()}.Compare(t)
	if res != nil {
		t.Errorf("Expected nil but got %v", res)
	}
}
