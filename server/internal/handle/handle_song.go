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
	var scannedDuration float64 = 0

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

		// 查找或创建歌单 (歌单名 = 文件夹名)
		// 移除权限控制，permission参数不再使用，默认公开
		currentPlaylist, err := model.FindOrCreatePlaylist(db, user.ID, subDirName)
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
				// 显式更新 Artists 关联 (针对更新场景)
				if !isCreated && len(song.Artists) > 0 {
					_ = db.Model(song).Association("Artists").Replace(song.Artists)
				}

				if isCreated {
					addedCount++
				} else {
					updatedCount++
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

	_ = model.UpdateSystemInfoStats(db, songCount, albumCount, artistCount, systemTotalDuration)

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
		// 404 Not Found if no cover
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
	if err := db.First(&artist, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "歌手不存在")
		return
	}
	ReturnSuccess(c, artist)
}

// GetAlbumDetail 获取专辑详情
func (*SongAuth) GetAlbumDetail(c *gin.Context) {
	idStr := c.Param("id")
	db := GetDB(c)
	var album model.Album
	if err := db.Preload("Artist").First(&album, idStr).Error; err != nil {
		ReturnError(c, g.ErrDbOp, "专辑不存在")
		return
	}
	ReturnSuccess(c, album)
}

// GetAllPlaylists 获取所有公开歌单
func (*SongAuth) GetAllPlaylists(c *gin.Context) {
	db := GetDB(c)

	playlists, err := model.GetPublicPlaylists(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
}

// GetUserPublicPlaylists 获取用户公开歌单
func (*SongAuth) GetUserPublicPlaylists(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}
	db := GetDB(c)

	playlists, err := model.GetUserPublicPlaylists(db, user.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
}

// GetUserPrivatePlaylists 获取用户私有歌单
func (*SongAuth) GetUserPrivatePlaylists(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		ReturnError(c, g.ErrUserNotExist, "用户未登录")
		return
	}
	db := GetDB(c)

	playlists, err := model.GetUserPrivatePlaylists(db, user.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, playlists)
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
