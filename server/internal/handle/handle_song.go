package handle

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	g "server/internal/global"
	"server/internal/model"
	"server/internal/utils/ai"
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

	if user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "权限不足，仅管理员可扫描")
		return
	}

	db := GetDB(c)
	basePath := g.GetConfig().BasicPath
	// 获取基础目录: basicPath
	// 修改扫描逻辑：不再扫描 username 目录，而是扫描 basicPath 下的所有文件夹（排除 data 文件夹）
	rootPath := filepath.Join(basePath.FilePath, basePath.FileName)

	// 支持的音频扩展名
	supportedExts := map[string]bool{
		".flac": true, ".mp3": true, ".wav": true, ".ogg": true, ".m4a": true,
	}

	addedCount := 0
	updatedCount := 0
	var scannedDuration float64 = 0

	// 读取根目录下的一级子目录
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			ReturnSuccess(c, gin.H{"message": "目录为空"})
			return
		}
		slog.Error("Failed to read root directory", "path", rootPath, "error", err)
		ReturnError(c, g.Err, "读取目录失败")
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // 忽略根目录下的文件
		}

		subDirName := entry.Name()
		// 跳过 data, avatar, podcast 文件夹
		if subDirName == "data" || subDirName == "avatar" || subDirName == "podcast" {
			continue
		}

		targetDir := filepath.Join(rootPath, subDirName)

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

			// 获取歌曲所在目录名作为歌单名
			dirName := filepath.Base(filepath.Dir(path))

			// 查找或创建歌单 (歌单名 = 歌曲所在目录名)
			currentPlaylist, err := model.FindOrCreatePlaylist(db, user.ID, dirName)
			if err != nil {
				slog.Error("Failed to create/find Playlist", "name", dirName, "error", err)
				return nil // 跳过该文件
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
			var songArtists []model.Artist
			if songArtist != "" {
				// 替换常见分隔符为 /
				tempInfo := songArtist
				tempInfo = strings.ReplaceAll(tempInfo, ";", "/")
				tempInfo = strings.ReplaceAll(tempInfo, "；", "/") // 中文分号
				tempInfo = strings.ReplaceAll(tempInfo, "、", "/") // 中文顿号

				names := strings.Split(tempInfo, "/")
				for _, name := range names {
					name = strings.TrimSpace(name)
					if name == "" {
						continue
					}

					artist, err := model.FindOrCreateArtist(db, name)
					if err != nil {
						slog.Error("Failed to create/find Artist", "name", name, "error", err)
						continue
					}
					songArtists = append(songArtists, *artist)

					// 使用第一个扫描到的作为主要关联ID (兼容旧逻辑)
					if artistID == nil {
						artistID = &artist.ID
					}
				}
			}

			// 2. 处理 Cover
			var coverID *int
			var coverUrl string // 用于歌单封面
			if m != nil && m.Picture() != nil {
				pic := m.Picture()
				// 保存目录: config/data/cover
				conf := g.GetConfig()
				saveDir := filepath.Join(conf.BasicPath.FilePath, conf.BasicPath.FileName, "data", "cover")
				hash, filename, width, height, err := imgtool.ProcessAndSaveCover(pic.Data, saveDir)
				if err == nil {
					// 数据库记录
					cover, err := model.FindOrCreateCover(db, hash, filename, pic.Ext, int64(len(pic.Data)), width, height)
					if err == nil {
						coverID = &cover.ID
						coverUrl = filename
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
			song.Artists = songArtists
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
			// 读取音频参数
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
			} else if ext == ".mp3" {
				if props, err := audio.ParseMp3Props(path); err == nil {
					song.Duration = props.Duration
					song.BitRate = props.BitRate
				} else {
					slog.Warn("Failed to parse MP3 props", "path", path, "error", err)
				}
			} else if ext == ".wav" {
				if props, err := audio.ParseWavProps(path); err == nil {
					song.Duration = props.Duration
					song.BitRate = props.BitRate
				} else {
					slog.Warn("Failed to parse WAV props", "path", path, "error", err)
				}
			} else {
				// TODO: 其他格式 (OGG, M4A) 暂未实现原生解析
				// 可以考虑引入 ffmpeg wrapper
				slog.Info("Duration parsing not implemented for this format yet", "format", ext)
			}

			// 保存
			isCreated, err := model.SaveSong(db, song)
			if err == nil {
				// 显式更新 Artists 关联 (针对更新场景)
				if !isCreated && len(song.Artists) > 0 {
					_ = db.Model(song).Association("Artists").Replace(song.Artists)
				}

				if isCreated {
					addedCount++
				} else {
					updatedCount++
				}

				// 尝试为关联的艺术家设置封面
				if coverUrl != "" {
					for _, artist := range songArtists {
						db.Model(&model.Artist{}).Where("id = ? AND (cover IS NULL OR cover = '')", artist.ID).Update("cover", coverUrl)
					}
					// 尝试为关联的专辑设置封面
					if song.AlbumID != nil {
						db.Model(&model.Album{}).Where("id = ? AND (cover IS NULL OR cover = '')", *song.AlbumID).Update("cover", coverUrl)
					}
				}

				// 累加时长
				scannedDuration += song.Duration
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

	// 更新用户总时长 (秒)
	user.TotalDuration = int64(scannedDuration)
	db.Save(user)

	// 更新系统统计
	var songCount, albumCount, artistCount int64
	db.Model(&model.Song{}).Count(&songCount)
	db.Model(&model.Album{}).Count(&albumCount)
	db.Model(&model.Artist{}).Count(&artistCount)

	// 计算所有用户总时长
	var systemTotalDuration int64
	db.Model(&model.User{}).Select("sum(total_duration)").Scan(&systemTotalDuration)

	// 计算所有歌曲总大小
	var systemTotalVolume int64
	db.Model(&model.Song{}).Select("sum(file_size)").Scan(&systemTotalVolume)

	_ = model.UpdateSystemInfoStats(db, songCount, albumCount, artistCount, systemTotalDuration, systemTotalVolume)

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

	// 增加播放次数 (异步执行，不阻塞播放)
	go func() {
		//重新获取一个新的DB实例(最好是新的会话)，虽然GORM DB是并发安全的，但为了避免上下文取消等问题
		// 这里直接用 db 即可，因为它是个 *gorm.DB
		db.Model(song).UpdateColumn("play_count", gorm.Expr("play_count + ?", 1))
	}()

	// 使用 Gin 的 File 响应，它自动处理 Range 头实现流式传输
	c.File(song.FilePath)
}

// GetSongOpeningAudio 获取歌曲开场白音频
func (*SongAuth) GetSongOpeningAudio(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	var song model.Song
	if err := db.First(&song, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	if song.OpeningAudioFile == "" {
		ReturnError(c, g.ErrFileNotExist, "开场白音频不存在")
		return
	}

	// 拼接音频文件路径: BasicPath/data/Podcast/filename
	// 注意：这里需要确保 song.OpeningAudioFile 只是文件名
	conf := g.GetConfig()
	audioPath := filepath.Join(conf.BasicPath.FilePath, conf.BasicPath.FileName, "data", "Podcast", song.OpeningAudioFile)

	// 检查文件是否存在
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		ReturnError(c, g.ErrFileNotExist, "音频文件丢失")
		return
	}

	c.File(audioPath)
}

// GetSongOpeningText 获取歌曲开场白文本
func (*SongAuth) GetSongOpeningText(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	var song model.Song
	if err := db.First(&song, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	ReturnSuccess(c, gin.H{
		"description": song.Description,
	})
}

// GetSongCover 获取歌曲封面
func (*SongAuth) GetSongCover(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	var song model.Song
	// Preload Cover specifically
	if err := db.Preload("Cover").First(&song, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	// 1. 优先尝试从音频源文件中提取完整的高清封面
	if song.FilePath != "" {
		f, err := os.Open(song.FilePath)
		if err == nil {
			defer f.Close()
			// 读取标签信息
			m, err := tag.ReadFrom(f)
			if err == nil && m != nil {
				pic := m.Picture()
				if pic != nil {
					// 直接返回原始图片数据
					c.Data(200, pic.MIMEType, pic.Data)
					return
				}
			} else {
				// 读取失败日志，仅调试用
				// slog.Warn("Failed to read tags from file", "path", song.FilePath, "error", err)
			}
		}
	}

	// 2. 如果源文件没有封面或读取失败，降级使用缓存的缩略图
	if song.CoverID == nil || song.Cover.Path == "" {
		// 尝试返回默认封面
		// 查找可能的路径 (兼容不同的运行目录)
		candidates := []string{
			"public/images/logo/favicon.png",
			"../public/images/logo/favicon.png",
			"../../public/images/logo/favicon.png",
			filepath.Join(g.GetConfig().BasicPath.FilePath, "localmusicplayer", "public", "images", "logo", "favicon.png"),
		}

		for _, path := range candidates {
			// 转换为绝对路径（如果是相对路径）
			absPath, _ := filepath.Abs(path)
			if _, err := os.Stat(absPath); err == nil {
				c.File(absPath)
				return
			}
		}

		// 404 Not Found if no cover and no default
		c.Status(404)
		return
	}

	// 假设封面存在 ./data/covers (需与 imgtool/scan 逻辑保持一致)
	coverPath := filepath.Join("./data/covers", song.Cover.Path)

	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		c.Status(404)
		return
	}

	c.File(coverPath)
}

// GetSongDetail 获取歌曲详细信息
func (*SongAuth) GetSongDetail(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)

	var song model.Song
	// Preload all associations
	// Artist: Preload singular if needed, but Artists (plural) is many2many
	if err := db.Preload("Artist").Preload("Artists").Preload("Album").Preload("Cover").First(&song, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	ReturnSuccess(c, song)
}

// GetArtistDetail 获取歌手详情
func (*SongAuth) GetArtistDetail(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)
	var artist model.Artist
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 30
	}

	if err := db.First(&artist, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌手不存在")
		return
	}

	if artist.Cover != "" && !strings.HasPrefix(artist.Cover, "/covers/") && !strings.HasPrefix(artist.Cover, "http") {
		artist.Cover = "/covers/" + artist.Cover
	}

	var total int64
	if err := db.Model(&model.Song{}).
		Joins("JOIN song_artists ON song_artists.song_id = song.id").
		Where("song_artists.artist_id = ?", artist.ID).
		Count(&total).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	var songsRaw []model.Song
	if err := db.Joins("JOIN song_artists ON song_artists.song_id = song.id").
		Where("song_artists.artist_id = ?", artist.ID).
		Order("song.id DESC").
		Limit(limit).Offset(offset).
		Preload("Cover").
		Find(&songsRaw).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	var songs []model.SimpleSongResponse
	for _, s := range songsRaw {
		artistId := 0
		albumId := 0
		if s.ArtistID != nil {
			artistId = *s.ArtistID
		}
		if s.AlbumID != nil {
			albumId = *s.AlbumID
		}

		coverUrl := ""
		if s.CoverID != nil && s.Cover.ID != 0 {
			coverUrl = "/covers/" + s.Cover.Path
		}

		songs = append(songs, model.SimpleSongResponse{
			ID:         s.ID,
			Title:      s.Title,
			ArtistName: s.ArtistName,
			AlbumTitle: s.AlbumName,
			AlbumName:  s.AlbumName,
			Duration:   s.Duration,
			Year:       s.Year,
			ArtistID:   artistId,
			AlbumID:    albumId,
			CoverUrl:   coverUrl,
			HasIntro:   s.OpeningAudioFile != "",
		})
	}

	ReturnSuccess(c, gin.H{
		"artist": artist,
		"list":   songs,
		"total":  total,
	})
}

// GetAlbumDetail 获取专辑详情
func (*SongAuth) GetAlbumDetail(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)
	var album model.Album
	if err := db.Preload("Artist").Preload("Songs").First(&album, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "专辑不存在")
		return
	}

	if album.Cover != "" && !strings.HasPrefix(album.Cover, "/covers/") && !strings.HasPrefix(album.Cover, "http") {
		album.Cover = "/covers/" + album.Cover
	}

	ReturnSuccess(c, album)
}

