你是一名音乐分析师。
请分析以下歌曲数据并输出摘要：


必须：
1. 仅输出 JSON
2. 不要解释
3. 不要编造销量或奖项

输出字段：
- artist
- title
- era
- genre_guess
- vocal_style
- emotional_tone（数组）
- listening_scene（数组）
- energy_level（1-10）
- keywords（数组）

歌曲：
{{song_json}}
