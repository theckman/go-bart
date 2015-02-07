// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package bartapi_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theckman/go-bart/api"
	. "gopkg.in/check.v1"
)

var exampleXml = `
<root>
	<somekey>hello!</somekey>
</root>
`

type xmlType struct {
	XMLName xml.Name `xml:"root"`
	Some    string   `xml:"somekey"`
}

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	srv *httptest.Server
	url string
	c   *bartapi.Client
}

var _ = Suite(&TestSuite{})

func (t *TestSuite) SetUpTest(c *C) {
	h := &handler{}
	t.srv = httptest.NewServer(h)
	t.url = t.srv.URL
	t.c = bartapi.New("testkey")
}

func (t *TestSuite) TearDownTest(c *C) {
	t.srv.Close()
}

func (t *TestSuite) TestSetBaseURL(c *C) {
	t.c.SetBaseURL("http://localhost")
	url := t.c.BaseURL()
	c.Check(url, Equals, "http://localhost")
}

func (t *TestSuite) TestKey(c *C) {
	k := "madness"
	cl := bartapi.New(k)
	c.Check(cl.Key(), Equals, k)
}

func (t *TestSuite) TestPull(c *C) {
	c.Assert(t.c.Key(), Equals, "testkey")

	t.c.SetBaseURL(fmt.Sprintf("%v/", t.url))
	c.Assert(t.c.BaseURL(), Equals, fmt.Sprintf("%v/", t.url))

	resp, err := t.c.Pull("test", nil)
	c.Assert(err, IsNil)

	var j map[string]interface{}

	err = json.Unmarshal(resp, &j)
	c.Assert(err, IsNil)

	// need to type assert the value from an interface{} to a string
	c.Check((j["key"]).(string), Equals, "testkey")
	c.Check((j["cmd"]).(string), Equals, "test")

	params := make(map[string]string)
	params["bacon"] = "good"
	params["salad"] = "bad"

	resp, err = t.c.Pull("test2", params)
	c.Assert(err, IsNil)

	j = make(map[string]interface{})

	err = json.Unmarshal(resp, &j)
	c.Assert(err, IsNil)
	c.Check((j["key"]).(string), Equals, "testkey")
	c.Check((j["cmd"]).(string), Equals, "test2")
	c.Check((j["bacon"]).(string), Equals, "good")
	c.Check((j["salad"]).(string), Equals, "bad")
}

func (t *TestSuite) TestDecode(c *C) {
	r := bytes.NewReader([]byte(exampleXml))
	x := &xmlType{}

	err := bartapi.Decode(r, x)
	c.Assert(err, IsNil)
	c.Check(x.Some, Equals, "hello!")

	r = bytes.NewReader([]byte(""))
	x = &xmlType{}

	err = bartapi.Decode(r, x)
	c.Assert(err, Not(IsNil))
}

type handler struct{}

func (*handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	params := make(map[string]string)

	for k, v := range req.Form {
		params[k] = v[0]
	}

	resp, err := json.Marshal(params)

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(rw, string(resp))
}
