import axios from "@/utils/request";
import { getSongCover, resolveCoverUrl } from "@/api/song";

/**
 * 搜索部分 API (Updated)
 */

/**
 * 获取热门搜索列表
 */
export const getSearchHot = () => {
  return Promise.resolve({ data: [] });
};

export const getSearchSongs = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/song", { params: { keywords, limit, offset } });
};

export const getSearchArtists = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/artist", { params: { keywords, limit, offset } });
};

export const getSearchAlbums = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/album", { params: { keywords, limit, offset } });
};

export const getSearchPlaylists = (keywords: string, limit: number = 30, offset: number = 1) => {
    return axios.get("/search/playlist", { params: { keywords, limit, offset } });
};

/**
 * 聚合搜索建议 (For Search/index.vue)
 */
export const getSearchSuggest = async (keywords: string) => {
    // Parallel requests
    const [songs, artists, albums, playlists] = await Promise.all([
        getSearchSongs(keywords, 10),
        getSearchArtists(keywords, 5),
        getSearchAlbums(keywords, 5),
        getSearchPlaylists(keywords, 5)
    ]);
    
    // Extract data from backend response structure: { code: 1000, data: { result: { songs: [] } } }
    // Note: Request interceptor returns response.data, so code/msg/data are available directly.
    
    // Process songs to match frontend component expectations (e.g. ID -> id)
    const rawSongs = (songs as any).data?.result?.songs || [];
    const processedSongs = rawSongs.map((song: any) => ({
        ...song,
        id: song.ID || song.id,
        album_title: song.album_name || song.album_title,
        cover_url: song.cover_url || (song.id ? getSongCover(song.id) : undefined),
    }));

    // Process artists: backend returns 'cover', frontend expects 'picUrl'
    const rawArtists = (artists as any).data?.result?.artists || [];
    const processedArtists = rawArtists.map((artist: any) => ({
        ...artist,
        id: artist.ID || artist.id,
        picUrl: artist.picUrl || artist.img1v1Url || resolveCoverUrl(artist.cover),
    }));

    // Process albums: backend returns 'cover' and 'title', frontend expects 'picUrl' and 'name'
    const rawAlbums = (albums as any).data?.result?.albums || [];
    const processedAlbums = rawAlbums.map((album: any) => ({
        ...album,
        id: album.ID || album.id,
        name: album.name || album.title,
        picUrl: album.picUrl || resolveCoverUrl(album.cover),
    }));

    // Process playlists: backend returns 'cover_url' and 'title', frontend expects 'coverImgUrl' and 'name'
    const rawPlaylists = (playlists as any).data?.result?.playlists || [];
    const processedPlaylists = rawPlaylists.map((playlist: any) => ({
        ...playlist,
        id: playlist.ID || playlist.id,
        name: playlist.name || playlist.title,
        coverImgUrl: playlist.coverImgUrl || playlist.picUrl || resolveCoverUrl(playlist.cover_url),
    }));

    return {
        songs: processedSongs,
        artists: processedArtists,
        albums: processedAlbums,
        playlists: processedPlaylists
    };
}

