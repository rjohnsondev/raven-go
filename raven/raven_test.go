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
	gs "github.com/rafrombrc/gospec/src/gospec"
)

func RavenSpec(c gs.Context) {
	c.Specify("sending udp", func() {
		dsn := "udp://59a868d4c312450ba46074dc48ae4bc5:aed77a1e07514208b90d43451083b9d4@localhost:9001/2"
		client, _ := NewClient(dsn)

		client.CaptureMessage("blah blah blah", "some other text",
			"how do exceptions actually work in this thing?")
		// uh.. i need to inject a mock

	})
}
