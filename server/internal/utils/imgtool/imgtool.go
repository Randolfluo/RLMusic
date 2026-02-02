package imgtool

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// ProcessAndSaveCover 处理并保存封面图片
// data: 图片原始二进制数据
// saveDir: 保存目录
// 返回: hash, filename, width, height, error
func ProcessAndSaveCover(data []byte, saveDir string) (string, string, int, int, error) {
	// 1. 计算 Hash
	hash := calculateHash(data)

	// 2. 解码图片
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", "", 0, 0, err
	}

	// 3. 调整大小 (150x150 thumbnail)
	// 使用 Lanczos3 算法进行高质量重采样
	// width: 150, height: 0 (保持纵横比) -> 或者 150x150 强制?
	// 用户要求: 存放 150x150px 的缩略图。通常是 Resize 到 150x150。
	thumbnail := resize.Resize(150, 150, img, resize.Lanczos3)

	// 4. 确定文件名
	ext := ".jpg"
	if format == "png" {
		ext = ".png"
	}
	filename := hash + ext
	fullPath := filepath.Join(saveDir, filename)

	// 5. 检查文件是否存在，不存在则保存
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// 确保目录存在
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return "", "", 0, 0, err
		}

		out, err := os.Create(fullPath)
		if err != nil {
			return "", "", 0, 0, err
		}
		defer out.Close()

		if format == "png" {
			err = png.Encode(out, thumbnail)
		} else {
			// 默认保存为 JPEG, 质量 80
			err = jpeg.Encode(out, thumbnail, &jpeg.Options{Quality: 80})
		}
		if err != nil {
			return "", "", 0, 0, err
		}
	}

	return hash, filename, 150, 150, nil
}

func calculateHash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
