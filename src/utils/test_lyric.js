const lyricFormat = (lrc, tlrc) => {
  const regex = /^\[([^\]]+)\]\s*(.*)$/;
  
  const parseLyric = (str) => {
    if(!str) return [];
    return str.split("\n")
      .map(line => {
        const match = line.match(regex);
        if (!match) return null;
        const [, timeStr, text] = match;
        const timeParts = timeStr.split(":");
        const min = parseInt(timeParts[0]);
        const sec = parseFloat(timeParts[1]);
        const time = min * 60 + sec;
        return { time, text: text.trim() };
      })
      .filter(item => item !== null && item.text !== "");
  };

  const lrcData = parseLyric(lrc);
  const tlrcData = parseLyric(tlrc);

  console.log("LRC Data:", lrcData);
  console.log("TLRC Data:", tlrcData);

  const result = lrcData.map(item => {
    const trans = tlrcData.find(t => Math.abs(t.time - item.time) < 0.5);
    return {
      time: item.time,
      lyric: item.text,
      lyricFy: trans ? trans.text : ""
    };
  });
  return result;
};

const lrc = `[00:01.00]Original Line 1
[00:05.50]Original Line 2`;

const tlrc = `[00:01.000]Translation Line 1
[00:05.600]Translation Line 2`;

const res = lyricFormat(lrc, tlrc);
console.log(JSON.stringify(res, null, 2));