// GetAllPlaylists 获取所有公开歌单
func (*SongAuth) GetAllPlaylists(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)

	playlists, total, err := model.GetPublicPlaylists(db, page, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"list":  playlists,
		"total": total,
	})
}

// GetUserPrivatePlaylists 获取用户私有歌单
func (*SongAuth) GetUserPrivatePlaylists(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)

	playlists, total, err := model.GetUserPrivatePlaylists(db, user.ID, page, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"list":  playlists,
		"total": total,
	})
}

type UpdatePlaylistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"` // Using pointer to distinguish between false and nil
}

// UpdatePlaylist 更新歌单信息(仅限Owner)
func (*SongAuth) UpdatePlaylist(c *gin.Context) {
	idStr := c.Param("id")
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	db := GetDB(c)
	var playlist model.Playlist
	if err := db.First(&playlist, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	// 权限检查: 仅 Owner 可修改
	if playlist.OwnerID != user.ID {
		ReturnError(c, g.ErrPermission, "无权修改此歌单")
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) > 0 {
		if err := db.Model(&playlist).Updates(updates).Error; err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
	}

	ReturnSuccess(c, "更新成功")
}

// ToggleLike 切换歌曲的"喜欢"状态
func (*SongAuth) ToggleLike(c *gin.Context) {
	songIDStr := c.Param("id")
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}
	db := GetDB(c)

	// 查找用户的"我喜欢的音乐"歌单
	var playlist model.Playlist
	err := db.Where("owner_id = ? AND title = ?", user.ID, "我喜欢的音乐").First(&playlist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，自动创建(私有)
			pl, err := model.CreatePlaylist(db, user.ID, "我喜欢的音乐", "因为热爱，所以收藏", false)
			if err != nil {
				ReturnError(c, g.ErrDbOp, "创建喜欢歌单失败")
				return
			}
			playlist = *pl
		} else {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
	}

	// 获取歌曲
	song, err := model.GetSongByID(db, songIDStr)
	if err != nil {
		ReturnError(c, g.ErrDbOp, "歌曲不存在")
		return
	}

	// 检查是否已在歌单中
	if model.IsSongInPlaylist(db, playlist.ID, song.ID) {
		// 存在则移除
		if err := model.RemoveSongFromPlaylist(db, &playlist, song); err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
		ReturnSuccess(c, gin.H{"liked": false, "message": "已取消喜欢"})
	} else {
		// 不存在则添加
		if err := model.AddSongToPlaylist(db, &playlist, song); err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
		ReturnSuccess(c, gin.H{"liked": true, "message": "已添加到喜欢"})
	}
}

// GetLikedSongs 获取用户喜欢的歌曲列表
func (*SongAuth) GetLikedSongs(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)

	// 1. Find "我喜欢的音乐" playlist
	var playlist model.Playlist
	err := db.Where("owner_id = ? AND title = ?", user.ID, "我喜欢的音乐").First(&playlist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No liked songs yet
			ReturnSuccess(c, gin.H{
				"list":  []model.Song{},
				"total": 0,
			})
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 2. Get songs in playlist
	playlistDetail, err := model.GetPlaylistDetail(db, strconv.Itoa(playlist.ID), page, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"list":  playlistDetail.Songs,
		"total": playlistDetail.Total,
	})
}

// GetPublicPlaylistDetail 获取公共歌单详情
func (*SongAuth) GetPublicPlaylistDetail(c *gin.Context) {
	idStr := c.Param("id")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)

	playlistDetail, err := model.GetPlaylistDetail(db, idStr, page, limit)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ReturnError(c, g.ErrDbOp, "歌单不存在")
		} else {
			ReturnError(c, g.ErrDbOp, err)
		}
		return
	}

	if !playlistDetail.IsPublic {
		ReturnError(c, g.ErrPermission, "该歌单为私有歌单，无法访问")
		return
	}

	ReturnSuccess(c, playlistDetail)
}

