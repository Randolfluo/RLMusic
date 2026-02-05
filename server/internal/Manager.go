package server

import (
	"server/internal/handle"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

const (
	apiBasePath = "/api"
)

// 全局路由处理器实例
var (
	userAuthAPI   = &handle.UserAuth{}
	fileAuthAPI   = &handle.FileAuth{}
	systemAuthAPI = &handle.SystemAuth{}
	songAuthAPI   = &handle.SongAuth{}
)

// RegisterHandlers 注册所有路由处理器
func RegisterHandlers(r *gin.Engine) {
	// 配置Swagger
	docs.SwaggerInfo.BasePath = apiBasePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerBaseHandler(r)
	registerAuthHandler(r)
}

// registerBaseHandler 注册基础路由,无需认证
func registerBaseHandler(r *gin.Engine) {
	base := r.Group(apiBasePath)

	// 用户认证相关
	auth := base.Group("/auth")
	{
		auth.POST("/login", userAuthAPI.Login)       // 用户登录
		auth.POST("/register", userAuthAPI.Register) // 用户注册
	}

	// 文件相关(无需认证)
	file := base.Group("/file")
	{
		file.POST("/initFolder", fileAuthAPI.InitFolder)
		file.GET("/avatar/:filename", userAuthAPI.GetAvatar)
	}

	// 歌曲相关(无需认证)
	song := base.Group("/song")
	{
		song.GET("/playlists/public", songAuthAPI.GetAllPlaylists)            // 获取所有公开歌单
		song.GET("/playlist/public/:id", songAuthAPI.GetPublicPlaylistDetail) // 获取公开歌单详情
		song.GET("/stream/:id", songAuthAPI.StreamSong)                       // 播放歌曲
		song.GET("/detail/:id", songAuthAPI.GetSongDetail)                    // 获取歌曲详情
		song.GET("/cover/:id", songAuthAPI.GetSongCover)                      // 获取歌曲封面
		song.GET("/lyric/:id", songAuthAPI.GetSongLyric)                      // 获取歌曲歌词
	}

	// 搜索相关(无需认证)
	search := base.Group("/search")
	{
		search.GET("/song", handle.SearchSong)         // 搜索歌曲
		search.GET("/artist", handle.SearchArtist)     // 搜索歌手
		search.GET("/album", handle.SearchAlbum)       // 搜索专辑
		search.GET("/playlist", handle.SearchPlaylist) // 搜索歌单
	}

	// 系统相关(无需认证)
	system := base.Group("/system")
	{
		system.GET("/stats", systemAuthAPI.GetStats) // 合并后的接口
	}
}

// registerAuthHandler 注册需要JWT认证的路由
func registerAuthHandler(r *gin.Engine) {
	auth := r.Group(apiBasePath)
	auth.Use(middleware.JWTAuth())

	// 用户认证相关(需要登录)
	userAuth := auth.Group("/auth")
	{
		userAuth.POST("/logout", userAuthAPI.Logout) // 退出登录
	}

	// 用户相关
	user := auth.Group("/user")
	{
		user.DELETE("", userAuthAPI.DeleteUser)        // 注销用户
		user.GET("/info", userAuthAPI.GetUserInfo)     // 获取用户信息
		user.POST("/avatar", userAuthAPI.UploadAvatar) // 上传头像
	}

	// 歌曲管理
	song := auth.Group("/song")
	{
		song.POST("/scan", songAuthAPI.ScanUserMusic)
		song.PUT("/playlist/:id", songAuthAPI.UpdatePlaylist) // 更新歌单信息
		// song.GET("/playlists/private", songAuthAPI.GetPrivatePlaylists) // Legacy (removed)
		song.GET("/playlists/user/public", songAuthAPI.GetUserPublicPlaylists)   // 获取用户公开歌单
		song.GET("/playlists/user/private", songAuthAPI.GetUserPrivatePlaylists) // 获取用户私有歌单
		song.GET("/playlist/private/:id", songAuthAPI.GetPrivatePlaylistDetail)  // 获取私有歌单详情
		song.POST("/like/:id", songAuthAPI.ToggleLike)                           // 点赞/取消点赞
		song.POST("/history", songAuthAPI.AddHistory)                            // 记录播放历史
		song.GET("/history", songAuthAPI.GetHistory)                             // 获取播放历史
		song.DELETE("/history", songAuthAPI.ClearHistory)                        // 清空播放历史
	}

	// 系统相关
	system := auth.Group("/system")
	{
		system.POST("/config", systemAuthAPI.UpdateConfig)
	}

}
