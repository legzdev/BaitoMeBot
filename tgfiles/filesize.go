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

import "github.com/amarnathcjd/gogram/telegram"

func GetFileSize(f any) int64 {
	switch f := f.(type) {
	case *telegram.MessageMediaDocument:
		return f.Document.(*telegram.DocumentObj).Size

	case *telegram.MessageMediaPhoto:
		photo, ok := f.Photo.(*telegram.PhotoObj)
		if ok {
			if len(photo.Sizes) == 0 {
				return 0
			}
			s, _ := getPhotoSize(photo.Sizes[len(photo.Sizes)-1])
			return s
		}
	}

	return 0
}

func getPhotoSize(sizes telegram.PhotoSize) (int64, string) {
	switch s := sizes.(type) {
	case *telegram.PhotoSizeObj:
		return int64(s.Size), s.Type
	case *telegram.PhotoStrippedSize:
		if len(s.Bytes) < 3 || s.Bytes[0] != 1 {
			return int64(len(s.Bytes)), s.Type
		}
		return int64(len(s.Bytes)) + 622, s.Type
	case *telegram.PhotoCachedSize:
		return int64(len(s.Bytes)), s.Type
	case *telegram.PhotoSizeEmpty:
		return 0, s.Type
	case *telegram.PhotoSizeProgressive:
		return int64(getMax(s.Sizes)), s.Type
	default:
		return 0, "w"
	}
}

func getMax(a []int32) int32 {
	if len(a) == 0 {
		return 0
	}
	maximum := a[0]
	for _, v := range a {
		if v > maximum {
			maximum = v
		}
	}
	return maximum
}
