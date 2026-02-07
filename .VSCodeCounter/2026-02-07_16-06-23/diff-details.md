# Diff Details

Date : 2026-02-07 16:06:23

Directory c:\\Users\\lrj32\\Desktop\\music player\\localmusicplayer

Total : 107 files,  6694 codes, 342 comments, 768 blanks, all 7804 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [README.md](/README.md) | Markdown | 42 | 0 | 17 | 59 |
| [components.d.ts](/components.d.ts) | TypeScript | 35 | 0 | 0 | 35 |
| [dist-electron/main.js](/dist-electron/main.js) | JavaScript | -71 | -2 | -1 | -74 |
| [dist-electron/preload.mjs](/dist-electron/preload.mjs) | JavaScript | -20 | -1 | -1 | -22 |
| [electron/main.ts](/electron/main.ts) | TypeScript | 25 | 3 | 4 | 32 |
| [package.json](/package.json) | JSON | 3 | 0 | 0 | 3 |
| [pnpm-lock.yaml](/pnpm-lock.yaml) | YAML | 31 | 0 | 4 | 35 |
| [public/images/logo/favicon.png](/public/images/logo/favicon.png) | XML | -1 | 0 | 0 | -1 |
| [public/images/logo/favicon.svg](/public/images/logo/favicon.svg) | XML | 1 | 0 | 0 | 1 |
| [public/loading.html](/public/loading.html) | HTML | 70 | 0 | 4 | 74 |
| [server/.air.toml](/server/.air.toml) | TOML | 22 | 0 | 5 | 27 |
| [server/api\_test\_scan.py](/server/api_test_scan.py) | Python | 66 | 8 | 15 | 89 |
| [server/cmd/main.go](/server/cmd/main.go) | Go | 3 | 1 | 1 | 5 |
| [server/config.toml](/server/config.toml) | TOML | 23 | 0 | 6 | 29 |
| [server/go.mod](/server/go.mod) | Go Module File | 9 | 0 | 0 | 9 |
| [server/go.sum](/server/go.sum) | Go Checksum File | 22 | 0 | 0 | 22 |
| [server/internal/Helper.go](/server/internal/Helper.go) | Go | 32 | 0 | 3 | 35 |
| [server/internal/Manager.go](/server/internal/Manager.go) | Go | 27 | 2 | 4 | 33 |
| [server/internal/global/config.go](/server/internal/global/config.go) | Go | 6 | 1 | 1 | 8 |
| [server/internal/global/result.go](/server/internal/global/result.go) | Go | -4 | -1 | -1 | -6 |
| [server/internal/handle/handle\_auth.go](/server/internal/handle/handle_auth.go) | Go | 113 | 20 | 23 | 156 |
| [server/internal/handle/handle\_file.go](/server/internal/handle/handle_file.go) | Go | -44 | -8 | -9 | -61 |
| [server/internal/handle/handle\_history.go](/server/internal/handle/handle_history.go) | Go | 88 | 6 | 17 | 111 |
| [server/internal/handle/handle\_search.go](/server/internal/handle/handle_search.go) | Go | 91 | 4 | 17 | 112 |
| [server/internal/handle/handle\_song.go](/server/internal/handle/handle_song.go) | Go | 290 | 34 | 54 | 378 |
| [server/internal/handle/handle\_system.go](/server/internal/handle/handle_system.go) | Go | 105 | 12 | 20 | 137 |
| [server/internal/middleware/auth.go](/server/internal/middleware/auth.go) | Go | 0 | -35 | -4 | -39 |
| [server/internal/middleware/stats.go](/server/internal/middleware/stats.go) | Go | 12 | 1 | 4 | 17 |
| [server/internal/model/History.go](/server/internal/model/History.go) | Go | 53 | 5 | 11 | 69 |
| [server/internal/model/Metadata.go](/server/internal/model/Metadata.go) | Go | 50 | 10 | 8 | 68 |
| [server/internal/model/Playlist.go](/server/internal/model/Playlist.go) | Go | 168 | 28 | 25 | 221 |
| [server/internal/model/Search.go](/server/internal/model/Search.go) | Go | 60 | 9 | 22 | 91 |
| [server/internal/model/Song.go](/server/internal/model/Song.go) | Go | 21 | 4 | 4 | 29 |
| [server/internal/model/SystemInfo.go](/server/internal/model/SystemInfo.go) | Go | 68 | 6 | 11 | 85 |
| [server/internal/model/SystemSetting.go](/server/internal/model/SystemSetting.go) | Go | -26 | -5 | -6 | -37 |
| [server/internal/model/User.go](/server/internal/model/User.go) | Go | 28 | 9 | 5 | 42 |
| [server/internal/model/User\_ext.go](/server/internal/model/User_ext.go) | Go | -5 | -1 | -3 | -9 |
| [server/internal/model/User\_test.go](/server/internal/model/User_test.go) | Go | -44 | -3 | -12 | -59 |
| [server/internal/model/ZBase.go](/server/internal/model/ZBase.go) | Go | 1 | 0 | 0 | 1 |
| [server/internal/utils/audio/flac.go](/server/internal/utils/audio/flac.go) | Go | 83 | 36 | 22 | 141 |
| [server/internal/utils/encrypt/Aes.go](/server/internal/utils/encrypt/Aes.go) | Go | 44 | 4 | 11 | 59 |
| [server/internal/utils/imgtool/imgtool.go](/server/internal/utils/imgtool/imgtool.go) | Go | 50 | 14 | 12 | 76 |
| [server/requests.http](/server/requests.http) | HTTP | 48 | 187 | 25 | 260 |
| [server/test/apifox\_batch\_import.json](/server/test/apifox_batch_import.json) | JSON | 797 | 0 | 0 | 797 |
| [server/test/auth.http](/server/test/auth.http) | HTTP | 26 | 8 | 8 | 42 |
| [server/test/convert\_to\_postman.py](/server/test/convert_to_postman.py) | Python | 132 | 16 | 37 | 185 |
| [server/test/playlist.http](/server/test/playlist.http) | HTTP | 27 | 13 | 10 | 50 |
| [server/test/search.http](/server/test/search.http) | HTTP | 9 | 5 | 5 | 19 |
| [server/test/song.http](/server/test/song.http) | HTTP | 31 | 16 | 16 | 63 |
| [server/test/system.http](/server/test/system.http) | HTTP | 20 | 9 | 8 | 37 |
| [server/test/user.http](/server/test/user.http) | HTTP | 1 | 1 | 0 | 2 |
| [src/App.vue](/src/App.vue) | Vue | 19 | 0 | 2 | 21 |
| [src/api/album.js](/src/api/album.js) | JavaScript | -29 | -17 | -5 | -51 |
| [src/api/artist.js](/src/api/artist.js) | JavaScript | -91 | -43 | -9 | -143 |
| [src/api/comment.js](/src/api/comment.js) | JavaScript | -42 | -25 | -5 | -72 |
| [src/api/home.js](/src/api/home.js) | JavaScript | -80 | -35 | -11 | -126 |
| [src/api/login.js](/src/api/login.js) | JavaScript | -88 | -31 | -9 | -128 |
| [src/api/login.ts](/src/api/login.ts) | TypeScript | 30 | 15 | 7 | 52 |
| [src/api/playlist.js](/src/api/playlist.js) | JavaScript | -108 | -56 | -12 | -176 |
| [src/api/playlist.ts](/src/api/playlist.ts) | TypeScript | 46 | 26 | 9 | 81 |
| [src/api/search.js](/src/api/search.js) | JavaScript | -30 | -17 | -5 | -52 |
| [src/api/search.ts](/src/api/search.ts) | TypeScript | 30 | 12 | 10 | 52 |
| [src/api/song.js](/src/api/song.js) | JavaScript | -72 | -33 | -9 | -114 |
| [src/api/song.ts](/src/api/song.ts) | TypeScript | 91 | 50 | 13 | 154 |
| [src/api/system.ts](/src/api/system.ts) | TypeScript | 34 | 10 | 5 | 49 |
| [src/api/user.js](/src/api/user.js) | JavaScript | -142 | -56 | -15 | -213 |
| [src/api/user.ts](/src/api/user.ts) | TypeScript | 35 | 9 | 5 | 49 |
| [src/api/video.js](/src/api/video.js) | JavaScript | -30 | -16 | -5 | -51 |
| [src/components/DataList/AllArtists.vue](/src/components/DataList/AllArtists.vue) | Vue | 54 | 0 | 5 | 59 |
| [src/components/DataList/PlayList.vue](/src/components/DataList/PlayList.vue) | Vue | 227 | 0 | 8 | 235 |
| [src/components/DataList/PlaylistGrid.vue](/src/components/DataList/PlaylistGrid.vue) | Vue | 144 | 0 | 13 | 157 |
| [src/components/DataModel/AddPlaylist.vue](/src/components/DataModel/AddPlaylist.vue) | Vue | 0 | 0 | 1 | 1 |
| [src/components/DataModel/PlaylistUpdate.vue](/src/components/DataModel/PlaylistUpdate.vue) | Vue | -1 | 0 | 0 | -1 |
| [src/components/Home/SystemStats.vue](/src/components/Home/SystemStats.vue) | Vue | 185 | 0 | 12 | 197 |
| [src/components/Nav/index.vue](/src/components/Nav/index.vue) | Vue | -10 | 0 | 4 | -6 |
| [src/components/Pagination/index.vue](/src/components/Pagination/index.vue) | Vue | 118 | 1 | 10 | 129 |
| [src/components/Player/BigPlayer.vue](/src/components/Player/BigPlayer.vue) | Vue | 695 | 3 | 23 | 721 |
| [src/components/Player/PlayerCover.vue](/src/components/Player/PlayerCover.vue) | Vue | 292 | 0 | 35 | 327 |
| [src/components/Player/PlayerRecord.vue](/src/components/Player/PlayerRecord.vue) | Vue | 82 | 0 | 11 | 93 |
| [src/components/Player/index.vue](/src/components/Player/index.vue) | Vue | 347 | -2 | 16 | 361 |
| [src/components/SearchInp/index.vue](/src/components/SearchInp/index.vue) | Vue | -501 | -1 | -16 | -518 |
| [src/router/index.ts](/src/router/index.ts) | TypeScript | -19 | 27 | 0 | 8 |
| [src/router/routes.ts](/src/router/routes.ts) | TypeScript | 99 | 10 | 1 | 110 |
| [src/store/musicData.ts](/src/store/musicData.ts) | TypeScript | 132 | 22 | 1 | 155 |
| [src/store/userData.ts](/src/store/userData.ts) | TypeScript | 2 | 0 | 0 | 2 |
| [src/utils/encrypt.ts](/src/utils/encrypt.ts) | TypeScript | 11 | 6 | 3 | 20 |
| [src/utils/lyricFormat.js](/src/utils/lyricFormat.js) | JavaScript | 36 | 4 | 8 | 48 |
| [src/utils/request.js](/src/utils/request.js) | JavaScript | -65 | -6 | -8 | -79 |
| [src/utils/request.ts](/src/utils/request.ts) | TypeScript | 143 | 27 | 21 | 191 |
| [src/utils/test\_lyric.js](/src/utils/test_lyric.js) | JavaScript | 37 | 0 | 8 | 45 |
| [src/utils/test\_lyric\_merge.js](/src/utils/test_lyric_merge.js) | JavaScript | 53 | 11 | 13 | 77 |
| [src/utils/timeTools.js](/src/utils/timeTools.js) | JavaScript | 3 | 0 | 1 | 4 |
| [src/views/Album/index.vue](/src/views/Album/index.vue) | Vue | 16 | 0 | 2 | 18 |
| [src/views/Artist/index.vue](/src/views/Artist/index.vue) | Vue | 16 | 0 | 2 | 18 |
| [src/views/History/HistoryView.vue](/src/views/History/HistoryView.vue) | Vue | 213 | 0 | 15 | 228 |
| [src/views/Home/index.vue](/src/views/Home/index.vue) | Vue | 99 | 0 | 10 | 109 |
| [src/views/Login/LoginView.vue](/src/views/Login/LoginView.vue) | Vue | 236 | 2 | 10 | 248 |
| [src/views/Playlist/PrivatePlaylists.vue](/src/views/Playlist/PrivatePlaylists.vue) | Vue | 62 | 0 | 7 | 69 |
| [src/views/Playlist/PublicPlaylists.vue](/src/views/Playlist/PublicPlaylists.vue) | Vue | 62 | 0 | 7 | 69 |
| [src/views/Playlist/index.vue](/src/views/Playlist/index.vue) | Vue | 337 | 0 | 34 | 371 |
| [src/views/Search/index.vue](/src/views/Search/index.vue) | Vue | 389 | 14 | 36 | 439 |
| [src/views/Setting/SettingView.vue](/src/views/Setting/SettingView.vue) | Vue | 279 | 0 | 13 | 292 |
| [src/views/Song/SongView.vue](/src/views/Song/SongView.vue) | Vue | 423 | 14 | 58 | 495 |
| [src/views/System/StatsView.vue](/src/views/System/StatsView.vue) | Vue | 30 | 0 | 5 | 35 |
| [src/views/User/index.vue](/src/views/User/index.vue) | Vue | 255 | 0 | 29 | 284 |
| [src/views/User/like.vue](/src/views/User/like.vue) | Vue | 19 | 0 | 2 | 21 |
| [vite.config.ts](/vite.config.ts) | TypeScript | 3 | 1 | 0 | 4 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details