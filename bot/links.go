// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package bot

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/tgfiles"
)

func GenerateLink(message *telegram.NewMessage) (string, error) {
	storedMessage, err := message.ForwardTo(config.TelegramChatID)
	if err != nil {
		return "", err
	}

	info := tgfiles.GetFileInfo(storedMessage)
	fileID := fmt.Sprint(storedMessage.ID)
	hash := tgfiles.GetShortHash(info)

	link := config.ServerHost.JoinPath("dl", fileID, info.Name)

	query := link.Query()
	query.Set("hash", hash)
	link.RawQuery = query.Encode()

	return link.String(), err
}
