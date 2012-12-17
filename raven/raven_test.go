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
	"code.google.com/p/gomock/gomock"
	"fmt"
	gs "github.com/rafrombrc/gospec/src/gospec"
	ts "heka/testsupport"
	"time"
)

func RavenSpec(c gs.Context) {

	t := new(ts.SimpleT)
	ctrl := gomock.NewController(t)

	c.Specify("sending udp", func() {
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
}
