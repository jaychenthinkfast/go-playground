package models

// RunResult 代码运行结果
type RunResult struct {
	Output    string `json:"output"`     // 运行输出
	Error     string `json:"error"`      // 错误信息
	ExitCode  int    `json:"exit_code"`  // 退出码
	Duration  int64  `json:"duration"`   // 运行时长（毫秒）
	Memory    int64  `json:"memory"`     // 内存使用（字节）
	CreatedAt int64  `json:"created_at"` // 创建时间戳
}
