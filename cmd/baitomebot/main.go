// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/bot"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/db"
	"github.com/legzdev/BaitoMeBot/server"
)

func main() {
	log.Println("Initializing...")
	ctx := context.Background()

	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}

	mainBot, err := bot.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = PrefetchTelegramChannel(mainBot.Client)
	if err != nil {
		log.Fatal(err)
	}

	handler := server.New(mainBot.Client)

	err = handler.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("=== Starting at %s ===", config.ServerAddress)

	err = http.ListenAndServe(config.ServerAddress, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func PrefetchTelegramChannel(client *telegram.Client) error {
	peer, err := db.Peers.GetChannel(config.TelegramChatID)
	if err != nil && peer != nil {
		return nil
	}

	input := &tg.InputChannel{
		ChannelID: config.TelegramChatID,
	}

	ctx := context.Background()
	ids := []tg.InputChannelClass{input}

	chats, err := client.API().ChannelsGetChannels(ctx, ids)
	if err != nil {
		return err
	}

	chatList := chats.GetChats()
	if len(chatList) == 0 {
		return errors.New("prefetch: cannot find channel peer")
	}

	chat := chatList[0]

	channel, ok := (chat).(*tg.Channel)
	if !ok {
		return errors.New("prefetch: invalid peer type")
	}

	peer = &tg.InputPeerChannel{
		ChannelID:  config.TelegramChatID,
		AccessHash: channel.AccessHash,
	}

	return db.Peers.Put(peer)
}