// GetPrivatePlaylistDetail 获取私有歌单详情(需要验证 Owner)
func (*SongAuth) GetPrivatePlaylistDetail(c *gin.Context) {
	idStr := c.Param("id")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)

	playlistDetail, err := model.GetPlaylistDetail(db, idStr, page, limit)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ReturnError(c, g.ErrDbOp, "歌单不存在")
		} else {
			ReturnError(c, g.ErrDbOp, err)
		}
		return
	}

	// 验证 Owner
	if playlistDetail.OwnerID != user.ID {
		ReturnError(c, g.ErrPermission, "无权访问此私有歌单")
		return
	}

	ReturnSuccess(c, playlistDetail)
}

// GetPlaylistDetail 获取歌单详情（随机返回歌曲，上限100首）
func (*SongAuth) GetPlaylistDetail(c *gin.Context) {
	idStr := c.Param("id")
	limitStr := c.Query("limit")

	limit := 100
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > 100 {
			limit = 100
		} else {
			limit = l
		}
	}

	db := GetDB(c)

	playlistDetail, err := model.GetPlaylistRandomSongs(db, idStr, limit)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ReturnError(c, g.ErrDbOp, "歌单不存在")
		} else {
			ReturnError(c, g.ErrDbOp, err)
		}
		return
	}

	user := GetCurrentUser(c)
	// 权限检查
	// 规则:
	// 1. 当前用户为管理员时，允许访问公共歌单 (隐含: 也能访问自己的私有歌单)
	// 2. 当前用户不为管理员时，只允许访问自己的私有歌单
	allowed := false
	isAdmin := user != nil && user.UserGroup == "admin"
	isOwner := user != nil && playlistDetail.OwnerID == user.ID

	if isAdmin {
		if playlistDetail.IsPublic || isOwner {
			allowed = true
		}
	} else {
		if isOwner && !playlistDetail.IsPublic {
			allowed = true
		}
	}

	if !allowed {
		ReturnError(c, g.ErrPermission, "无权访问此歌单(权限受限)")
		return
	}

	// 构造返回数据
	type SongInfo struct {
		Title      string `json:"title"`
		Year       string `json:"year"`
		ArtistName string `json:"artist"`
		AlbumName  string `json:"album"`
	}

	var songs []SongInfo
	for _, song := range playlistDetail.Songs {
		songs = append(songs, SongInfo{
			Title:      song.Title,
			Year:       song.Year,
			ArtistName: song.ArtistName,
			AlbumName:  song.AlbumName,
		})
	}

	ReturnSuccess(c, gin.H{
		"playlist": playlistDetail.Title,
		"songs":    songs,
	})
}

// GetPlaylistAIAnalysis 获取歌单AI分析
func (*SongAuth) GetPlaylistAIAnalysis(c *gin.Context) {
	idStr := c.Param("id")

	db := GetDB(c)

	// 使用 GetPlaylistRandomSongs 获取歌单随机歌曲，限制100首，用于分析
	playlistDetail, err := model.GetPlaylistRandomSongs(db, idStr, 100)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ReturnError(c, g.ErrDbOp, "歌单不存在")
		} else {
			ReturnError(c, g.ErrDbOp, err)
		}
		return
	}

	user := GetCurrentUser(c)
	// 权限检查
	// 管理员可为所有歌单生成描述，普通用户只能为自己的私有歌单生成描述。
	allowed := false
	isAdmin := user != nil && user.UserGroup == "admin"
	isOwner := user != nil && playlistDetail.OwnerID == user.ID

	if isAdmin {
		allowed = true // Admin can access all
	} else {
		// Regular user can only access their own private playlists (and public if they own it? User request says "own private playlists")
		// "普通用户只能为自己的私有歌单生成描述"
		if isOwner && !playlistDetail.IsPublic {
			allowed = true
		}
	}

	if !allowed {
		ReturnError(c, g.ErrPermission, "无权访问此歌单(权限受限)")
		return
	}

	// 读取 prompt1.md
	promptBytes, err := os.ReadFile("prompts/prompt_playlist1.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_playlist1.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板")
		return
	}
	promptTemplate := string(promptBytes)

	// 构造 AI 分析所需的简化数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range playlistDetail.Songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 替换 prompt 中的占位符
	finalPrompt := strings.Replace(promptTemplate, "{{json数据}}", string(jsonData), 1)

	// 调用 AI 接口
	reply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 尝试解析返回的 JSON (如果是 JSON 格式的话)
	var result interface{}
	// 清理可能的 markdown 代码块标记 ```json ... ```
	cleanReply := strings.TrimSpace(reply)
	cleanReply = strings.TrimPrefix(cleanReply, "```json")
	cleanReply = strings.TrimPrefix(cleanReply, "```")
	cleanReply = strings.TrimSuffix(cleanReply, "```")

	if err := json.Unmarshal([]byte(cleanReply), &result); err == nil {
		ReturnSuccess(c, result)
	} else {
		// 如果不是 JSON，直接返回字符串
		ReturnSuccess(c, gin.H{
			"analysis": reply,
		})
	}
}

// internalGetAnalysis 内部方法：获取歌单分析结果
func (h *SongAuth) internalGetAnalysis(c *gin.Context, idStr string) (string, *g.Result) {
	db := GetDB(c)

	// 使用 GetPlaylistRandomSongs 获取歌单随机歌曲，限制100首，用于分析
	playlistDetail, err := model.GetPlaylistRandomSongs(db, idStr, 100)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &g.ErrDbOp
		}
		return "", &g.ErrDbOp
	}

	user := GetCurrentUser(c)
	// 权限检查
	allowed := false
	isAdmin := user != nil && user.UserGroup == "admin"
	isOwner := user != nil && playlistDetail.OwnerID == user.ID

	if isAdmin {
		if playlistDetail.IsPublic || isOwner {
			allowed = true
		}
	} else {
		if isOwner && !playlistDetail.IsPublic {
			allowed = true
		}
	}

	if !allowed {
		return "", &g.ErrPermission
	}

	// 读取 prompt1.md
	promptBytes, err := os.ReadFile("prompts/prompt_playlist1.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_playlist1.md", "error", err)
		return "", &g.Err
	}
	promptTemplate := string(promptBytes)

	// 构造 AI 分析所需的简化数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range playlistDetail.Songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return "", &g.ErrRequest
	}

	// 替换 prompt 中的占位符
	finalPrompt := strings.Replace(promptTemplate, "{{json数据}}", string(jsonData), 1)

	// 调用 AI 接口
	reply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		return "", &g.Err
	}
	return reply, nil
}

// GetPlaylistAIDescription 获取歌单AI描述
func (h *SongAuth) GetPlaylistAIDescription(c *gin.Context) {
	idStr := c.Param("id")

	// 1. 获取分析结果
	analysisReply, errRes := h.internalGetAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 2.// 读取 prompt2.md
	promptBytes, err := os.ReadFile("prompts/prompt_playlist2.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_playlist2.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板")
		return
	}
	promptTemplate := string(promptBytes)

	// 3. 替换 prompt2 中的占位符 {{summary_json}}
	finalPrompt := strings.Replace(promptTemplate, "{{summary_json}}", analysisReply, 1)

	// 4. 调用 AI 接口生成描述
	description, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 5. 保存描述到数据库
	db := GetDB(c)
	if err := db.Model(&model.Playlist{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		slog.Error("Failed to update playlist description", "error", err)
		// 即使保存失败，也返回生成的描述给前端
	}

	ReturnSuccess(c, gin.H{
		"description": description,
	})
}

