/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Rob Miller (rmiller@mozilla.com)
#
# ***** END LICENSE BLOCK *****/
package raven

import (
	"bytes"
	"code.google.com/p/gomock/gomock"
	"fmt"
	gs "github.com/rafrombrc/gospec/src/gospec"
	ts "heka/testsupport"
	"net/http"
	"net/url"
	"path"
	"time"
)

func build_request(timestamp time.Time, pkey string, project string, url *url.URL, packet []byte) *http.Request {
	apiURL := url
	apiURL.Path = path.Join(apiURL.Path, "/api/"+project+"/store/")
	apiURL.User = nil
	location := apiURL.String()

	buf := bytes.NewBuffer(packet)
	req, _ := http.NewRequest("POST", location, buf)

	authHeader := AuthHeader(timestamp, pkey)
	req.Header.Add("X-Sentry-Auth", authHeader)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Connection", "close")
	req.Header.Add("Accept-Encoding", "identity")

	return req
}

func RavenSpec(c gs.Context) {

	t := new(ts.SimpleT)
	ctrl := gomock.NewController(t)

	c.Specify("udp transport writes to a network connection", func() {
		dsn := "udp://someuser:somepass@localhost:801/2"
		client, _ := NewClient(dsn)

		udp_transport := client.sentryTransport.(*UdpSentryTransport)

		origClient := udp_transport.Client

		// Clobber the client with a mock network connection
		mock_conn := ts.NewMockConn(ctrl)
		udp_transport.Client = mock_conn
		defer func() {
			udp_transport := client.sentryTransport.(*UdpSentryTransport)
			udp_transport.Client = origClient
		}()

		timestamp := time.Now().UTC()

		str_packet := "some-data-string"
		expected_msg := []byte(fmt.Sprintf(UDP_TEMPLATE,
			AuthHeader(timestamp, udp_transport.PublicKey),
			str_packet))

		mock_conn.EXPECT().Write(expected_msg)
		mock_conn.EXPECT().Close()
		client.sentryTransport.Send([]byte(str_packet), timestamp)

	})

	c.Specify("http transport works on simple POST, no redirects", func() {
        // I"ve skipped this test until I can sort out the gomock
        // issue with HTTP Client
		return
		/*
			dsn := "http://someuser:somepass@localhost:801/2"
			client, _ := NewClient(dsn)

			http_transport := client.sentryTransport.(*HttpSentryTransport)

			origClient := http_transport.Client

			// Clobber the client with a mock network connection
			mock := NewMockHttpClient(ctrl)
			http_transport.Client = mock
			defer func() {
				http_transport := client.sentryTransport.(*HttpSentryTransport)
				http_transport.Client = origClient
			}()

			timestamp := time.Now().UTC()

			str_packet := "some-data-string"

			req := build_request(timestamp, client.PublicKey,
			client.Project, client.URL, []byte(str_packet))

			mock.EXPECT().Do(gomock.Any())

			client.sentryTransport.Send([]byte(str_packet), timestamp)
		*/
	})
}
