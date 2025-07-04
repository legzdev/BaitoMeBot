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
	"sync"

	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/errs"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"
	"golang.org/x/net/proxy"
)

type Bot struct {
	Client *telegram.Client
	mux    sync.Mutex
	muxes  map[string]*Mutex
}

func New(ctx context.Context) (bot *Bot, err error) {
	bot = &Bot{
		muxes: map[string]*Mutex{},
	}

	bot.NewTGMutex()

	dp := tg.NewUpdateDispatcher()
	dp.OnNewMessage(bot.MessageHandler())

	bot.Client, err = NewWithToken(ctx, config.TelegramBotToken, dp)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func NewWithToken(ctx context.Context, token string, uh telegram.UpdateHandler) (*telegram.Client, error) {
	opts := telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: "bot.session",
		},
		UpdateHandler: uh,
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: proxy.Dial,
		}),
	}

	client := telegram.NewClient(int(config.TelegramAppID), config.TelegramAppHash, opts)
	initChan := make(chan bool)

	go client.Run(ctx, func(ctx context.Context) error {
		initChan <- true
		for {
			select {
			case <-ctx.Done():
				return nil
			}
		}
	})

	<-initChan

	status, err := client.Auth().Status(ctx)
	if err != nil {
		return nil, errs.Wrap(err, "auth status")
	}

	if !status.Authorized {
		_, err = client.Auth().Bot(ctx, token)
		if err != nil {
			return nil, errs.Wrap(err, "auth login")
		}
	}

	return client, err
}
