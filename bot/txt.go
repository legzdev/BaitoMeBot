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
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/db"
)

func OnTxtCommand(message *telegram.NewMessage) error {
	command := message.GetCommand()

	if strings.HasPrefix(command, "/txtcaptionfull") {
		return onTxt(message, db.StateTxtCaptionFull)
	} else if strings.HasPrefix(command, "/txtcaption") {
		return onTxt(message, db.StateTxtCaption)
	}

	return onTxt(message, db.StateTxt)
}

func onTxt(message *telegram.NewMessage, state db.State) error {
	senderID := message.SenderID()

	buffer := db.GetBuffer(senderID)
	if buffer == nil {
		bufferName := NameFromArgs(message.MessageText())
		bufferName = GetBufferName("", bufferName)

		db.SetBuffer(senderID, bufferName)
		db.SetState(senderID, state)

		text := "<b>Switched to TXT (buffer) mode.</b>\n"
		text += fmt.Sprintf("📝 <b>Name:</b> <code>%s</code>\n", bufferName)
		text += fmt.Sprintf("🔭 <b>Mode:</b> %s", state.String())

		opts := telegram.SendOptions{
			ParseMode: telegram.HTML,
		}

		_, err := message.Reply(text, opts)
		return err
	}

	if buffer.Len() == 0 {
		db.DelBuffer(senderID)
		db.DelState(senderID)

		text := "<b>Switched to normal mode.</b>\n"
		text += "Cleaned empty buffer"

		opts := telegram.SendOptions{
			ParseMode: telegram.HTML,
		}

		_, err := message.Reply(text, opts)
		return err
	}

	newName := NameFromArgs(message.MessageText())
	fileName := GetBufferName(buffer.Name, newName)

	opts := telegram.MediaOptions{
		FileName: fileName,
	}

	_, err := message.ReplyMedia(buffer.Bytes(), opts)
	if err != nil {
		return err
	}

	db.DelBuffer(senderID)
	db.DelState(senderID)

	return nil
}

func NameFromArgs(text string) string {
	textParts := strings.Split(text, " ")
	if len(textParts) < 2 {
		return ""
	}

	return strings.Join(textParts, " ")
}

func GetBufferName(currentName string, newName string) string {
	var name string

	if newName != "" {
		name = newName
	} else if currentName != "" {
		name = currentName
	} else {
		name = time.Now().Format(time.DateTime) + ".txt"
	}

	if !strings.HasSuffix(name, ".txt") {
		name += ".txt"
	}

	return name
}