// GenerateAllPublicPlaylistsDescription 生成所有公共歌单的AI描述
func (h *SongAuth) GenerateAllPublicPlaylistsDescription(c *gin.Context) {
	// 1. 管理员权限检查
	user := GetCurrentUser(c)
	if user == nil || user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "只有管理员可以执行此操作")
		return
	}

	db := GetDB(c)

	// 2. 获取所有公共歌单 ID
	var playlistIDs []int
	if err := db.Model(&model.Playlist{}).Where("is_public = ?", true).Pluck("id", &playlistIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "获取公共歌单失败")
		return
	}

	// 3. 异步执行生成任务，避免阻塞 HTTP 请求
	go func(ids []int) {
		slog.Info("开始批量生成公共歌单描述", "count", len(ids))
		for _, id := range ids {
			idStr := strconv.Itoa(id)
			// 复用 existing logic
			// 注意：internalGetAnalysis 包含权限检查，但管理员有权访问所有公共歌单，所以没问题
			// 为了复用，我们需要构造一个假的 context 或者重构 internalGetAnalysis 不依赖 context
			// 这里简单起见，我们重构 internalGetAnalysis 的逻辑，或者直接在此处调用核心逻辑

			// 为了代码复用，我们提取核心生成逻辑
			if err := h.generateAndSaveDescription(db, idStr, user); err != nil {
				slog.Error("生成歌单描述失败", "id", id, "error", err)
			} else {
				slog.Info("生成歌单描述成功", "id", id)
			}
		}
		slog.Info("批量生成公共歌单描述完成")
	}(playlistIDs)

	ReturnSuccess(c, gin.H{
		"message": fmt.Sprintf("已开始后台生成 %d 个公共歌单的描述", len(playlistIDs)),
	})
}

// generateAndSaveDescription 生成并保存单个歌单描述 (内部核心逻辑)
func (h *SongAuth) generateAndSaveDescription(db *gorm.DB, idStr string, user *model.User) error {
	// 1. 获取歌单详情 (Limit 100)
	playlistDetail, err := model.GetPlaylistRandomSongs(db, idStr, 100)
	if err != nil {
		return err
	}

	// 2. 权限检查 (复用逻辑)
	allowed := false
	isAdmin := user != nil && user.UserGroup == "admin"
	isOwner := user != nil && playlistDetail.OwnerID == user.ID

	if isAdmin {
		if playlistDetail.IsPublic || isOwner {
			allowed = true
		}
	} else {
		if isOwner && !playlistDetail.IsPublic {
			allowed = true
		}
	}

	if !allowed {
		return fmt.Errorf("permission denied")
	}

	// 3. 读取 prompt1.md
	promptBytes, err := os.ReadFile("prompts/prompt_playlist1.md")
	if err != nil {
		return err
	}
	promptTemplate1 := string(promptBytes)

	// 4. 构造 AI 分析所需的简化数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range playlistDetail.Songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return err
	}

	// 5. 生成分析 (Step 1)
	finalPrompt1 := strings.Replace(promptTemplate1, "{{json数据}}", string(jsonData), 1)
	analysisReply, err := ai.ChatWithQwen(finalPrompt1)
	if err != nil {
		return err
	}

	// 6. 读取 prompt2.md
	promptBytes2, err := os.ReadFile("prompts/prompt_playlist2.md")
	if err != nil {
		return err
	}
	promptTemplate2 := string(promptBytes2)

	// 7. 生成描述 (Step 2)
	finalPrompt2 := strings.Replace(promptTemplate2, "{{summary_json}}", analysisReply, 1)
	description, err := ai.ChatWithQwen(finalPrompt2)
	if err != nil {
		return err
	}

	// 8. 保存描述
	if err := db.Model(&model.Playlist{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		return err
	}

	return nil
}

// GetArtistAIAnalysis 获取艺术家AI分析
func (h *SongAuth) GetArtistAIAnalysis(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 调用内部方法获取分析结果
	analysisReply, errRes := h.internalGetArtistAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 尝试解析返回的 JSON
	var result interface{}
	cleanReply := strings.TrimSpace(analysisReply)
	cleanReply = strings.TrimPrefix(cleanReply, "```json")
	cleanReply = strings.TrimPrefix(cleanReply, "```")
	cleanReply = strings.TrimSuffix(cleanReply, "```")

	if err := json.Unmarshal([]byte(cleanReply), &result); err == nil {
		ReturnSuccess(c, result)
	} else {
		ReturnSuccess(c, gin.H{
			"analysis": analysisReply,
		})
	}
}

// GetArtistDetailRandom 获取艺术家详情（随机返回歌曲，上限100首）
func (*SongAuth) GetArtistDetailRandom(c *gin.Context) {
	idStr := c.Param("id")
	limitStr := c.Query("limit")

	limit := 100
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > 100 {
			limit = 100
		} else {
			limit = l
		}
	}

	db := GetDB(c)

	// 1. 获取艺术家信息
	var artist model.Artist
	if err := db.First(&artist, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "艺术家不存在")
		return
	}

	// 2. 获取随机歌曲
	songs, err := model.GetArtistRandomSongs(db, idStr, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 3. 构造返回数据
	type SongInfo struct {
		Title      string `json:"title"`
		Year       string `json:"year"`
		ArtistName string `json:"artist"`
		AlbumName  string `json:"album"`
	}

	var songInfos []SongInfo
	for _, song := range songs {
		songInfos = append(songInfos, SongInfo{
			Title:      song.Title,
			Year:       song.Year,
			ArtistName: song.ArtistName,
			AlbumName:  song.AlbumName,
		})
	}

	ReturnSuccess(c, gin.H{
		"artist": artist.Name,
		"songs":  songInfos,
	})
}

// GetAlbumDetailRandom 获取专辑详情（随机返回歌曲，上限100首）
func (*SongAuth) GetAlbumDetailRandom(c *gin.Context) {
	idStr := c.Param("id")
	limitStr := c.Query("limit")

	limit := 100
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > 100 {
			limit = 100
		} else {
			limit = l
		}
	}

	db := GetDB(c)

	// 1. 获取专辑信息
	var album model.Album
	if err := db.First(&album, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "专辑不存在")
		return
	}

	// 2. 获取随机歌曲
	songs, err := model.GetAlbumRandomSongs(db, idStr, limit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 3. 构造返回数据
	type SongInfo struct {
		Title      string `json:"title"`
		Year       string `json:"year"`
		ArtistName string `json:"artist"`
		AlbumName  string `json:"album"`
	}

	var songInfos []SongInfo
	for _, song := range songs {
		songInfos = append(songInfos, SongInfo{
			Title:      song.Title,
			Year:       song.Year,
			ArtistName: song.ArtistName,
			AlbumName:  song.AlbumName,
		})
	}

	ReturnSuccess(c, gin.H{
		"album": album.Title,
		"songs": songInfos,
	})
}

// internalGetArtistAnalysis 内部方法：获取艺术家分析结果
func (h *SongAuth) internalGetArtistAnalysis(c *gin.Context, idStr string) (string, *g.Result) {
	db := GetDB(c)

	// 1. 获取艺术家随机歌曲 (Limit 100)
	songs, err := model.GetArtistRandomSongs(db, idStr, 100)
	if err != nil {
		return "", &g.ErrDbOp
	}

	if len(songs) == 0 {
		return "", &g.Err
	}

	// 2. 读取 prompt_artist1.md
	promptBytes, err := os.ReadFile("prompts/prompt_artist1.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_artist1.md", "error", err)
		return "", &g.Err
	}
	promptTemplate := string(promptBytes)

	// 3. 构造 AI 分析所需的简化数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return "", &g.ErrRequest
	}

	// 4. 替换 prompt 中的占位符
	finalPrompt := strings.Replace(promptTemplate, "{{json数据}}", string(jsonData), 1)

	// 5. 调用 AI 接口
	reply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		return "", &g.Err
	}
	return reply, nil
}

// GetArtistAIDescription 获取艺术家AI描述
func (h *SongAuth) GetArtistAIDescription(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 获取分析结果
	analysisReply, errRes := h.internalGetArtistAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 3. 读取 prompt_artist2.md
	promptBytes, err := os.ReadFile("prompts/prompt_artist2.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_artist2.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板")
		return
	}
	promptTemplate := string(promptBytes)

	// 4. 替换 prompt_artist2 中的占位符 {{summary_json}}
	finalPrompt := strings.Replace(promptTemplate, "{{summary_json}}", analysisReply, 1)

	// 5. 调用 AI 接口生成描述
	description, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 6. 保存描述到数据库
	db := GetDB(c)
	if err := db.Model(&model.Artist{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		slog.Error("Failed to update artist description", "error", err)
	}

	ReturnSuccess(c, gin.H{
		"description": description,
	})
}

