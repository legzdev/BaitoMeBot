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
	"github.com/legzdev/BaitoMeBot/db"
)

func OnMessage(message *telegram.NewMessage) error {
	if !message.IsMedia() {
		return nil
	}

	link, err := GenerateLink(message)
	if err != nil {
		return err
	}

	keyboard := telegram.NewKeyboard()
	keyboard.AddRow(
		telegram.Button.URL("Open", link),
	)

	text := fmt.Sprintf("<code>%s</code>", link)
	opts := &telegram.SendOptions{
		ParseMode:   telegram.HTML,
		ReplyMarkup: keyboard.Build(),
		ReplyID:     message.ID,
	}

	_, err = Bot.SendMessage(message.Client, message.ChannelID(), text, opts)
	if err != nil {
		return err
	}

	buffer := db.GetBuffer(message.SenderID())
	if buffer != nil {
		buffer.WriteString(link + "\n")
	}

	return nil
}

func HelpHandler(message *telegram.NewMessage) error {
	text := "Hello, send me any file to get a direct streamble link to that file."
	_, err := message.Reply(text)
	return err
}
