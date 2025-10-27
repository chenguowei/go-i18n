package i18n

import (
	"fmt"
	"runtime"
)

const (
	// Version 版本号
	Version = "1.0.0"
	// BuildTime 构建时间
	BuildTime = "2024-01-01T00:00:00Z"
	// GitCommit Git提交号
	GitCommit = "unknown"
)

// BuildInfo 构建信息
type BuildInfo struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// GetBuildInfo 获取构建信息
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		GoVersion: runtime.Version(),
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String 返回版本信息字符串
func (b BuildInfo) String() string {
	return fmt.Sprintf("GoI18n-Gin v%s (Go: %s, OS/Arch: %s/%s)",
		b.Version, b.GoVersion, b.OS, b.Arch)
}

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}

// IsCompatible 检查版本兼容性
func IsCompatible(minVersion string) bool {
	// 简单的版本比较，实际项目中可以使用更复杂的版本比较库
	return Version >= minVersion
}