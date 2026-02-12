package model

// Podcast 播客/节目模型
type Podcast struct {
	ID       uint             `gorm:"primarykey"`
	Title    string           `json:"title" gorm:"index;comment:播客标题"`
	Cover    string           `json:"cover" gorm:"comment:播客封面URL"` // 类似专辑封面
	IsPublic bool             `gorm:"index" json:"is_public"`
	Episodes []PodcastEpisode `json:"episodes" gorm:"foreignKey:PodcastID"`
}

// PodcastEpisode 播客单集关联表
type PodcastEpisode struct {
	ID        uint   `gorm:"primarykey"`
	PodcastID uint   `json:"podcast_id" gorm:"index;not null;comment:所属播客ID"`
	SongID    string `json:"song_id" gorm:"index;not null;size:36;comment:关联的音频文件ID"`

	IntroAudio string `json:"intro_audio" gorm:"comment:开场白音频文件路径"`

	Song Song `json:"song" gorm:"foreignKey:SongID"`
}
