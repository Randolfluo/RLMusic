package handle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
)

type SongAuth struct{}

// ScanUserMusic 扫描用户音乐文件并录入数据库
func (*SongAuth) ScanUserMusic(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)
	basePath := g.GetConfig().BasicPath
	// 获取用户根目录: basicPath/username
	userPath := filepath.Join(basePath.FilePath, basePath.FileName, user.Username)

	// 定义要扫描的子文件夹及其对应的公开权限
	// public: 公开
	// group: 暂时设为非公开 (或对应特定业务逻辑)
	// private: 私有
	folders := map[string]string{
		"public":  "public",
		"group":   "group",
		"private": "private",
	}

	// 支持的音频扩展名
	supportedExts := map[string]bool{
		".mp3": true, ".flac": true, ".wav": true, ".ogg": true, ".m4a": true,
	}

	addedCount := 0
	updatedCount := 0

	for subDir, permission := range folders {
		targetDir := filepath.Join(userPath, subDir)

		// 检查目录是否存在，不存在则跳过
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				slog.Error("Error accessing path", "path", path, "error", err)
				return nil
			}
			if info.IsDir() {
				return nil
			}

			// 检查后缀
			ext := strings.ToLower(filepath.Ext(path))
			if !supportedExts[ext] {
				return nil
			}

			// 读取音频元数据
			f, err := os.Open(path)
			if err != nil {
				slog.Error("Failed to open file", "path", path, "error", err)
				return nil
			}
			defer f.Close()

			var songTitle, songArtist, songAlbum string
			var coverData []byte
			var coverMime string
			var year string
			var trackNum, discNum int

			m, err := tag.ReadFrom(f)
			if err == nil {
				songTitle = m.Title()
				songArtist = m.Artist()
				songAlbum = m.Album()
				if m.Year() > 0 {
					year = strconv.Itoa(m.Year())
				}
				trackNum, _ = m.Track()
				discNum, _ = m.Disc()

				if pic := m.Picture(); pic != nil {
					coverData = pic.Data
					coverMime = pic.MIMEType
				}
			} else {
				slog.Warn("Failed to read metadata tags, using filename", "path", path, "error", err)
			}

			// 默认值
			if songTitle == "" {
				songTitle = strings.TrimSuffix(info.Name(), ext)
			}
			if songArtist == "" {
				songArtist = "Unknown Artist"
			}
			if songAlbum == "" {
				songAlbum = "Unknown Album"
			}

			// 1. 处理 Artist
			artist, err := model.FindOrCreateArtist(db, songArtist)
			if err != nil {
				slog.Error("Failed to create/find Artist", "name", songArtist, "error", err)
				return nil
			}

			// 2. 处理 Cover
			var coverID *int
			if len(coverData) > 0 {
				// 计算 hash 去重
				hash := sha256.Sum256(coverData)
				checksum := hex.EncodeToString(hash[:])

				if cid, err := model.FindOrCreateCover(db, coverData, coverMime, checksum); err == nil {
					coverID = cid
				}
			}

			// 3. 处理 Album
			album, err := model.FindOrCreateAlbum(db, songAlbum, artist.ID, coverID)
			if err != nil {
				slog.Error("Failed to create/find Album", "title", songAlbum, "error", err)
				return nil
			}

			// 4. 处理 Song
			// 查找是否存在（根据路径）
			song, err := model.FindSongByPath(db, path)
			if err != nil {
				song = &model.Song{} // 新建
			}

			// 准备数据
			song.Title = songTitle
			song.ArtistID = &artist.ID
			song.AlbumID = &album.ID
			song.CoverImageID = coverID
			song.TrackNum = trackNum
			song.DiscNum = discNum
			song.Year = year
			// 文件信息
			song.FilePath = path
			song.FileName = info.Name()
			song.FileSize = info.Size()
			song.Format = strings.TrimPrefix(ext, ".")
			// 权限信息 (如果是新歌，或者覆盖更新)
			song.OwnerID = &user.ID
			song.Permission = permission

			// 保存
			isCreated, err := model.SaveSong(db, song)
			if err == nil {
				if isCreated {
					addedCount++
				} else {
					updatedCount++
				}
			} else {
				slog.Error("Failed to save song", "title", songTitle, "error", err)
			}

			// 自动添加到对应权限的歌单
			playlistName := strings.Title(permission) // Public, Private, Group
			playlist, err := model.FindOrCreatePlaylist(db, user.ID, playlistName, permission)
			if err != nil {
				slog.Error("Failed to create/find Playlist", "name", playlistName, "error", err)
			} else {
				// 关联歌曲到歌单 (如果尚未关联)
				if err := model.AddSongToPlaylist(db, playlist, song); err != nil {
					slog.Error("Failed to add song to playlist", "playlist", playlistName, "song", song.Title, "error", err)
				}
			}

			return nil
		})

		if err != nil {
			slog.Error("Walk folder failed", "path", targetDir, "error", err)
		}
	}

	ReturnSuccess(c, gin.H{
		"added":   addedCount,
		"updated": updatedCount,
		"message": "扫描完成",
	})
}

// StreamSong 流式传输歌曲
func (*SongAuth) StreamSong(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	song, err := model.GetSongByID(db, idStr)
	if err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	// 权限检查
	currentUser := GetCurrentUser(c)
	// 如果不是公开的，且当前用户不是所有者
	if song.Permission != "public" && (currentUser == nil || *song.OwnerID != currentUser.ID) {
		ReturnError(c, g.ErrPermission, "无权访问此歌曲")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
		ReturnError(c, g.ErrFileNotExist, "音频文件丢失")
		return
	}

	// 使用 Gin 的 File 响应，它自动处理 Range 头实现流式传输
	c.File(song.FilePath)
}

// GetSongCover 获取封面图片
func (*SongAuth) GetSongCover(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	song, err := model.GetSongWithCover(db, idStr)
	if err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	// 简单的权限检查 (封面通常比较宽容，或者跟随歌曲权限)
	// if !song.IsPublic && ...

	if song.CoverImageID == nil || len(song.CoverImage.Data) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, song.CoverImage.MimeType, song.CoverImage.Data)
}

// GetPlayList 获取播放列表 (示例：获取所有歌曲)
func (*SongAuth) GetPlayList(c *gin.Context) {
	db := GetDB(c)
	user := GetCurrentUser(c)

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	songs, total, err := model.GetSongsList(db, user.ID, page, pageSize)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"list":     songs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetMetaCover 获取封面图片 (直接通过封面ID)
func (*SongAuth) GetMetaCover(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	db := GetDB(c)

	cover, err := model.GetCover(db, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if len(cover.Data) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// 设置缓存控制
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存1年，因为 checksum 应该不变
	c.Header("ETag", fmt.Sprintf(`"%s"`, cover.Checksum))

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, cover.Checksum) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	c.Data(http.StatusOK, cover.MimeType, cover.Data)
}

// GetPlaylists 获取所有歌单 (包括自己的和公开的)
func (*SongAuth) GetPlaylists(c *gin.Context) {
	db := GetDB(c)
	user := GetCurrentUser(c)

	playlists, err := model.GetUserPlaylists(db, user.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
}