// GenerateAllArtistDescriptions 批量生成艺术家AI描述
func (h *SongAuth) GenerateAllArtistDescriptions(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)

	// 2. 获取所有艺术家 ID
	var artistIDs []int
	if err := db.Model(&model.Artist{}).Pluck("id", &artistIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "获取艺术家失败")
		return
	}

	// 3. 异步执行
	go func(ids []int) {
		slog.Info("开始批量生成艺术家描述", "count", len(ids))
		for _, id := range ids {
			idStr := strconv.Itoa(id)
			if err := h.generateAndSaveArtistDescription(db, idStr); err != nil {
				slog.Error("生成艺术家描述失败", "id", id, "error", err)
			} else {
				slog.Info("生成艺术家描述成功", "id", id)
			}
		}
		slog.Info("批量生成艺术家描述完成")
	}(artistIDs)

	ReturnSuccess(c, gin.H{
		"message": fmt.Sprintf("已开始后台生成 %d 个艺术家的描述", len(artistIDs)),
	})
}

// generateAndSaveArtistDescription 生成并保存单个艺术家描述 (内部核心逻辑)
func (h *SongAuth) generateAndSaveArtistDescription(db *gorm.DB, idStr string) error {
	// 1. 获取艺术家随机歌曲
	songs, err := model.GetArtistRandomSongs(db, idStr, 100)
	if err != nil {
		return err
	}
	if len(songs) == 0 {
		return fmt.Errorf("no songs found")
	}

	// 2. 读取 prompt_artist1.md
	promptBytes1, err := os.ReadFile("prompts/prompt_artist1.md")
	if err != nil {
		return err
	}
	promptTemplate1 := string(promptBytes1)

	// 3. 构造 AI 分析数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return err
	}

	// 4. Step 1: 分析
	finalPrompt1 := strings.Replace(promptTemplate1, "{{json数据}}", string(jsonData), 1)
	analysisReply, err := ai.ChatWithQwen(finalPrompt1)
	if err != nil {
		return err
	}

	// 5. 读取 prompt_artist2.md
	promptBytes2, err := os.ReadFile("prompts/prompt_artist2.md")
	if err != nil {
		return err
	}
	promptTemplate2 := string(promptBytes2)

	// 6. Step 2: 描述
	finalPrompt2 := strings.Replace(promptTemplate2, "{{summary_json}}", analysisReply, 1)
	description, err := ai.ChatWithQwen(finalPrompt2)
	if err != nil {
		return err
	}

	// 7. 保存
	if err := db.Model(&model.Artist{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		return err
	}

	return nil
}

// internalGetAlbumAnalysis 内部方法：获取专辑分析结果
func (h *SongAuth) internalGetAlbumAnalysis(c *gin.Context, idStr string) (string, *g.Result) {
	db := GetDB(c)

	// 1. 获取专辑随机歌曲 (Limit 100)
	songs, err := model.GetAlbumRandomSongs(db, idStr, 100)
	if err != nil {
		return "", &g.ErrDbOp
	}

	if len(songs) == 0 {
		return "", &g.Err
	}

	// 2. 读取 prompt_album1.md
	promptBytes, err := os.ReadFile("prompts/prompt_album1.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_album1.md", "error", err)
		return "", &g.Err
	}
	promptTemplate := string(promptBytes)

	// 3. 构造 AI 分析所需的简化数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return "", &g.ErrRequest
	}

	// 4. 替换 prompt 中的占位符
	// 注意：prompt_album1.md 应该包含 {{json数据}} 占位符，
	// 但用户提供的 prompt_album1.md 内容似乎没有显示该占位符？
	// 假设用户会添加或者默认追加。为了稳妥，我们检查一下。
	// 根据用户提供的 prompt_artist1.md，它是包含 {{json数据}} 的。
	// 这里我们假设 prompt_album1.md 也需要这个占位符。
	// 如果原文件没有，我们可能需要追加。
	// 暂时假设用户提供的 prompt 模板是完整的或我们将数据附在最后。
	// 不过根据之前的逻辑，应该是替换 {{json数据}}。
	// 让我们假设 prompt_album1.md 应该有。如果没有，我们直接附在后面。
	finalPrompt := promptTemplate
	if strings.Contains(promptTemplate, "{{json数据}}") {
		finalPrompt = strings.Replace(promptTemplate, "{{json数据}}", string(jsonData), 1)
	} else {
		finalPrompt = promptTemplate + "\n\n专辑数据：\n" + string(jsonData)
	}

	// 5. 调用 AI 接口
	reply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		return "", &g.Err
	}
	return reply, nil
}

// GetAlbumAIAnalysis 获取专辑AI分析
func (h *SongAuth) GetAlbumAIAnalysis(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 调用内部方法获取分析结果
	analysisReply, errRes := h.internalGetAlbumAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 尝试解析返回的 JSON
	var result interface{}
	cleanReply := strings.TrimSpace(analysisReply)
	cleanReply = strings.TrimPrefix(cleanReply, "```json")
	cleanReply = strings.TrimPrefix(cleanReply, "```")
	cleanReply = strings.TrimSuffix(cleanReply, "```")

	if err := json.Unmarshal([]byte(cleanReply), &result); err == nil {
		ReturnSuccess(c, result)
	} else {
		ReturnSuccess(c, gin.H{
			"analysis": analysisReply,
		})
	}
}

// GetAlbumAIDescription 获取专辑AI描述
func (h *SongAuth) GetAlbumAIDescription(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 获取分析结果
	analysisReply, errRes := h.internalGetAlbumAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 3. 读取 prompt_album2.md
	promptBytes, err := os.ReadFile("prompts/prompt_album2.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_album2.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板")
		return
	}
	promptTemplate := string(promptBytes)

	// 4. 替换 prompt_album2 中的占位符 {{summary_json}}
	finalPrompt := strings.Replace(promptTemplate, "{{summary_json}}", analysisReply, 1)

	// 5. 调用 AI 接口生成描述
	description, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 6. 保存描述到数据库
	db := GetDB(c)
	if err := db.Model(&model.Album{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		slog.Error("Failed to update album description", "error", err)
	}

	ReturnSuccess(c, gin.H{
		"description": description,
	})
}

// GenerateAllAlbumDescriptions 批量生成专辑AI描述
func (h *SongAuth) GenerateAllAlbumDescriptions(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)

	// 2. 获取所有专辑 ID
	var albumIDs []int
	if err := db.Model(&model.Album{}).Pluck("id", &albumIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "获取专辑失败")
		return
	}

	// 3. 异步执行
	go func(ids []int) {
		slog.Info("开始批量生成专辑描述", "count", len(ids))
		for _, id := range ids {
			idStr := strconv.Itoa(id)
			if err := h.generateAndSaveAlbumDescription(db, idStr); err != nil {
				slog.Error("生成专辑描述失败", "id", id, "error", err)
			} else {
				slog.Info("生成专辑描述成功", "id", id)
			}
		}
		slog.Info("批量生成专辑描述完成")
	}(albumIDs)

	ReturnSuccess(c, gin.H{
		"message": fmt.Sprintf("已开始后台生成 %d 个专辑的描述", len(albumIDs)),
	})
}

// generateAndSaveAlbumDescription 生成并保存单个专辑描述 (内部核心逻辑)
func (h *SongAuth) generateAndSaveAlbumDescription(db *gorm.DB, idStr string) error {
	// 1. 获取专辑随机歌曲
	songs, err := model.GetAlbumRandomSongs(db, idStr, 100)
	if err != nil {
		return err
	}
	if len(songs) == 0 {
		return fmt.Errorf("no songs found")
	}

	// 2. 读取 prompt_album1.md
	promptBytes1, err := os.ReadFile("prompts/prompt_album1.md")
	if err != nil {
		return err
	}
	promptTemplate1 := string(promptBytes1)

	// 3. 构造 AI 分析数据
	type AISongData struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Year   string `json:"year"`
	}
	var aiSongs []AISongData
	for _, s := range songs {
		aiSongs = append(aiSongs, AISongData{
			Title:  s.Title,
			Artist: s.ArtistName,
			Album:  s.AlbumName,
			Year:   s.Year,
		})
	}

	jsonData, err := json.Marshal(aiSongs)
	if err != nil {
		return err
	}

	// 4. Step 1: 分析
	finalPrompt1 := promptTemplate1
	if strings.Contains(promptTemplate1, "{{json数据}}") {
		finalPrompt1 = strings.Replace(promptTemplate1, "{{json数据}}", string(jsonData), 1)
	} else {
		finalPrompt1 = promptTemplate1 + "\n\n专辑数据：\n" + string(jsonData)
	}

	analysisReply, err := ai.ChatWithQwen(finalPrompt1)
	if err != nil {
		return err
	}

	// 5. 读取 prompt_album2.md
	promptBytes2, err := os.ReadFile("prompts/prompt_album2.md")
	if err != nil {
		return err
	}
	promptTemplate2 := string(promptBytes2)

	// 6. Step 2: 描述
	finalPrompt2 := strings.Replace(promptTemplate2, "{{summary_json}}", analysisReply, 1)
	description, err := ai.ChatWithQwen(finalPrompt2)
	if err != nil {
		return err
	}

	// 7. 保存
	if err := db.Model(&model.Album{}).Where("id = ?", idStr).Update("description", description).Error; err != nil {
		return err
	}

	return nil
}

