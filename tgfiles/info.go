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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/db"
)

const (
	TimeLayout = "2006-01-02_15-04-05"
)

type FileInfo struct {
	ID       int
	Name     string
	MimeType string
	Size     int64
	Location tg.InputFileLocationClass
}

func GetFileInfo(message *tg.Message, userID int64) (*FileInfo, error) {
	info := &FileInfo{ID: message.ID}

	switch v := message.Media.(type) {
	case *tg.MessageMediaDocument:
		document, ok := v.Document.AsNotEmpty()
		if !ok {
			break
		}

		for _, attr := range document.Attributes {
			documentFileName, ok := attr.(*tg.DocumentAttributeFilename)
			if ok {
				info.Name = documentFileName.FileName
				break
			}
		}

		info.MimeType = document.MimeType
		info.Size = document.Size
		info.Location = document.AsInputDocumentFileLocation()

	case *tg.MessageMediaPhoto:
		photo, ok := v.Photo.AsNotEmpty()
		if !ok {
			break
		}

		info.Name = fmt.Sprintf("photo_%s.jpeg", time.Now().Format(TimeLayout))
		info.MimeType = "image/jpeg"

		photoSizes := len(photo.Sizes)
		if photoSizes == 0 {
			return nil, errors.New("photo without sizes")
		}

		photoSize, ok := photo.Sizes[photoSizes-1].AsNotEmpty()
		if !ok {
			return nil, errors.New("empty photo size")
		}

		size, ok := photoSize.(*tg.PhotoSize)
		if !ok {
			return nil, errors.New("photo with size 0")
		}

		info.Size = int64(size.GetSize())
		info.Location = &tg.InputPhotoFileLocation{
			ID:            photo.ID,
			AccessHash:    photo.AccessHash,
			FileReference: photo.FileReference,
			ThumbSize:     size.GetType(),
		}
	}

	state := db.GetState(userID)
	if state == db.StateTxtCaption || state == db.StateTxtCaptionFull {
		caption := ""

		if caption == "" {
			// do nothing
		} else if state == db.StateTxtCaption {
			captionLines := strings.Split(caption, "\n")
			info.Name = strings.TrimSpace(captionLines[0]) + GetFileExtension(info.Name)

		} else if state == db.StateTxtCaptionFull {
			caption = strings.ReplaceAll(caption, "\n", "")
			info.Name = caption + GetFileExtension(info.Name)
		}
	}

	if info.Name == "" {
		info.Name = fmt.Sprintf("file_%s.unknown", time.Now().Format(TimeLayout))
	}

	return info, nil
}

func GetFileExtension(fileName string) string {
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex != -1 {
		return fileName[lastDotIndex:]
	}

	return ".unknown"
}
