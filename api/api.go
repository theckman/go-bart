// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// Package bartapi is for polling and XML decoding the Bay Area Rapid Transit
// API endpoints. The struct types to decode the XML in to are not provided
// by bartapi. Consumers of bartapi maintain their own structs.
//
// NOTE: This package is meant to be consumed as part of go-bart.
// If you are looking to interact with the BART API you should use that instead.
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

// Endpoint is a string which contains the
// URL of a specific BART API endpoint.
type Endpoint string

// PublicAPIKey is the public key provided by BART for unregistered
// use of their API. To note, by using this key you automatically
// agree to their license agreement.
const PublicAPIKey = "MW9S-E7SL-26DU-VV8V"

// AdvisoryEndpoint is the endpoint for requesting BART Service Advisories.
//
// This endpoint handles commands to get any current advisories (delays),
// any elevator status information, as well as the number of trains active
// in the system.
const AdvisoryEndpoint Endpoint = "http://api.bart.gov/api/bsa.aspx"

// EstimatesEndpoint is the endpoint for getting Real-Time estimates for
// departure times at a given station for specific routes.
const EstimatesEndpoint Endpoint = "http://api.bart.gov/api/etd.aspx"

// RouteEndpoint is the endpoint for getting information about specific routes.
// This includes pulling information about a single route, or information
// on all BART routes.
const RouteEndpoint Endpoint = "http://api.bart.gov/api/route.aspx"

// ScheduleEndpoint is the endpoint for getting information about schedules,
// and things BART thinks are related to schedules (like fare).
//
// This includes: trip planning based on arrival and departure times, estimated
// load factor for a given trains, fare calculation, full route schedules,
// station schedules, and holiday / special schedule notices.
const ScheduleEndpoint Endpoint = "http://api.bart.gov/api/sched.aspx"

// StationEndpoint is the endpoint for getting station information.
// This includes a list of all stations, station general information,
// and inforation about access to and from the station as well as surrounding
// neighborhood information.
const StationEndpoint Endpoint = "http://api.bart.gov/api/stn.aspx"

// Client is the BART API client
type Client struct {
	key string
	url Endpoint
}

// New returns a new BART API client.
func New(key string, url Endpoint) *Client {
	return &Client{key: key, url: url}
}

// URL returns the endpoint being used by the client.
func (c *Client) URL() Endpoint {
	return c.url
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

	params.WriteString(fmt.Sprintf("%v?cmd=%v&key=%v", string(c.url), cmd, c.key))

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
