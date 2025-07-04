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
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/config"
)

type Mutex struct {
	sync.Mutex
	lastUsage time.Time
}

func (bot *Bot) NewTGMutex() {
	go func() {
		for {
			bot.mux.Lock()

			for mutexID, mutex := range bot.muxes {
				if time.Since(mutex.lastUsage) > config.TimeBetweenChecks {
					delete(bot.muxes, mutexID)
				}
			}

			bot.mux.Unlock()
			time.Sleep(config.TimeBetweenChecks)
		}
	}()
}

func (bot *Bot) GetChatMutex(chatID string) *Mutex {
	bot.mux.Lock()
	defer bot.mux.Unlock()

	userMutex, ok := bot.muxes[chatID]
	if !ok {
		userMutex = new(Mutex)
		bot.muxes[chatID] = userMutex
	} else {
		elapsed := time.Since(userMutex.lastUsage)
		if elapsed < config.TimeBetweenMessages {
			time.Sleep(config.TimeBetweenMessages - elapsed)
		}
	}

	return userMutex
}

func (bot *Bot) SendMessage(ctx context.Context, req *tg.MessagesSendMessageRequest) (upd tg.UpdatesClass, err error) {
	bot.do(peerKey(req.Peer), func() {
		req.RandomID = rand.Int63()
		upd, err = bot.Client.API().MessagesSendMessage(ctx, req)
	})
	return upd, err
}

func (bot *Bot) SendMedia(ctx context.Context, req *tg.MessagesSendMediaRequest) (upd tg.UpdatesClass, err error) {
	bot.do(peerKey(req.Peer), func() {
		req.RandomID = rand.Int63()
		upd, err = bot.Client.API().MessagesSendMedia(ctx, req)
	})
	return upd, err
}

func (bot *Bot) ForwardMessage(ctx context.Context, req *tg.MessagesForwardMessagesRequest) (upd tg.UpdatesClass, err error) {
	bot.do(peerKey(req.ToPeer), func() {
		req.RandomID = []int64{rand.Int63()}
		upd, err = bot.Client.API().MessagesForwardMessages(ctx, req)
	})
	return upd, err
}

func (bot *Bot) do(peerKey string, callback func()) {
	mux := bot.GetChatMutex(peerKey)
	mux.Lock()
	defer mux.Unlock()

	callback()
	mux.lastUsage = time.Now()
}

func peerKey(peer tg.InputPeerClass) string {
	switch v := peer.(type) {
	case *tg.InputPeerChannel:
		return fmt.Sprintf("channel(%d)", v.ChannelID)
	case *tg.InputPeerUser:
		return fmt.Sprintf("user(%d)", v.UserID)
	}
	return "unknown(0)"
}
