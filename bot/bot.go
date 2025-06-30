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
	"log"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/legzdev/BaitoMeBot/config"
)

var Bot = NewTGMutex()

func New() (*telegram.Client, error) {
	return NewWithToken(config.TelegramBotToken)
}

func NewWithToken(token string) (*telegram.Client, error) {
	clientConfig := telegram.ClientConfig{
		AppHash:      config.TelegramAppHash,
		AppID:        config.TelegramAppID,
		LogLevel:     telegram.LogWarn,
		FloodHandler: FloodHandler,
	}

	client, err := telegram.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	client.On("message:/start", AuthMiddleware(HelpHandler))
	client.On("message:/help", AuthMiddleware(HelpHandler))
	client.On("message:/txt", AuthMiddleware(OnTxtCommand))
	client.On("message", AuthMiddleware(OnMessage), telegram.FilterMedia)

	err = client.ConnectBot(token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Not sure if this works as expected
func FloodHandler(err error) bool {
	floodWait := telegram.GetFloodWait(err)

	duration := time.Duration(floodWait) * time.Second
	log.Printf("ratelimit reached (floodwait %v)", duration)

	time.Sleep(duration)
	return true
}
