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
	"errors"
	"fmt"
	"log"

	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/db"
	"github.com/legzdev/BaitoMeBot/tgfiles"
)

func (bot *Bot) GenerateLink(ctx context.Context, fromPeer *tg.InputPeerUser, msgID int) (string, error) {
	toPeer, err := db.Peers.GetChannel(config.TelegramChatID)
	if err != nil {
		return "", err
	}

	req := &tg.MessagesForwardMessagesRequest{
		FromPeer: &tg.InputPeerUser{
			AccessHash: fromPeer.AccessHash,
			UserID:     fromPeer.UserID,
		},
		ToPeer: toPeer,
		ID:     []int{msgID},
	}

	update, err := bot.ForwardMessage(ctx, req)
	if err != nil {
		return "", err
	}

	updates, ok := update.(*tg.Updates)
	if !ok || len(updates.Updates) < 2 {
		return "", fmt.Errorf("link: invalid forward updates length %d", len(updates.Updates))
	}

	log.Printf("updates %#v\n", updates)

	storedMessageUpdate, okUpd := updates.Updates[1].(*tg.UpdateNewChannelMessage)
	storedMessage, okMsg := storedMessageUpdate.Message.(*tg.Message)

	if !okUpd || !okMsg {
		return "", errors.New("link: invalid message id or message type")
	}

	info, err := tgfiles.GetFileInfo(storedMessage, fromPeer.UserID)
	if err != nil {
		return "", err
	}

	link := config.ServerHost.JoinPath("dl")
	link = link.JoinPath(fmt.Sprint(storedMessage.ID))
	link = link.JoinPath(info.Name)

	query := link.Query()
	query.Set("hash", tgfiles.GetShortHash(info))
	link.RawQuery = query.Encode()

	return link.String(), err
}
