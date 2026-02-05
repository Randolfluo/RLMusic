/**
 * 格式化歌词字符串为时间轴和歌词文本的数组
 * @param {string} lrc 歌词字符串
 * @param {string} tlrc 翻译歌词字符串
 * @returns {Array<{ time: number, lyric: string, lyricFy: string }>} 时间轴和歌词文本的数组
 */
const lyricFormat = (lrc, tlrc) => {
  // 匹配时间轴和歌词文本的正则表达式
  const regex = /^\[(\d{2}):(\d{2})\.(\d{2,3})\](.*)$/;
  
  // 辅助函数：解析歌词字符串
  const parseLyric = (str) => {
    if(!str) return [];
    return str.split("\n")
      .map(line => {
        const match = line.match(regex);
        if (!match) return null;
        const [, minStr, secStr, msStr, text] = match;
        const min = parseInt(minStr, 10);
        const sec = parseInt(secStr, 10);
        // 处理2位或3位毫秒数
        const ms = parseFloat("0." + msStr);
        const time = min * 60 + sec + ms;
        return { time, text: text.trim() };
      })
      .filter(item => item !== null && item.text !== "");
  };

  const lrcData = parseLyric(lrc);
  const tlrcData = parseLyric(tlrc);
  const result = [];

  // 第1种模式：单文件中包含双语（同一时间戳出现两次）
  for (let i = 0; i < lrcData.length; i++) {
    const current = lrcData[i];
    
    // 如果当前行和下一行时间戳极度接近（<0.05s），且下一行也是文本
    if (i < lrcData.length - 1) {
      const next = lrcData[i+1];
      if (Math.abs(current.time - next.time) < 0.05) {
        // 判定为原文+译文的组合
        result.push({
          time: current.time,
          lyric: current.text,
          lyricFy: next.text
        });
        i++; // 跳过下一行（因为它是译文，已经合并了）
        continue;
      }
    }

    // 第2种模式：标准双文件（如果有 tlrcData）
    // 尝试从 tlrcData 查找翻译
    let transText = "";
    if (tlrcData.length > 0) {
       const trans = tlrcData.reduce((prev, curr) => {
            const diff = Math.abs(curr.time - current.time);
            if (diff > 1.0) return prev;
            if (!prev) return curr;
            return diff < Math.abs(prev.time - current.time) ? curr : prev;
        }, null);
        if (trans) transText = trans.text;
    }

    result.push({
      time: current.time,
      lyric: current.text,
      lyricFy: transText
    });
  }

  // 检查是否为纯音乐，是则返回空数组
  if (result.length && /纯音乐，请欣赏/.test(result[0].lyric)) {
    console.log("该歌曲为纯音乐");
    return [];
  }
  return result;
};



export default lyricFormat;
