// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ServerInfo struct {
	OK     bool   `json:"ok"`
	Uptime string `json:"uptime"`
}

func (server *Server) ServerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	info := ServerInfo{
		OK:     true,
		Uptime: time.Since(server.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(&info)
	if err != nil {
		log.Println(err)
	}
}
