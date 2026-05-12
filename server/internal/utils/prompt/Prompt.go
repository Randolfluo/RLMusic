package prompt

import (
	"os"
	"path/filepath"
)

// Read 按照「可执行文件目录优先、CWD 兜底」的顺序读取 prompts 目录下指定文件。
// 这样无论二进制是从 server/ 还是 server/bin/ 启动，都能找到 prompts 目录。
//
// 查找顺序:
//  1. <exeDir>/prompts/<name>          —— 部署形态：与可执行文件同级
//  2. <exeDir>/../prompts/<name>       —— 部署形态：可执行文件位于 bin/ 子目录
//  3. ./prompts/<name>                 —— 开发形态：go run 时 CWD 为 server/
func Read(name string) ([]byte, error) {
	candidates := make([]string, 0, 3)

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "prompts", name),
			filepath.Join(exeDir, "..", "prompts", name),
		)
	}
	candidates = append(candidates, filepath.Join("prompts", name))

	var lastErr error
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
