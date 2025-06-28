// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"bytes"
	"sync"
)

var (
	buffers    map[UserID]*UserBuffer
	buffersMux sync.Mutex
)

type UserBuffer struct {
	*bytes.Buffer
	Name string
}

func NewUserBuffer(name string) *UserBuffer {
	return &UserBuffer{Name: name, Buffer: new(bytes.Buffer)}
}

func SetBuffer(userID int64, name string) {
	buffersMux.Lock()
	defer buffersMux.Unlock()

	buffers[userID] = NewUserBuffer(name)
}

func GetBuffer(userID int64) *UserBuffer {
	buffersMux.Lock()
	defer buffersMux.Unlock()

	return buffers[userID]
}

func DelBuffer(userID int64) {
	buffersMux.Lock()
	defer buffersMux.Unlock()

	delete(buffers, userID)
}
