// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package tgfiles

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/legzdev/BaitoMeBot/config"
)

func GetShortHash(info FileInfo) string {
	return GetHash(info)[:config.HashLength]
}

func GetHash(info FileInfo) string {
	hasher := md5.New()

	io.WriteString(hasher, fmt.Sprint(info.ID))
	io.WriteString(hasher, fmt.Sprint(info.Size))
	io.WriteString(hasher, info.MimeType)

	return hex.EncodeToString(hasher.Sum(nil))
}
