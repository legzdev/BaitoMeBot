// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package server_test

import (
	"testing"

	"github.com/legzdev/BaitoMeBot/server"
)

type RangeHeaderData struct {
	Header   string
	FileSize int64
	Range    server.Range
}

func TestParseRangeHeader(t *testing.T) {
	tests := []RangeHeaderData{
		{
			Header:   "",
			FileSize: 500,
			Range:    server.Range{Start: 0, End: 499},
		},
		{
			Header:   "bytes=900-",
			FileSize: 1000,
			Range:    server.Range{Start: 900, End: 999},
		},
		{
			Header:   "bytes=-100",
			FileSize: 500,
			Range:    server.Range{Start: 0, End: 100},
		},
		{
			Header:   "bytes=250-",
			FileSize: 500,
			Range:    server.Range{Start: 250, End: 499},
		},
	}

	for _, data := range tests {
		r, err := server.ParseRangeHeader(data.Header, data.FileSize)
		if err != nil {
			t.Fatal(err)
		}

		if r.Start != data.Range.Start || r.End != data.Range.End {
			t.Errorf("ParseRangeHeader(%q, %d): (%d, %d)", data.Header, data.FileSize, r.Start, r.End)
		}
	}
}
