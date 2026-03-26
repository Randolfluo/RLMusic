# 数据库表结构详细设计

> 本文档详细列出了 RLMusic 系统中所有数据库表的设计细节，包括表名、字段名、数据类型、主键、非空约束及字段含义。

## 4.1 用户表 (User)

**表名**: `user`
**说明**: 存储系统注册用户的账号、密码、权限及偏好信息。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 用户唯一标识 | int | 是 | 是 |
| `created_at` | 创建时间 | datetime | 否 | 否 |
| `updated_at` | 更新时间 | datetime | 否 | 否 |
| `username` | 登录用户名 | varchar(50) | 否 | 是 |
| `password` | 登录密码 | varchar(100) | 否 | 是 |
| `email` | 用户邮箱 | varchar(255) | 否 | 否 |
| `user_group` | 用户组/权限 | varchar(50) | 否 | 是 |
| `avatar` | 头像路径 | varchar(255) | 否 | 否 |
| `listening_duration` | 累计听歌时长 | int64 | 否 | 是 |
| `total_duration` | 曲库总时长 | int64 | 否 | 是 |
| `last_login` | 最后登录时间 | datetime | 否 | 否 |
| `is_delete` | 软删除标记 | bool | 否 | 是 |
| `ip_src` | 登录 IP | varchar(50) | 否 | 否 |

## 4.2 歌曲表 (Song)

**表名**: `song`
**说明**: 存储本地扫描到的音乐文件元数据及文件路径信息。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 歌曲唯一标识 | int | 是 | 是 |
| `title` | 歌曲标题 | varchar(255) | 否 | 否 |
| `artist_name` | 歌手名称 | varchar(255) | 否 | 否 |
| `album_name` | 专辑名称 | varchar(255) | 否 | 否 |
| `artist_id` | 关联歌手ID | int | 否 | 否 |
| `album_id` | 关联专辑ID | int | 否 | 否 |
| `cover_id` | 关联封面ID | int | 否 | 否 |
| `track_num` | 轨道号 | int | 否 | 否 |
| `disc_num` | 碟号 | int | 否 | 否 |
| `year` | 年份 | varchar(20) | 否 | 否 |
| `file_path` | 文件存储路径 | varchar(500) | 否 | 是 |
| `file_name` | 文件名 | varchar(255) | 否 | 否 |
| `file_size` | 文件大小 | int64 | 否 | 否 |
| `format` | 音频格式 | varchar(20) | 否 | 是 |
| `duration` | 时长 | float64 | 否 | 否 |
| `sample_rate` | 采样率 | int | 否 | 否 |
| `bit_depth` | 位深 | int | 否 | 否 |
| `channels` | 声道数 | int | 否 | 否 |
| `bit_rate` | 比特率 | int | 否 | 否 |
| `play_count` | 播放次数 | int | 否 | 是 |
| `description` | AI描述/开场白 | text | 否 | 否 |
| `opening_audio_file` | AI开场白音频 | varchar(500) | 否 | 否 |
| `is_delete` | 软删除标记 | bool | 否 | 是 |

## 4.3 歌单表 (Playlist)

**表名**: `playlist`
**说明**: 存储用户创建的歌单或系统生成的歌单信息。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 歌单唯一标识 | int | 是 | 是 |
| `title` | 歌单标题 | varchar(255) | 否 | 是 |
| `description` | 歌单描述 | text | 否 | 否 |
| `has_intro` | 是否有开场白 | bool | 否 | 是 |
| `is_public` | 是否公开 | bool | 否 | 否 |
| `cover_url` | 歌单封面 | varchar(500) | 否 | 否 |
| `owner_id` | 创建者ID | int | 否 | 否 |
| `play_count` | 播放次数 | int | 否 | 是 |
| `total_songs` | 歌曲总数 | int | 否 | 是 |

## 4.4 艺术家表 (Artist)

**表名**: `artist`
**说明**: 存储从音频标签中解析出的歌手信息。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 歌手唯一标识 | int | 是 | 是 |
| `name` | 歌手名称 | varchar(255) | 否 | 是 |
| `description` | 歌手简介 | text | 否 | 否 |
| `cover` | 歌手封面 | varchar(500) | 否 | 否 |

## 4.5 专辑表 (Album)

**表名**: `album`
**说明**: 存储从音频标签中解析出的专辑信息。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 专辑唯一标识 | int | 是 | 是 |
| `title` | 专辑标题 | varchar(255) | 否 | 是 |
| `description` | 专辑简介 | text | 否 | 否 |
| `release_date` | 发行日期 | datetime | 否 | 否 |
| `cover` | 专辑封面 | varchar(500) | 否 | 否 |
| `artist_id` | 关联歌手ID | int | 否 | 否 |

## 4.6 封面表 (Cover)

**表名**: `cover`
**说明**: 存储提取出的封面图片信息，用于去重和管理。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 封面唯一标识 | int | 是 | 是 |
| `hash` | 图片哈希 | varchar(64) | 否 | 否 |
| `path` | 存储路径 | varchar(500) | 否 | 否 |
| `format` | 图片格式 | varchar(10) | 否 | 否 |
| `size` | 文件大小 | int64 | 否 | 否 |
| `width` | 宽度 | int | 否 | 否 |
| `height` | 高度 | int | 否 | 否 |

## 4.7 播放历史表 (History)

**表名**: `history`
**说明**: 记录用户的歌曲播放历史。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | 记录唯一标识 | uint | 是 | 是 |
| `user_id` | 用户ID | int | 否 | 否 |
| `song_id` | 歌曲ID | int | 否 | 否 |
| `created_at` | 播放时间 | datetime | 否 | 否 |

## 4.8 系统信息表 (SystemInfo)

**表名**: `system_info`
**说明**: 存储系统级的统计信息（Key-Value 结构）。

| 字段名称 | 字段含义 | 数据类型 | 主键 | 非空 |
| :--- | :--- | :--- | :--- | :--- |
| `key` | 统计项键名 | varchar(100) | 是 | 是 |
| `value` | 统计项值 | text | 否 | 否 |

## 4.9 中间关联表

以下表为 Gorm 自动维护的多对多关联表，无独立模型文件，但在数据库中实际存在。

### 4.9.1 歌单-歌曲关联表 (playlist_songs)
| 字段名称 | 字段含义 | 数据类型 | 主键 |
| :--- | :--- | :--- | :--- |
| `playlist_id` | 歌单ID | int | 是 |
| `song_id` | 歌曲ID | int | 是 |

### 4.9.2 歌曲-歌手关联表 (song_artists)
| 字段名称 | 字段含义 | 数据类型 | 主键 |
| :--- | :--- | :--- | :--- |
| `song_id` | 歌曲ID | int | 是 |
| `artist_id` | 歌手ID | int | 是 |

### 4.9.3 用户-收藏歌单关联表 (user_subscribed_playlists)
| 字段名称 | 字段含义 | 数据类型 | 主键 |
| :--- | :--- | :--- | :--- |
| `user_id` | 用户ID | int | 是 |
| `playlist_id` | 歌单ID | int | 是 |
