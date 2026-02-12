# 播客 (Podcast) 数据库模型设计

基于现有的歌单 (Playlist) 和歌曲 (Song) 模型，我们设计一个适用于播客场景的数据库模型。

## 核心概念

1.  **Podcast (播客/节目)**: 类似于 `Playlist`，代表一个播客专辑或系列。
2.  **PodcastEpisode (播客单集)**: 关联 `Song` (实际音频文件) 和 `Podcast`，并增加额外的元数据（如开场白）。

## 数据库表设计

### 1. Podcast (播客系列表)

该表存储播客的元数据，结构与 `playlists` 表非常相似。

```go
type Podcast struct {
    g.GVA_MODEL
    Title       string `json:"title" gorm:"index;comment:播客标题"`
    Description string `json:"description" gorm:"comment:播客描述"`
    Cover       string `json:"cover" gorm:"comment:播客封面URL"` // 类似专辑封面
    Author      string `json:"author" gorm:"index;comment:主播/作者"`
    UserID      uint   `json:"user_id" gorm:"index;comment:创建者ID"`
    IsPublic    bool   `json:"is_public" gorm:"default:false;comment:是否公开"`
    
    // 关联
    User     User      `json:"user" gorm:"foreignKey:UserID"`
    Episodes []PodcastEpisode `json:"episodes" gorm:"foreignKey:PodcastID"`
}
```

### 2. PodcastEpisode (播客单集关联表)

这是核心差异所在。在普通歌单中，我们可能只需要 `playlist_songs` 中间表来关联。但在播客中，我们需要为每一集指定一个**开场白 (Intro)**。

我们不直接修改 `Song` 表，因为同一个音频文件（Song）可能被用在不同的播客中，或者作为纯音乐播放时不需要开场白。我们将开场白信息放在关联表中。

```go
type PodcastEpisode struct {
    g.GVA_MODEL
    PodcastID uint `json:"podcast_id" gorm:"index;not null;comment:所属播客ID"`
    SongID    string `json:"song_id" gorm:"index;not null;size:36;comment:关联的音频文件ID"` // 假设 Song ID 是 UUID 字符串
    
    // 排序
    SortOrder int `json:"sort_order" gorm:"default:0;comment:播放顺序"`
    
    // 特有字段：开场白
    // 这里存储的是开场白音频文件的静态资源路径 (URL 或 相对路径)
    // 例如: "/podcast/intro_12345.flac"
    IntroAudio string `json:"intro_audio" gorm:"comment:开场白音频文件路径"`
    
    // 可选：单集特定的标题或描述（如果不同于 Song 本身的信息）
    Title       string `json:"title" gorm:"comment:单集标题(覆盖Song标题)"`
    Description string `json:"description" gorm:"comment:单集简介"`

    // 关联
    Podcast Podcast `json:"podcast" gorm:"foreignKey:PodcastID"`
    Song    Song    `json:"song" gorm:"foreignKey:SongID"`
}
```

## 业务逻辑说明

1.  **创建播客**: 用户创建一个 `Podcast` 条目。
2.  **添加单集**: 
    *   用户选择一个现有的音频文件 (`Song`)。
    *   系统调用 TTS 接口生成开场白音频，保存到静态资源目录 (`/podcast`)。
    *   创建 `PodcastEpisode` 记录，填入 `PodcastID`, `SongID` 以及生成的 `IntroAudio` 路径。
3.  **播放逻辑**:
    *   前端获取播客详情，得到 `Episodes` 列表。
    *   遍历播放时，对于每个 Episode：
        1.  先加载并播放 `IntroAudio` (如果存在)。
        2.  Intro 播放完毕后，无缝衔接播放 `Song` 对应的音频文件。

## API 接口设计 (示例)

*   `POST /api/podcast`: 创建播客
*   `POST /api/podcast/:id/episode`: 添加单集
    *   Request: `{ "song_id": "...", "intro_text": "本期节目我们将欣赏..." }`
    *   Backend: 生成 TTS -> 保存文件 -> 创建 PodcastEpisode
*   `GET /api/podcast/:id`: 获取播客详情 (包含 Episodes 列表)

## 总结

通过引入 `PodcastEpisode` 中间实体，我们灵活地实现了“每首歌曲前添加开场白”的需求，而无需污染通用的 `Song` 模型，同时也复用了现有的音频文件管理能力。
