// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot.
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package tgfiles

import (
	"github.com/amarnathcjd/gogram/telegram"
)

type FileInfo struct {
	ID       int32
	Name     string
	MimeType string
	Size     int64
}

func GetFileInfo(message *telegram.NewMessage) FileInfo {
	media := message.Media()

	return FileInfo{
		ID:       message.ID,
		Name:     GetFileName(message, message.SenderID()),
		MimeType: GetMimeType(media),
		Size:     GetFileSize(media),
	}
}
