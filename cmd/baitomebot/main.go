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
	"log"
	"net/http"

	"github.com/legzdev/BaitoMeBot/bot"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/db"
	"github.com/legzdev/BaitoMeBot/server"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db.Init()

	client, err := bot.New()
	if err != nil {
		log.Fatal(err)
	}

	handler := server.New(client)

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
