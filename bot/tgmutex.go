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
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/config"
)

type TGMutex struct {
	mux   sync.Mutex
	muxes map[int64]*Mutex
}

type Mutex struct {
	sync.Mutex
	lastUsage time.Time
}

func NewTGMutex() *TGMutex {
	mux := &TGMutex{
		muxes: make(map[int64]*Mutex),
	}

	go func() {
		for {
			mux.mux.Lock()

			for userID, mutex := range mux.muxes {
				if time.Since(mutex.lastUsage) > config.TimeBetweenChecks {
					delete(mux.muxes, userID)
				}
			}

			mux.mux.Unlock()
			time.Sleep(config.TimeBetweenChecks)
		}
	}()

	return mux
}

func (tgm *TGMutex) GetChatMutex(chatID int64) *Mutex {
	tgm.mux.Lock()
	defer tgm.mux.Unlock()

	userMutex, ok := tgm.muxes[chatID]
	if !ok {
		userMutex = new(Mutex)
		tgm.muxes[chatID] = userMutex
	}

	return userMutex
}

func (tgm *TGMutex) SendMessage(c *telegram.Client, chatID int64, text string, opts ...*telegram.SendOptions) (*telegram.NewMessage, error) {
	mux := tgm.GetChatMutex(chatID)
	mux.Lock()
	defer mux.Unlock()

	message, err := c.SendMessage(chatID, text, opts...)
	time.Sleep(config.TimeBetweenMessages)

	return message, err
}

func (tgm *TGMutex) SendMedia(c *telegram.Client, chatID int64, media any, opts ...*telegram.MediaOptions) (*telegram.NewMessage, error) {
	mux := tgm.GetChatMutex(chatID)
	mux.Lock()
	defer mux.Unlock()

	message, err := c.SendMedia(chatID, media, opts...)
	time.Sleep(config.TimeBetweenMessages)

	return message, err
}

func (tgm *TGMutex) ForwardTo(c *telegram.Client, dstID, srcID int64, msgID int32) (*telegram.NewMessage, error) {
	mux := tgm.GetChatMutex(dstID)
	mux.Lock()
	defer mux.Unlock()

	messages, err := c.Forward(dstID, srcID, []int32{msgID})
	time.Sleep(config.TimeBetweenMessages)

	if err != nil {
		return nil, err
	}

	return &messages[0], nil
}