// GetSongLyric 获取歌曲歌词
func (*SongAuth) GetSongLyric(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的ID")
		return
	}

	db := GetDB(c)
	var song model.Song
	if err := db.First(&song, id).Error; err != nil {
		ReturnError(c, g.Err, "歌曲不存在")
		return
	}

	// 1. 尝试读取同名 .lrc 文件
	ext := filepath.Ext(song.FilePath)
	basePath := strings.TrimSuffix(song.FilePath, ext)
	lrcPath := basePath + ".lrc"
	transPath := basePath + ".fy.lrc" // 约定翻译歌词文件名为 .fy.lrc

	lrcContent, errLrc := os.ReadFile(lrcPath)
	transContent, errTrans := os.ReadFile(transPath)

	if errLrc == nil {
		tlyric := ""
		if errTrans == nil {
			tlyric = string(transContent)
		}

		ReturnSuccess(c, gin.H{
			"lrc":    gin.H{"lyric": string(lrcContent)},
			"tlyric": gin.H{"lyric": tlyric},
			"source": "file",
		})
		return
	}

	// 2. 尝试读取内嵌歌词
	f, err := os.Open(song.FilePath)
	if err != nil {
		ReturnError(c, g.Err, "无法打开文件")
		return
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		// 无法读取元数据，返回无歌词
		ReturnSuccess(c, gin.H{
			"lrc":    gin.H{"lyric": "[00:00.000] 暂无歌词"},
			"tlyric": gin.H{"lyric": ""},
		})
		return
	}

	lyric := m.Lyrics()
	if lyric == "" {
		ReturnSuccess(c, gin.H{
			"lrc":    gin.H{"lyric": "[00:00.000] 暂无歌词"},
			"tlyric": gin.H{"lyric": ""},
		})
		return
	}

	ReturnSuccess(c, gin.H{
		"lrc":    gin.H{"lyric": lyric},
		"tlyric": gin.H{"lyric": ""},
		"source": "tag",
	})
}

// SubscribePlaylist 收藏歌单
func (*SongAuth) SubscribePlaylist(c *gin.Context) {
	idStr := c.Param("id")
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	playlistID, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的歌单ID")
		return
	}

	db := GetDB(c)
	var playlist model.Playlist
	if err := db.First(&playlist, playlistID).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	if err := db.Model(user).Association("SubscribedPlaylists").Append(&playlist); err != nil {
		ReturnError(c, g.ErrDbOp, "收藏失败")
		return
	}

	ReturnSuccess(c, gin.H{"message": "收藏成功"})
}

// UnsubscribePlaylist 取消收藏歌单
func (*SongAuth) UnsubscribePlaylist(c *gin.Context) {
	idStr := c.Param("id")
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	playlistID, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的歌单ID")
		return
	}

	db := GetDB(c)
	var playlist model.Playlist
	if err := db.First(&playlist, playlistID).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	if err := db.Model(user).Association("SubscribedPlaylists").Delete(&playlist); err != nil {
		ReturnError(c, g.ErrDbOp, "取消收藏失败")
		return
	}

	ReturnSuccess(c, gin.H{"message": "取消收藏成功"})
}

// GetSubscribedPlaylists 获取收藏的歌单
func (*SongAuth) GetSubscribedPlaylists(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)
	playlists, err := model.GetSubscribedPlaylists(db, user.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, "获取收藏列表失败")
		return
	}

	ReturnSuccess(c, playlists)
}

// CheckIsSubscribed 检查是否已收藏
func (*SongAuth) CheckIsSubscribed(c *gin.Context) {
	idStr := c.Param("id")
	user := GetCurrentUser(c)
	if user == nil {
		// 未登录视为未收藏，但不报错，方便前端处理
		ReturnSuccess(c, gin.H{"is_subscribed": false})
		return
	}

	playlistID, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的歌单ID")
		return
	}

	db := GetDB(c)
	isSubscribed, err := model.IsPlaylistSubscribed(db, user.ID, playlistID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{"is_subscribed": isSubscribed})
}

// CreatePlaylistRequest 创建歌单请求
type CreatePlaylistRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// CreatePrivatePlaylist 创建私有歌单
func (*SongAuth) CreatePrivatePlaylist(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	var req CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查歌单是否已存在
	var count int64
	if err := db.Model(&model.Playlist{}).Where("title = ? AND owner_id = ?", req.Title, user.ID).Count(&count).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if count > 0 {
		ReturnError(c, g.ErrRequest, "歌单已存在")
		return
	}

	// 创建私有歌单 (IsPublic = false)
	playlist, err := model.CreatePlaylist(db, user.ID, req.Title, req.Description, false)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlist)
}

// AddSongsRequest 批量添加歌曲请求
type AddSongsRequest struct {
	PlaylistID int   `json:"playlist_id" binding:"required"`
	SongIDs    []int `json:"song_ids" binding:"required"`
}

// AddSongsToPlaylist 批量添加歌曲到歌单
func (*SongAuth) AddSongsToPlaylist(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	var req AddSongsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if len(req.SongIDs) == 0 {
		ReturnSuccess(c, "未选择任何歌曲")
		return
	}

	db := GetDB(c)

	// 1. 检查歌单权限
	var playlist model.Playlist
	if err := db.First(&playlist, req.PlaylistID).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	if playlist.OwnerID != user.ID {
		ReturnError(c, g.ErrPermission, "无权修改此歌单")
		return
	}

	// 2. 查找所有歌曲
	var songs []model.Song
	if err := db.Find(&songs, req.SongIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "查询歌曲失败")
		return
	}

	if len(songs) == 0 {
		ReturnError(c, g.ErrRequest, "未找到有效的歌曲")
		return
	}

	// 过滤已存在的歌曲
	var existingSongIDs []int
	// 假设中间表名为 playlist_songs
	if err := db.Table("playlist_songs").Where("playlist_id = ? AND song_id IN ?", playlist.ID, req.SongIDs).Pluck("song_id", &existingSongIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "查询歌单歌曲失败")
		return
	}

	existingMap := make(map[int]bool)
	for _, id := range existingSongIDs {
		existingMap[id] = true
	}

	var newSongs []model.Song
	for _, song := range songs {
		if !existingMap[song.ID] {
			newSongs = append(newSongs, song)
		}
	}

	if len(newSongs) == 0 {
		ReturnSuccess(c, gin.H{
			"message": "所选歌曲已在歌单中",
			"count":   0,
		})
		return
	}

	// 3. 批量添加
	if err := model.AddSongsToPlaylist(db, &playlist, newSongs); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 4. 设置歌单封面（始终使用最新添加歌曲的封面）
	for _, s := range newSongs {
		// 需要加载 Song 的 Cover 关联，因为 newSongs 此时可能只有基本信息
		var fullSong model.Song
		if err := db.Preload("Cover").First(&fullSong, s.ID).Error; err == nil {
			if fullSong.CoverID != nil && fullSong.Cover.ID != 0 {
				coverUrl := fullSong.Cover.Path
				if err := db.Model(&playlist).Update("cover_url", coverUrl).Error; err == nil {
					break // 成功设置后退出
				}
			}
		}
	}

	ReturnSuccess(c, gin.H{
		"message": "添加成功",
		"count":   len(newSongs),
	})
}

