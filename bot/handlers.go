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
	"context"
	"log"

	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/db"
)

// func HelpHandler(message *telegram.NewMessage) error {
// 	text := "Hello, send me any file to get a direct streamble link to that file."
// 	_, err := message.Reply(text)
// 	return err
// }
//
//

func (bot *Bot) MessageHandler() tg.NewMessageHandler {
	callback := func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		message, ok := update.Message.(*tg.Message)
		if !ok || message.Out {
			return nil
		}

		messagePeer, ok := message.PeerID.(*tg.PeerUser)
		if !ok {
			return nil
		}

		peer := &tg.InputPeerUser{
			UserID:     messagePeer.UserID,
			AccessHash: e.Users[messagePeer.UserID].AccessHash,
		}

		err := db.Peers.Put(peer)
		if err != nil {
			return err
		}

		_, isMedia := message.GetMedia()
		if isMedia {
			link, err := bot.GenerateLink(ctx, peer, message.ID)
			if err != nil {
				return err
			}

			log.Println(link)

			req := &tg.MessagesSendMessageRequest{
				Peer:    peer,
				Message: link,
			}

			_, err = bot.SendMessage(ctx, req)
			if err != nil {
				return err
			}

		} else {
		}

		return nil

		// _, ok = message.GetMedia()
		// if ok {
		// 	opts := &tg.MessagesForwardMessagesRequest{}
		//
		// 	forward, err := client.API().MessagesForwardMessages(ctx, opts)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	return func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		err := callback(ctx, e, update)
		if err != nil {
			log.Printf("ERRO: %v", err)
		}
		return err
	}
}
