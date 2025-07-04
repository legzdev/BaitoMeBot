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
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gotd/td/tg"
	"github.com/legzdev/BaitoMeBot/config"
	"github.com/legzdev/BaitoMeBot/db"
	"github.com/legzdev/BaitoMeBot/tgfiles"
)

func (server *Server) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	client := server.GetWorker()
	if client == nil {
		http.Error(w, "server bussy", http.StatusServiceUnavailable)
		return
	}

	fileID := r.PathValue("file_id")

	messageID, err := strconv.Atoi(fileID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	peer, err := db.Peers.GetChannel(config.TelegramChatID)
	if err != nil {
		log.Printf("failed to get channel from peers cache: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messagesReq := &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  peer.ChannelID,
			AccessHash: peer.AccessHash,
		},
		ID: []tg.InputMessageClass{
			&tg.InputMessageID{ID: messageID},
		},
	}

	messagesRes, err := client.API().ChannelsGetMessages(ctx, messagesReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	messages, ok := messagesRes.(*tg.MessagesChannelMessages)
	if !ok {
		log.Println("messages response are not channel messages")
		return
	}

	if len(messages.Messages) == 0 {
		http.NotFound(w, r)
		return
	}

	message, ok := messages.Messages[0].(*tg.Message)
	if !ok {
		http.NotFound(w, r)
		return
	}

	fileInfo, err := tgfiles.GetFileInfo(message, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := r.PathValue("file_name")
	if fileName == "" {
		fileName = fileInfo.Name
	}

	hash := tgfiles.GetShortHash(fileInfo)
	if hash != r.URL.Query().Get("hash") {
		http.NotFound(w, r)
		return
	}

	media := message.Media
	if media == nil {
		http.NotFound(w, r)
		return
	}

	rang, err := ParseRequestRange(r, fileInfo.Size)
	if err != nil {
		status := http.StatusRequestedRangeNotSatisfiable
		http.Error(w, http.StatusText(status), status)
		return
	}

	headers := w.Header()
	headers.Set("Content-Type", fileInfo.MimeType)

	contentLength := (rang.End - rang.Start) + 1
	headers.Set("Content-Length", fmt.Sprint(contentLength))

	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", fileName)
	headers.Set("Content-Disposition", contentDisposition)

	if r.Method == http.MethodHead {
		headers.Set("Accept-Ranges", "bytes")
		return
	}

	if rang.IsFromHeader {
		contentRange := fmt.Sprintf("bytes %d-%d/%d", rang.Start, rang.End, fileInfo.Size)
		headers.Set("Content-Range", contentRange)
		w.WriteHeader(http.StatusPartialContent)
	}

	reader, err := NewTelegramReader(ctx, client, message, rang.End)
	if err != nil {
		return
	}

	reader.Seek(rang.Start, io.SeekStart)

	chunkSize := GetChunkSize(contentLength)
	buffer := make([]byte, chunkSize)

	_, err = CopyBuffer(w, reader, buffer)
	if err != nil {
		log.Printf("ERRO: CopyBuffer: %v", err)
		return
	}
}

func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	var written int64

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, err := dst.Write(buf[:nr])
			if err != nil {
				return 0, err
			}

			written += int64(nw)
			continue
		}

		if errors.Is(err, io.EOF) {
			return written, nil
		}

		return 0, err
	}
}

func GetChunkSize(length int64) int {
	if length >= config.ChunkSize {
		return config.ChunkSize
	}

	if (length % config.AlignSize) == 0 {
		return int(length)
	}

	n := int(length / config.AlignSize)
	return (config.AlignSize * n) + config.AlignSize
}
