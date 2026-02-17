# Database Models Schema

This document describes the database schema for the LocalMusicPlayer server, based on the GORM models.

## Tables

### User
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | User ID | Primary Key, Auto Increment |
| `created_at` | datetime | Creation Time | |
| `updated_at` | datetime | Update Time | |
| `username` | varchar(50) | Username | Unique |
| `password` | varchar(100) | Encrypted Password | |
| `email` | varchar(255) | Email Address | |
| `user_group` | varchar(50) | User Role | Default: 'user' |
| `avatar` | varchar(255) | Avatar URL | |
| `listening_duration` | int64 | Total Listening Duration (ms) | Default: 0 |
| `total_duration` | int64 | Total Songs Duration (ms) | Default: 0 |
| `last_login` | datetime | Last Login Time | |
| `is_delete` | bool | Soft Delete Flag | |
| `ip_src` | varchar(50) | Last Login IP | |

**Relations:**
- `SubscribedPlaylists` (Many-to-Many): Users can subscribe to multiple Playlists (`user_subscribed_playlists` table).

---

### Song
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | Song ID | Primary Key, Auto Increment |
| `title` | varchar(255) | Song Title | Index |
| `artist_name` | varchar(255) | Artist Name (Denormalized) | Index |
| `album_name` | varchar(255) | Album Name (Denormalized) | Index |
| `artist_id` | int | Artist ID | Index |
| `album_id` | int | Album ID | Index |
| `cover_id` | int | Cover ID | Index |
| `track_num` | int | Track Number | |
| `disc_num` | int | Disc Number | |
| `year` | varchar(20) | Release Year | |
| `file_path` | varchar(500) | File Path | Not Null |
| `file_name` | varchar(255) | File Name | |
| `file_size` | int64 | File Size (bytes) | |
| `format` | varchar(20) | Audio Format | Default: 'flac' |
| `duration` | float64 | Duration (seconds) | |
| `sample_rate` | int | Sample Rate (Hz) | |
| `bit_depth` | int | Bit Depth | |
| `channels` | int | Channels | |
| `bit_rate` | int | Bit Rate (kbps) | |
| `play_count` | int | Play Count | Default: 0 |
| `is_delete` | bool | Soft Delete Flag | Default: false |

**Relations:**
- `Artist` (Belongs To): A song belongs to an Artist.
- `Artists` (Many-to-Many): A song can have multiple Artists (`song_artists` table).
- `Album` (Belongs To): A song belongs to an Album.
- `Cover` (Belongs To): A song has a Cover.

---

### Playlist
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | Playlist ID | Primary Key, Auto Increment |
| `title` | varchar(255) | Playlist Title | Not Null, Index |
| `description` | text | Description | |
| `is_public` | bool | Is Public | Index |
| `cover_url` | varchar(500) | Cover URL | |
| `owner_id` | int | Owner User ID | Index |
| `play_count` | int | Play Count | Default: 0 |
| `total_songs` | int | Total Songs Count | Default: 0 |

**Relations:**
- `Songs` (Many-to-Many): A playlist contains multiple Songs (`playlist_songs` table).

---

### History
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | uint | History ID | Primary Key |
| `user_id` | int | User ID | Index |
| `song_id` | int | Song ID | Index |
| `created_at` | datetime | Play Time | |

**Relations:**
- `Song` (Belongs To): A history record links to a Song.

---

### Artist
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | Artist ID | Primary Key, Auto Increment |
| `name` | varchar(255) | Artist Name | Not Null, Index |
| `description` | text | Description | |
| `cover` | varchar(500) | Cover URL | |

**Relations:**
- `Songs` (Many-to-Many): An artist has multiple Songs (`song_artists` table).

---

### Album
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | Album ID | Primary Key, Auto Increment |
| `title` | varchar(255) | Album Title | Not Null, Index |
| `description` | text | Description | |
| `release_date` | datetime | Release Date | |
| `cover` | varchar(500) | Cover URL | |
| `artist_id` | int | Artist ID | Index |

**Relations:**
- `Artist` (Belongs To): An album belongs to an Artist.
- `Songs` (Has Many): An album has multiple Songs.

---

### Cover
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | int | Cover ID | Primary Key, Auto Increment |
| `hash` | varchar(64) | Image Hash | Unique Index |
| `path` | varchar(500) | Image Path | |
| `format` | varchar(10) | Image Format | |
| `size` | int64 | File Size | |
| `width` | int | Width | |
| `height` | int | Height | |

---

### SystemInfo
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `key` | varchar(100) | Config Key | Primary Key |
| `value` | text | Config Value | |

---

### Podcast
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | uint | Podcast ID | Primary Key |
| `title` | string | Podcast Title | Index |
| `cover` | string | Cover URL | |
| `is_public` | bool | Is Public | Index |

**Relations:**
- `Episodes` (Has Many): A podcast has multiple Episodes.

---

### PodcastEpisode
| Column | Type | Description | Constraints |
| :--- | :--- | :--- | :--- |
| `id` | uint | Episode ID | Primary Key |
| `podcast_id` | uint | Podcast ID | Index, Not Null |
| `song_id` | string | Song ID | Index, Not Null, Size: 36 |
| `intro_audio` | string | Intro Audio Path | |

**Relations:**
- `Song` (Belongs To): An episode is linked to a Song.