// RemoveSongsRequest 批量移除歌曲请求
type RemoveSongsRequest struct {
	PlaylistID int   `json:"playlist_id" binding:"required"`
	SongIDs    []int `json:"song_ids" binding:"required"`
}

// RemoveSongsFromPlaylist 从歌单中批量移除歌曲
func (*SongAuth) RemoveSongsFromPlaylist(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	var req RemoveSongsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if len(req.SongIDs) == 0 {
		ReturnSuccess(c, "未选择任何歌曲")
		return
	}

	db := GetDB(c)

	// 1. 检查歌单权限
	var playlist model.Playlist
	if err := db.First(&playlist, req.PlaylistID).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	// 检查是否为公共歌单
	if playlist.IsPublic {
		ReturnError(c, g.ErrPermission, "公共歌单禁止删除歌曲")
		return
	}

	// 检查所有权
	if playlist.OwnerID != user.ID {
		ReturnError(c, g.ErrPermission, "无权操作此歌单")
		return
	}

	// 2. 查找歌曲 (只需要 ID 即可，GORM Delete Association 需要 Model 对象)
	var songs []model.Song
	if err := db.Find(&songs, req.SongIDs).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "查询歌曲失败")
		return
	}

	if len(songs) == 0 {
		ReturnError(c, g.ErrRequest, "未找到有效的歌曲")
		return
	}

	// 3. 批量移除
	if err := model.RemoveSongsFromPlaylist(db, &playlist, songs); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"message": "移除成功",
		"count":   len(songs),
	})
}

// DeletePrivatePlaylist 删除私有歌单
func (*SongAuth) DeletePrivatePlaylist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的ID")
		return
	}

	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	db := GetDB(c)

	// 1. 查找歌单
	var playlist model.Playlist
	if err := db.First(&playlist, id).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	// 2. 检查权限 (必须是私有歌单且是所有者)
	if playlist.IsPublic {
		ReturnError(c, g.ErrPermission, "无法删除公共歌单")
		return
	}

	if playlist.OwnerID != user.ID {
		ReturnError(c, g.ErrPermission, "无权删除此歌单")
		return
	}

	// 3. 执行删除
	if err := model.DeletePlaylist(db, id); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, "删除成功")
}

// internalGetSongAnalysis 内部方法：获取歌曲分析结果
func (h *SongAuth) internalGetSongAnalysis(c *gin.Context, idStr string) (string, *g.Result) {
	db := GetDB(c)

	// 1. 获取歌曲信息
	song, err := model.GetSongByID(db, idStr)
	if err != nil {
		return "", &g.ErrDbOp
	}

	// 2. 读取 prompts/prompt_podcast1.md
	promptBytes, err := os.ReadFile("prompts/prompt_podcast1.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_podcast1.md", "error", err)
		return "", &g.Err
	}
	promptTemplate := string(promptBytes)

	// 3. 构造 AI 分析所需的简化数据
	type SongData struct {
		Title  string `json:"title"`
		Year   string `json:"year"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	}
	songData := SongData{
		Title:  song.Title,
		Year:   song.Year,
		Artist: song.ArtistName,
		Album:  song.AlbumName,
	}

	jsonData, err := json.Marshal(songData)
	if err != nil {
		return "", &g.ErrRequest
	}

	// 4. 替换 prompt 中的占位符
	finalPrompt := strings.Replace(promptTemplate, "{{song_json}}", string(jsonData), 1)

	// 5. 调用 AI 接口
	reply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		return "", &g.Err
	}
	return reply, nil
}

// GetSongAIAnalysis 获取歌曲AI分析 (Step 1)
func (h *SongAuth) GetSongAIAnalysis(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 调用内部方法获取分析结果
	analysisReply, errRes := h.internalGetSongAnalysis(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取分析失败")
		return
	}

	// 尝试解析返回的 JSON
	var result interface{}
	cleanReply := strings.TrimSpace(analysisReply)
	cleanReply = strings.TrimPrefix(cleanReply, "```json")
	cleanReply = strings.TrimPrefix(cleanReply, "```")
	cleanReply = strings.TrimSuffix(cleanReply, "```")

	if err := json.Unmarshal([]byte(cleanReply), &result); err == nil {
		ReturnSuccess(c, result)
	} else {
		ReturnSuccess(c, gin.H{
			"analysis": analysisReply,
		})
	}
}

// internalGetSongDraft 内部方法：获取开场白草稿 (Step 2)
func (h *SongAuth) internalGetSongDraft(c *gin.Context, idStr string) (string, *g.Result) {
	// 1. 获取分析结果 (Step 1)
	analysisReply, errRes := h.internalGetSongAnalysis(c, idStr)
	if errRes != nil {
		return "", errRes
	}

	// 2. 读取 prompts/prompt_podcast2.md
	promptBytes2, err := os.ReadFile("prompts/prompt_podcast2.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_podcast2.md", "error", err)
		return "", &g.Err
	}
	promptTemplate2 := string(promptBytes2)

	// 3. 生成开场白草稿 (Step 2)
	finalPrompt2 := strings.Replace(promptTemplate2, "{{analysis_json}}", analysisReply, 1)
	openingDraft, err := ai.ChatWithQwen(finalPrompt2)
	if err != nil {
		return "", &g.Err
	}
	return openingDraft, nil
}

// GetSongAIDraft 获取歌曲AI开场白草稿 (Step 2)
func (h *SongAuth) GetSongAIDraft(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 获取草稿
	openingDraft, errRes := h.internalGetSongDraft(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取草稿失败")
		return
	}

	ReturnSuccess(c, gin.H{
		"draft": openingDraft,
	})
}

// GetSongAIOpeningRemark 获取歌曲AI最终开场白 (Step 3)
func (h *SongAuth) GetSongAIOpeningRemark(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 获取草稿 (Step 2)
	openingDraft, errRes := h.internalGetSongDraft(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取草稿失败")
		return
	}

	// 3. 读取 prompts/prompt_podcast3.md
	promptBytes3, err := os.ReadFile("prompts/prompt_podcast3.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_podcast3.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板3")
		return
	}
	promptTemplate3 := string(promptBytes3)

	// 4. 生成最终语音脚本 (Step 3)
	// 注意：prompt_podcast3.md 中的占位符是 {{analysis_json}}，但实际上我们传入的是 Step 2 生成的文本
	finalPrompt3 := strings.Replace(promptTemplate3, "{{analysis_json}}", openingDraft, 1)
	finalScript, err := ai.ChatWithQwen(finalPrompt3)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	ReturnSuccess(c, gin.H{
		"opening_remark": finalScript,
		"draft":          openingDraft, // 可选：同时也返回草稿方便调试
	})
}

