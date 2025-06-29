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
	"io"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/errs"
)

type TelegramReader struct {
	io.ReadSeeker
	client        *telegram.Client
	media         telegram.MessageMedia
	location      telegram.InputFileLocation
	dcID          int
	currentOffset int64
	maxOffset     int64
	end           int64
}

func NewTelegramReader(client *telegram.Client, media telegram.MessageMedia, end int64) (*TelegramReader, error) {
	location, dcID, _, _, err := telegram.GetFileLocation(media)
	if err != nil {
		return nil, err
	}

	return &TelegramReader{
		client:    client,
		media:     media,
		location:  location,
		dcID:      int(dcID),
		end:       end,
		maxOffset: end + 1,
	}, nil
}

func (r *TelegramReader) Read(buffer []byte) (int, error) {
	if r.currentOffset == r.maxOffset {
		return 0, io.EOF
	}

	if r.currentOffset > r.maxOffset {
		return 0, fmt.Errorf("invalid offset %d (end %d)", r.currentOffset, r.end)
	}

	var paddingPrefix int
	var paddingSuffix int

	bufferSize := len(buffer)
	if (bufferSize % config.AlignSize) != 0 {
		return 0, fmt.Errorf("buffer length should be divisible by %d", config.AlignSize)
	}

	offset := r.currentOffset
	limit := bufferSize

	if (offset % int64(config.AlignSize)) != 0 {
		align := offset / config.AlignSize
		newOffset := align * config.AlignSize

		paddingPrefix = int(offset - newOffset)
		paddingSuffix = config.AlignSize - paddingPrefix

		offset = newOffset
		limit = paddingPrefix + bufferSize + paddingSuffix
	}

	sender := r.client.MTProto
	var err error

	if sender.GetDC() != r.dcID {
		sender, err = r.client.CreateExportedSender(r.dcID, false)
		if err != nil {
			return 0, err
		}
	}

	params := &telegram.UploadGetFileParams{
		Location:     r.location,
		Offset:       offset,
		Limit:        int32(limit),
		Precise:      true,
		CdnSupported: false,
	}

	part, err := sender.MakeRequest(params)
	if err != nil {
		return 0, errs.Wrapf(err, "request error (off=%d lim=%d end=%d)", r.currentOffset, limit, r.end)
	}

	switch res := part.(type) {
	case *telegram.UploadFileObj:
		resSize := bufferSize

		resLen := len(res.Bytes)
		if resSize > resLen {
			resSize = resLen
		}

		if (r.currentOffset + int64(resSize)) > r.maxOffset {
			resSize = int(r.maxOffset - r.currentOffset)
		}

		resEnd := paddingPrefix + resSize
		resBuf := res.Bytes[paddingPrefix:resEnd]

		written := copy(buffer, resBuf)
		r.currentOffset += int64(written)

		return written, nil

	case *telegram.UploadFileCdnRedirect:
		return 0, errors.New("cdn redirect not implemented")
	}

	return 0, errors.New("unexpected result")
}

func (r *TelegramReader) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64

	switch whence {
	case io.SeekStart:
		newOffset = offset

	case io.SeekCurrent:
		newOffset = r.currentOffset + offset

	case io.SeekEnd:
		newOffset = r.end + offset

	default:
		return 0, fmt.Errorf("seek: invalid whence %d", whence)
	}

	if newOffset > r.maxOffset {
		newOffset = r.maxOffset
	}

	r.currentOffset = newOffset
	return r.currentOffset, nil
}
