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
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/db"
)

func GetFileName(message *telegram.NewMessage, userID int64) string {
	var fileName string

	state := db.GetState(userID)
	if state == db.StateTxtCaption || state == db.StateTxtCaptionFull {
		caption := message.Text()

		if caption == "" {
			// do nothing
		} else if state == db.StateTxtCaption {
			captionLines := strings.Split(caption, "\n")
			fileName = strings.TrimSpace(captionLines[0]) + GetFileExtension(message.Media())

		} else if state == db.StateTxtCaptionFull {
			caption = strings.ReplaceAll(caption, "\n", "")
			fileName = caption + GetFileExtension(message.Media())
		}
	}

	if fileName == "" {
		fileName = telegram.GetFileName(message.Media())
	}

	return fileName
}