// GetSongAIOpeningRemarkTTS 获取歌曲AI开场白并转语音 (Step 4)
func (h *SongAuth) GetSongAIOpeningRemarkTTS(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")

	// 2. 获取草稿 (Step 2)
	openingDraft, errRes := h.internalGetSongDraft(c, idStr)
	if errRes != nil {
		ReturnError(c, *errRes, "获取草稿失败")
		return
	}

	// 3. 读取 prompts/prompt_podcast3.md
	promptBytes3, err := os.ReadFile("prompts/prompt_podcast3.md")
	if err != nil {
		slog.Error("Failed to read prompts/prompt_podcast3.md", "error", err)
		ReturnError(c, g.Err, "无法读取提示词模板3")
		return
	}
	promptTemplate3 := string(promptBytes3)

	// 4. 生成最终语音脚本 (Step 3)
	finalPrompt3 := strings.Replace(promptTemplate3, "{{analysis_json}}", openingDraft, 1)
	finalScript, err := ai.ChatWithQwen(finalPrompt3)
	if err != nil {
		ReturnError(c, g.Err, err)
		return
	}

	// 5. 调用 TTS 生成语音 (Step 4)
	// 默认参数配置
	params := &ai.Params{
		Voice:      "Neil",
		Format:     "flac",
		SampleRate: 24000,
		Volume:     50,
		Rate:       1.0,
		Pitch:      1.0,
	}

	audioPath, err := ai.SynthesizeAudio(finalScript, params)
	if err != nil {
		slog.Error("TTS generation failed", "error", err)
		ReturnError(c, g.Err, "语音合成失败")
		return
	}

	// 6. 保存开场白和语音文件路径到数据库
	audioFileName := filepath.Base(audioPath)
	db := GetDB(c)
	updates := map[string]interface{}{
		"description":        finalScript,
		"opening_audio_file": audioFileName,
	}
	if err := db.Model(&model.Song{}).Where("id = ?", idStr).Updates(updates).Error; err != nil {
		slog.Error("Failed to update song description and audio file", "id", idStr, "error", err)
		// 不返回错误，因为生成已经成功
	}

	// 尝试找到该歌曲所属的歌单，并更新歌单的 HasIntro 为 true
	// 注意：一首歌曲可能属于多个歌单，这里我们更新包含该歌曲的所有歌单
	// 或者，我们可以只更新"当前上下文"的歌单，但接口没有传 playlist_id。
	// 简单起见，我们更新包含此歌曲的所有歌单的 HasIntro 为 true
	// 这样任何包含此歌曲的歌单都会被标记为有开场白（这是合理的，因为只要有一首有开场白，歌单就可以认为有开场白内容）
	var playlistIDs []int
	if err := db.Table("playlist_songs").Where("song_id = ?", idStr).Pluck("playlist_id", &playlistIDs).Error; err == nil && len(playlistIDs) > 0 {
		if err := db.Model(&model.Playlist{}).Where("id IN ?", playlistIDs).Update("has_intro", true).Error; err != nil {
			slog.Error("Failed to update playlists HasIntro status", "song_id", idStr, "error", err)
		}
	}

	ReturnSuccess(c, gin.H{
		"opening_remark": finalScript,
		"audio_path":     audioFileName,
	})
}

// BatchGenerateSongIntros 批量生成歌单内歌曲的开场白
func (h *SongAuth) BatchGenerateSongIntros(c *gin.Context) {
	// 1. 权限检查: 允许所有登录用户
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	idStr := c.Param("id")
	db := GetDB(c)

	// 2. 获取歌单内的所有歌曲
	// 这里我们需要获取歌单的所有歌曲，而不是随机 100 首
	// 复用 GetPlaylistDetail 的逻辑，但 limit 设大一点或者获取全部
	// 由于 GetPlaylistRandomSongs 可能会乱序，这里我们最好按顺序获取或者全部获取
	// 考虑到歌单可能很大，建议异步处理

	// 获取歌单信息用于权限检查
	var playlist model.Playlist
	if err := db.First(&playlist, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	// 权限检查
	allowed := false
	isAdmin := user.UserGroup == "admin"
	isOwner := playlist.OwnerID == user.ID

	if isAdmin {
		if playlist.IsPublic || isOwner {
			allowed = true
		}
	} else {
		if isOwner && !playlist.IsPublic {
			allowed = true
		}
	}

	if !allowed {
		ReturnError(c, g.ErrPermission, "无权访问此歌单")
		return
	}

	// 获取歌单所有歌曲 ID
	var songIDs []int
	err := db.Table("playlist_songs").Where("playlist_id = ?", idStr).Pluck("song_id", &songIDs).Error
	if err != nil {
		ReturnError(c, g.ErrDbOp, "获取歌单歌曲失败")
		return
	}

	if len(songIDs) == 0 {
		ReturnError(c, g.ErrRequest, "歌单为空")
		return
	}

	// 3. 异步执行批量生成任务
	go func(ids []int, userID int) {
		slog.Info("开始批量生成歌曲开场白", "playlist_id", idStr, "count", len(ids))
		for _, songID := range ids {
			// 检查是否已经生成过（可选，这里强制重新生成或者跳过已存在的）
			// 这里假设全部生成

			// 调用内部逻辑生成单个歌曲的开场白
			// 由于 internalGetSongDraft 等方法依赖 gin.Context，我们需要重构或者模拟
			// 或者提取核心逻辑为独立函数。
			// 这里我们直接调用核心生成逻辑

			if err := h.generateAndSaveSongIntro(db, strconv.Itoa(songID)); err != nil {
				slog.Error("生成歌曲开场白失败", "song_id", songID, "error", err)
			} else {
				slog.Info("生成歌曲开场白成功", "song_id", songID)
			}
		}
		slog.Info("批量生成歌曲开场白完成", "playlist_id", idStr)
	}(songIDs, user.ID)

	// 更新歌单 HasIntro 状态为 true
	if err := db.Model(&model.Playlist{}).Where("id = ?", idStr).Update("has_intro", true).Error; err != nil {
		slog.Error("Failed to update playlist HasIntro status", "id", idStr, "error", err)
	}

	ReturnSuccess(c, gin.H{
		"message": fmt.Sprintf("已开始后台生成 %d 首歌曲的开场白", len(songIDs)),
	})
}

// generateAndSaveSongIntro 生成并保存单个歌曲开场白 (内部核心逻辑)
func (h *SongAuth) generateAndSaveSongIntro(db *gorm.DB, idStr string) error {
	// 1. 获取歌曲信息
	song, err := model.GetSongByID(db, idStr)
	if err != nil {
		return err
	}

	// 2. Step 1: 分析
	promptBytes, err := os.ReadFile("prompts/prompt_podcast1.md")
	if err != nil {
		return err
	}
	promptTemplate := string(promptBytes)

	type SongData struct {
		Title  string `json:"title"`
		Year   string `json:"year"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	}
	songData := SongData{
		Title:  song.Title,
		Year:   song.Year,
		Artist: song.ArtistName,
		Album:  song.AlbumName,
	}
	jsonData, err := json.Marshal(songData)
	if err != nil {
		return err
	}
	finalPrompt := strings.Replace(promptTemplate, "{{song_json}}", string(jsonData), 1)
	analysisReply, err := ai.ChatWithQwen(finalPrompt)
	if err != nil {
		return err
	}

	// 3. Step 2: 草稿
	promptBytes2, err := os.ReadFile("prompts/prompt_podcast2.md")
	if err != nil {
		return err
	}
	promptTemplate2 := string(promptBytes2)
	finalPrompt2 := strings.Replace(promptTemplate2, "{{analysis_json}}", analysisReply, 1)
	openingDraft, err := ai.ChatWithQwen(finalPrompt2)
	if err != nil {
		return err
	}

	// 4. Step 3: 最终脚本
	promptBytes3, err := os.ReadFile("prompts/prompt_podcast3.md")
	if err != nil {
		return err
	}
	promptTemplate3 := string(promptBytes3)
	finalPrompt3 := strings.Replace(promptTemplate3, "{{analysis_json}}", openingDraft, 1)
	finalScript, err := ai.ChatWithQwen(finalPrompt3)
	if err != nil {
		return err
	}

	// 5. Step 4: TTS
	params := &ai.Params{
		Voice:      "Neil",
		Format:     "flac",
		SampleRate: 24000,
		Volume:     50,
		Rate:       1.0,
		Pitch:      1.0,
	}
	audioPath, err := ai.SynthesizeAudio(finalScript, params)
	if err != nil {
		return err
	}

	// 6. 保存
	audioFileName := filepath.Base(audioPath)
	updates := map[string]interface{}{
		"description":        finalScript,
		"opening_audio_file": audioFileName,
	}
	if err := db.Model(&model.Song{}).Where("id = ?", idStr).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DeletePublicPlaylist 删除公共歌单 (仅限管理员)
func (*SongAuth) DeletePublicPlaylist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, "无效的ID")
		return
	}

	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}

	// 检查是否为管理员
	if user.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, "无权操作")
		return
	}

	db := GetDB(c)

	// 1. 查找歌单
	var playlist model.Playlist
	if err := db.First(&playlist, id).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌单不存在")
		return
	}

	// 2. 检查类型
	if !playlist.IsPublic {
		ReturnError(c, g.ErrPermission, "此接口仅用于删除公共歌单")
		return
	}

	// 3. 执行删除
	if err := model.DeletePlaylist(db, id); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, "删除成功")
}
