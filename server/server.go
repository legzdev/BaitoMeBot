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
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/legzdev/BaitoMeBot/bot"
)

type Server struct {
	*http.ServeMux
	mux             sync.RWMutex
	startTime       time.Time
	extraTokens     []string
	Workers         []*telegram.Client
	WorkersCount    int
	NextWorkerIndex int
}

func New(client *telegram.Client) *Server {
	tokens := GetExtraBotTokens()

	// Default client + extra tokens
	workersCount := 1 + len(tokens)

	server := &Server{
		ServeMux:     http.NewServeMux(),
		startTime:    time.Now(),
		extraTokens:  make([]string, workersCount-1),
		Workers:      make([]*telegram.Client, workersCount),
		WorkersCount: workersCount,
	}

	server.Workers[0] = client

	server.HandleFunc("/{$}", server.ServerInfo)
	server.HandleFunc("/dl/{file_id}/{file_name}", server.Download)

	return server
}

func GetExtraBotTokens() []string {
	var tokens []string

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "BOT_TOKEN") {
			continue
		}

		envParts := strings.Split(env, "=")
		envValue := envParts[1]

		tokens = append(tokens, envValue)
	}

	return tokens
}

func (server *Server) Init() error {
	ctx := context.Background()

	for index, env := range server.extraTokens {
		envValue := strings.Split(env, "=")[1]

		worker, err := bot.NewWithToken(ctx, envValue, nil)
		if err != nil {
			return err
		}

		server.Workers[index+1] = worker
	}

	return nil
}

func (server *Server) GetWorker() *telegram.Client {
	server.mux.Lock()
	defer server.mux.Unlock()

	worker := server.Workers[server.NextWorkerIndex]

	server.NextWorkerIndex++
	if server.NextWorkerIndex == server.WorkersCount {
		server.NextWorkerIndex = 0
	}

	return worker
}
