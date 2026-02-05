const lyricFormat = (lrc, tlrc) => {
    // 匹配时间轴和歌词文本的正则表达式，修复正则表达式以匹配所有内容
    const regex = /^\[(\d{2}):(\d{2})\.(\d{2,3})\](.*)$/;

    // 辅助函数：解析歌词字符串
    const parseLyric = (str) => {
        if (!str) return [];
        return str.split("\n")
            .map(line => {
                const match = line.match(regex);
                if (!match) return null;
                const [, minStr, secStr, msStr, text] = match;
                const min = parseInt(minStr, 10);
                const sec = parseInt(secStr, 10);
                const ms = parseFloat("0." + msStr);
                const time = min * 60 + sec + ms;
                return { time, text: text.trim() };
            })
            .filter(item => item !== null && item.text !== "");
    };

    const lrcData = parseLyric(lrc);
    const tlrcData = parseLyric(tlrc);

    // console.log("LRC Data:", lrcData);

    const result = [];
    
    // 这种格式：同一时间戳的两行，第一行是原文，第二行是译文
    // [00:16.82]いたいけなモーション
    // [00:16.82]讨人怜爱的动作
    
    for (let i = 0; i < lrcData.length; i++) {
        const current = lrcData[i];
        
        // 如果当前行和下一行时间戳完全相同（或者极度接近）
        // 并且这一行还没有被处理过
        if (i < lrcData.length - 1) {
            const next = lrcData[i+1];
            // 判断时间差是否小于 0.05s
            if (Math.abs(current.time - next.time) < 0.05) {
                // 将这两行合并为一行：current是原文，next是译文
                result.push({
                    time: current.time,
                    lyric: current.text,
                    lyricFy: next.text
                });
                i++; // 跳过下一行
                continue;
            }
        }
        
        // 正常情况：尝试从 tlrcData 查找翻译（兼容标准双文件格式）
        const trans = tlrcData.find(t => Math.abs(t.time - current.time) < 1.0);
        result.push({
             time: current.time,
             lyric: current.text,
             lyricFy: trans ? trans.text : ""
        });
    }

    return result;
};

const lrc = `[00:00.00]作词 : とみー
[00:01.00]作曲 : とみー
[00:16.82]いたいけなモーション
[00:16.82]讨人怜爱的动作
[00:18.38]振り切れるテンション
[00:18.38]超越计量的紧张
[00:19.92]意外、意外`;

const tlrc = ""; // tlyric 为空

const res = lyricFormat(lrc, tlrc);
console.log(JSON.stringify(res, null, 2));
