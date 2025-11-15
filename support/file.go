package support

import "fmt"

// BytesToHumanReadable converts bytes to human-readable format (KB, MB, GB, etc.)
func BytesToHumanReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatSize Helper function to format file size in human-readable format
func FormatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	var suffix string
	var value float64

	switch {
	case size >= TB:
		suffix = "TB"
		value = float64(size) / TB
	case size >= GB:
		suffix = "GB"
		value = float64(size) / GB
	case size >= MB:
		suffix = "MB"
		value = float64(size) / MB
	case size >= KB:
		suffix = "KB"
		value = float64(size) / KB
	default:
		suffix = "B"
		value = float64(size)
	}

	if value == float64(int64(value)) {
		return fmt.Sprintf("%d%s", int64(value), suffix)
	}
	return fmt.Sprintf("%.2f%s", value, suffix)
}
