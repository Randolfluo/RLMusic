# RLMusic 后端 Model 数据库字典

本文档根据 server/internal/model 下的 Go 模型导出。

说明:
- 主键: 字段是否为主键。
- 非空: 字段是否被显式约束为非空，或由主键隐式非空。
- 表名遵循 GORM 默认命名规则。

## 表: users

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | User ID | Yes | Yes |
| created_at | Create time | No | No |
| updated_at | Update time | No | No |
| username | Username (unique) | No | No |
| password | Password hash | No | No |
| email | Email | No | No |
| user_group | User group | No | No |
| avatar | Avatar path | No | No |
| listening_duration | Total listened duration | No | No |
| total_duration | Total song duration | No | No |
| last_login | Last login time | No | No |
| is_delete | Soft delete flag | No | No |
| ip_src | Last login IP source | No | No |

## 表: system_infos

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| key | System metric key | Yes | Yes |
| value | System metric value | No | No |

## 表: artists

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | Artist ID | Yes | Yes |
| name | Artist name | No | Yes |
| description | Artist description | No | No |
| cover | Cover URL | No | No |

## 表: albums

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | Album ID | Yes | Yes |
| title | Album title | No | Yes |
| description | Album description | No | No |
| release_date | Release date | No | No |
| cover | Cover URL | No | No |
| artist_id | Related artist ID (nullable FK) | No | No |

## 表: covers

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | Cover ID | Yes | Yes |
| hash | Content hash (unique) | No | No |
| path | File path | No | No |
| format | Image format | No | No |
| size | File size | No | No |
| width | Image width | No | No |
| height | Image height | No | No |

## 表: songs

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | Song ID | Yes | Yes |
| title | Song title | No | No |
| artist_name | Denormalized artist name | No | No |
| album_name | Denormalized album name | No | No |
| artist_id | Related artist ID (nullable FK) | No | No |
| album_id | Related album ID (nullable FK) | No | No |
| cover_id | Related cover ID (nullable FK) | No | No |
| track_num | Track number | No | No |
| disc_num | Disc number | No | No |
| year | Year | No | No |
| file_path | File path | No | Yes |
| file_name | File name | No | No |
| file_size | File size | No | No |
| format | Audio format | No | No |
| duration | Duration in seconds | No | No |
| sample_rate | Sample rate | No | No |
| bit_depth | Bit depth | No | No |
| channels | Channels | No | No |
| bit_rate | Bit rate | No | No |
| play_count | Play count | No | No |
| description | AI description/opening text | No | No |
| opening_audio_file | AI opening audio file path | No | No |
| is_delete | Soft delete flag | No | No |

## 表: playlists

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | Playlist ID | Yes | Yes |
| title | Playlist title | No | Yes |
| description | Playlist description | No | No |
| has_intro | Has opening intro | No | No |
| is_public | Public visibility | No | No |
| cover_url | Playlist cover URL | No | No |
| owner_id | Owner user ID | No | No |
| play_count | Play count | No | No |
| total_songs | Total songs | No | No |

## 表: histories

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| id | History ID | Yes | Yes |
| user_id | User ID | No | No |
| song_id | Song ID | No | No |
| created_at | Play time | No | No |

## 关联表: song_artists

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| song_id | Song ID | Yes (composite) | Yes |
| artist_id | Artist ID | Yes (composite) | Yes |

## 关联表: playlist_songs

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| playlist_id | Playlist ID | Yes (composite) | Yes |
| song_id | Song ID | Yes (composite) | Yes |

## 关联表: user_subscribed_playlists

| 字段名称 | 字段说明 | 主键 | 非空 |
| --- | --- | --- | --- |
| user_id | User ID | Yes (composite) | Yes |
| playlist_id | Playlist ID | Yes (composite) | Yes |

## 来源模型

- `server/internal/model/User.go`
- `server/internal/model/SystemInfo.go`
- `server/internal/model/Metadata.go`
- `server/internal/model/Song.go`
- `server/internal/model/Playlist.go`
- `server/internal/model/History.go`