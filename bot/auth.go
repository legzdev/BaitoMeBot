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

// type MiddlewareFunc = func(*telegram.NewMessage) error
//
// func AuthMiddleware(callback MiddlewareFunc) MiddlewareFunc {
// 	return func(message *telegram.NewMessage) error {
// 		var whiteListEnabled bool
//
// 		client := message.Client
// 		senderID := message.SenderID()
//
// 		if config.AllowedUsers != nil {
// 			if slices.Contains(config.AllowedUsers, senderID) {
// 				return callback(message)
// 			}
// 			whiteListEnabled = true
// 		}
//
// 		if config.AuthChannelID != 0 {
// 			member, err := client.GetChatMember(config.AuthChannelID, senderID)
// 			if err != nil {
// 				return err
// 			}
//
// 			switch member.Status {
// 			case telegram.Creator:
// 			case telegram.Admin:
// 			case telegram.Member:
// 				return callback(message)
// 			}
//
// 			whiteListEnabled = true
// 		}
//
// 		if whiteListEnabled {
// 			text := "You are not allowed to use this bot.\n"
// 			text += "However, it's open-source and you can host it yourself."
//
// 			sourceButton := telegram.Button.URL(
// 				"Source Code", config.SourceURL,
// 			)
//
// 			keyboard := telegram.NewKeyboard()
// 			keyboard.AddRow(sourceButton)
//
// 			opts := telegram.SendOptions{
// 				ReplyMarkup: keyboard.Build(),
// 			}
//
// 			_, err := message.Reply(text, opts)
// 			return err
// 		}
//
// 		return callback(message)
// 	}
// }
