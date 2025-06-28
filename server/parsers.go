// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Range struct {
	Start        int64
	End          int64
	IsFromHeader bool
}

func ParseRequestRange(r *http.Request, fileSize int64) (Range, error) {
	return ParseRangeHeader(r.Header.Get("Range"), fileSize)
}

func ParseRangeHeader(header string, fileSize int64) (Range, error) {
	result := Range{End: fileSize - 1}
	if header == "" {
		return result, nil
	}

	result.IsFromHeader = true

	// Reject multi-range requests
	if strings.Contains(header, ",") {
		return result, errors.New("range: multi-range requests are not implemented")
	}

	headerParts := strings.Split(header, "=")
	if len(headerParts) != 2 {
		return result, errors.New("range: bad request")
	}

	unit := headerParts[0]
	if unit != "bytes" {
		return result, fmt.Errorf("invalid range unit %q", unit)
	}

	parts := strings.Split(headerParts[1], "-")
	if len(parts) != 2 {
		return result, errors.New("missing range")
	}

	var start int64
	var err error

	if parts[0] != "" {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return result, err
		}
	}

	var end int64

	if parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return result, err
		}
	}

	result.Start = start
	if end != 0 {
		result.End = end
	}

	if result.Start > result.End || result.End >= fileSize {
		return result, errors.New("invalid range")
	}

	return result, nil
}
