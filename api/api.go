// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package bartapi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"code.google.com/p/go-charset/charset"

	// for the charset package we need to load
	// the data in for it to use.
	_ "code.google.com/p/go-charset/data"
)

const (
	// PublicAPIKey is the public key provided by BART for unregistered
	// use of their API. To note, by using this key you automatically
	// agree to their license agreement.
	PublicAPIKey = "MW9S-E7SL-26DU-VV8V"

	// URL is the base URL for the API endpoint
	URL = "http://api.bart.gov/api/bsa.aspx"
)

// Client is the BART API client
type Client struct {
	key, baseURL string
}

// New returns a new BART API client.
func New(key string) *Client {
	return &Client{key: key, baseURL: URL}
}

// SetBaseURL sets the base URL for the API client.
// "Base" meaning where the query params are passed.
func (c *Client) SetBaseURL(u string) {
	c.baseURL = u
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Key returns the API key of the client.
func (c *Client) Key() string {
	return c.key
}

// Pull does an HTTP GET request against the API endpoint.
// You need to provide the command (cmd) to send the API.
// You can add more query params using the "query" map
// if you need to, otherwise use nil.
func (c *Client) Pull(cmd string, query map[string]string) ([]byte, error) {
	var params bytes.Buffer

	params.WriteString(fmt.Sprintf("%v?cmd=%v&key=%v", c.baseURL, cmd, c.key))

	for k, v := range query {
		params.WriteString(fmt.Sprintf("&%v=%v", k, v))
	}

	resp, err := http.Get(params.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// Decode is a function to help with decoding the XML provided by BART.
// Because of their encoding format, we need to set the CharsetReader in
// this function. r is the data to parse, and v is data structure
// to parse it in to.
func Decode(r io.Reader, v interface{}) error {
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReader
	return d.Decode(v)
}
