package handle

import (
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/audio"
	"server/internal/utils/imgtool"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// 支持的音频扩展名
	supportedExts := map[string]bool{
		".flac": true, //".wav": true, ".ogg": true, ".m4a": true,
	}

	addedCount := 0
	updatedCount := 0

	// 读取用户目录下的一级子目录
	entries, err := os.ReadDir(userPath)
	if err != nil {
		// 如果目录不存在，可能只是用户还没传文件，不报错
		if os.IsNotExist(err) {
			ReturnSuccess(c, gin.H{"message": "目录为空"})
			return
		}
		slog.Error("Failed to read user directory", "path", userPath, "error", err)
		ReturnError(c, g.Err, "读取目录失败")
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // 忽略根目录下的文件
		}

		subDirName := entry.Name()
		targetDir := filepath.Join(userPath, subDirName)

		// 确定权限: "public" 文件夹为公开，其他默认为私有
		permission := "private"
		if subDirName == "public" {
			permission = "public"
		}

		// 查找或创建歌单 (歌单名 = 文件夹名)
		currentPlaylist, err := model.FindOrCreatePlaylist(db, user.ID, subDirName, permission)
		if err != nil {
			slog.Error("Failed to create/find Playlist", "name", subDirName, "error", err)
			continue // 歌单创建失败，跳过该文件夹的扫描? 或者继续扫描但不加歌单? 这里选择跳过
		}

		err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
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

			// ... (rest of the logic)

			// 读取音频元数据
			f, err := os.Open(path)
			if err != nil {
				slog.Error("Failed to open file", "path", path, "error", err)
				return nil
			}
			defer f.Close()

			var songTitle, songArtist, songAlbum string
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

			} else {
				slog.Warn("Failed to read metadata tags, using filename", "path", path, "error", err)
			}

			// 默认值
			if songTitle == "" {
				songTitle = strings.TrimSuffix(info.Name(), ext)
			}
			// 如果数据为空，对应数据库值也为空 (不设置默认 Unknown)

			// 1. 处理 Artist
			var artistID *int
			if songArtist != "" {
				artist, err := model.FindOrCreateArtist(db, songArtist)
				if err != nil {
					slog.Error("Failed to create/find Artist", "name", songArtist, "error", err)
					return nil
				}
				artistID = &artist.ID
			}

			// 2. 处理 Cover
			var coverID *int
			var coverUrl string // 用于歌单封面
			if m != nil && m.Picture() != nil {
				pic := m.Picture()
				// 保存目录: ./data/covers (相对运行目录)
				hash, filename, width, height, err := imgtool.ProcessAndSaveCover(pic.Data, "./data/covers")
				if err == nil {
					// 数据库记录
					cover, err := model.FindOrCreateCover(db, hash, filename, pic.Ext, int64(len(pic.Data)), width, height)
					if err == nil {
						coverID = &cover.ID
						coverUrl = "/covers/" + filename
					} else {
						slog.Error("Failed to find/create Cover", "error", err)
					}
				} else {
					slog.Warn("Failed to process cover image", "path", path, "error", err)
				}
			}

			// 3. 处理 Album
			var albumID *int
			if songAlbum != "" {
				album, err := model.FindOrCreateAlbum(db, songAlbum, artistID)
				if err != nil {
					slog.Error("Failed to create/find Album", "title", songAlbum, "error", err)
					return nil
				}
				albumID = &album.ID
			}

			// 4. 处理 Song
			// 查找是否存在（根据路径）
			song, err := model.FindSongByPath(db, path)
			if err != nil {
				song = &model.Song{} // 新建
			}

			// 准备数据
			song.Title = songTitle
			song.ArtistName = songArtist
			song.AlbumName = songAlbum
			song.ArtistID = artistID
			song.AlbumID = albumID
			song.CoverID = coverID

			song.TrackNum = trackNum
			song.DiscNum = discNum
			song.Year = year
			// 文件信息
			song.FilePath = path
			song.FileName = info.Name()
			song.FileSize = info.Size()
			song.Format = strings.TrimPrefix(ext, ".")

			// 读取音频参数 (目前仅支持 FLAC)
			if ext == ".flac" {
				if props, err := audio.ParseFlacProps(path); err == nil {
					song.Duration = props.Duration
					song.SampleRate = props.SampleRate
					song.BitDepth = props.BitDepth
					song.Channels = props.Channels
					song.BitRate = props.BitRate
				} else {
					slog.Warn("Failed to parse FLAC props", "path", path, "error", err)
				}
			}

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

			// 关联歌曲到当前目录对应的歌单
			if err := model.AddSongToPlaylist(db, currentPlaylist, song); err != nil {
				slog.Error("Failed to add song to playlist", "playlist", currentPlaylist.Title, "song", song.Title, "error", err)
			}

			// 如果歌单封面为空，且当前歌曲有封面，设置该封面为歌单封面
			if currentPlaylist.CoverUrl == "" && coverUrl != "" {
				if err := db.Model(currentPlaylist).Update("cover_url", coverUrl).Error; err != nil {
					slog.Warn("Failed to update playlist cover", "playlist", currentPlaylist.Title, "error", err)
				}
				currentPlaylist.CoverUrl = coverUrl
			}

			return nil
		})

		if err != nil {
			slog.Error("Walk folder failed", "path", targetDir, "error", err)
		}
	}
	slog.Info("Music scan completed", "user", user.Username, "added", addedCount, "updated", updatedCount)
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

	// 检查文件是否存在
	if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
		ReturnError(c, g.ErrFileNotExist, "音频文件丢失")
		return
	}

	// 使用 Gin 的 File 响应，它自动处理 Range 头实现流式传输
	c.File(song.FilePath)
}

// GetSongCover 获取封面图片 (已移除)
// func (*SongAuth) GetSongCover(c *gin.Context) { ... }

// GetPublicPlaylists 获取公共歌单
func (*SongAuth) GetPublicPlaylists(c *gin.Context) {
	db := GetDB(c)

	playlists, err := model.GetPublicPlaylists(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
}

// GetPrivatePlaylists 获取私人歌单
func (*SongAuth) GetPrivatePlaylists(c *gin.Context) {
	db := GetDB(c)
	user := GetCurrentUser(c)

	playlists, err := model.GetPrivatePlaylists(db, user.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
}

// GetPlaylistDetail 获取歌单详情(含歌曲)
func (*SongAuth) GetPlaylistDetail(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	playlistDetail, err := model.GetPlaylistDetail(db, idStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ReturnError(c, g.ErrDbOp, "歌单不存在")
		} else {
			ReturnError(c, g.ErrDbOp, err)
		}
		return
	}

	ReturnSuccess(c, playlistDetail)
}
