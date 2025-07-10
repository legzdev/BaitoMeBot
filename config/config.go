// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SourceURL = "https://github.com/legzdev/BaitoMeBot"

	ChunkSize = 1024 * 1024 // 1 MB
	AlignSize = 4 * 1024    // 4 KB
)

var (
	TelegramAppHash  string
	TelegramAppID    int32
	TelegramBotToken string
	TelegramChatID   int64

	AllowedUsers  []int64
	AuthChannelID int64
	HashLength    int
	ServerAddress string
	ServerHost    *url.URL

	TimeBetweenMessages time.Duration
	TimeBetweenChecks   time.Duration
)

type ErrEnvNotFound struct {
	Name string
}

func (e *ErrEnvNotFound) Error() string {
	return fmt.Sprintf("config: env %q not found", e.Name)
}

func Load() (err error) {
	TelegramAppHash = os.Getenv("TELEGRAM_API_HASH")
	if TelegramAppHash == "" {
		return &ErrEnvNotFound{Name: "TELEGRAM_API_HASH"}
	}

	TelegramAppID = getEnvNumber[int32]("TELEGRAM_API_ID", 0)
	if TelegramAppID == 0 {
		return &ErrEnvNotFound{Name: "TELEGRAM_API_ID"}
	}

	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if TelegramBotToken == "" {
		return &ErrEnvNotFound{Name: "TELEGRAM_BOT_TOKEN"}
	}

	TelegramChatID = getEnvNumber[int64]("TELEGRAM_CHAT_ID", 0)
	if TelegramChatID == 0 {
		return &ErrEnvNotFound{Name: "TELEGRAM_CHAT_ID"}
	}

	allowedUsers := os.Getenv("ALLOWED_USERS")
	if allowedUsers != "" {
		for idStr := range strings.SplitSeq(allowedUsers, ",") {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return err
			}

			AllowedUsers = append(AllowedUsers, id)
		}
	}

	ServerAddress = getEnvString("SERVER_ADDRESS", ":8080")
	serverHost := getEnvString("SERVER_HOST", DefaultServerHost())

	ServerHost, err = url.Parse(serverHost)
	if err != nil {
		return err
	}

	HashLength = getEnvNumber("HASH_LENGTH", 6)

	TimeBetweenMessages = getEnvDuration("TIME_BETWEEN_MESSAGES", 1*time.Second)
	TimeBetweenChecks = getEnvDuration("TIME_BETWEEN_CHECKS", 1*time.Minute)

	return nil
}

func DefaultServerHost() string {
	req, err := http.NewRequest(http.MethodGet, "http://ip.me", nil)
	if err != nil {
		return ""
	}

	// Get response in plain text
	req.Header.Set("User-Agent", "curl/1.18.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	ip := strings.TrimRight(string(body), "\n")
	address := "http://" + ip + ":8080"

	return address
}

func getEnvString(name string, value string) string {
	env := os.Getenv(name)
	if env != "" {
		return env
	}
	return value
}

type Number interface {
	int | int32 | int64
}

func getEnvNumber[T Number](name string, value T) T {
	env := os.Getenv(name)
	if env == "" {
		return value
	}

	number, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		log.Fatalf("env %s: invalid number %s", name, env)
	}

	return T(number)
}

func getEnvDuration(name string, value time.Duration) time.Duration {
	env := os.Getenv(name)
	if env == "" {
		return value
	}

	duration, err := time.ParseDuration(env)
	if err != nil {
		log.Fatalf("env %s: invalid duration %v", name, env)
	}

	return duration
}
