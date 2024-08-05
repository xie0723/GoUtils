package GoUtils

import "strings"

// ConcatString 高性能字符串拼接
func ConcatString(sl ...string) string {
	n := 0
	for i := 0; i < len(sl); i++ {
		n += len(sl[i])
	}

	var b strings.Builder
	b.Grow(n)
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}