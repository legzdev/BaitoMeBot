// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package tgfiles

import (
	"mime"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func GetMimeType(media telegram.MessageMedia) string {
	fileExt := GetFileExtension(media)

	mimeType := mime.TypeByExtension(fileExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType
}

func GetFileExtension(media telegram.MessageMedia) string {
	fileName := telegram.GetFileName(media)

	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex != -1 {
		return fileName[lastDotIndex:]
	}

	return ".unknown"
}
