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

import "sync"

type State uint8

const (
	StateDefault State = iota
	StateTxt
	StateTxtCaption
	StateTxtCaptionFull
)

func (state State) String() string {
	switch state {
	case StateTxt:
		return "Default (filename)"
	case StateTxtCaption:
		return "Caption (oneline)"
	case StateTxtCaptionFull:
		return "Caption (full)"
	}

	return ""
}

var (
	states    map[UserID]State
	statesMux sync.Mutex
)

func SetState(userID int64, state State) {
	statesMux.Lock()
	defer statesMux.Unlock()

	states[userID] = state
}

func GetState(userID int64) State {
	statesMux.Lock()
	defer statesMux.Unlock()

	state, ok := states[userID]
	if ok {
		return state
	}

	return StateDefault
}

func DelState(userID int64) {
	statesMux.Lock()
	defer statesMux.Unlock()

	delete(states, userID)
}

func InTxtMode(userID int64) bool {
	state := GetState(userID)
	return state == StateTxt || state == StateTxtCaption
}
