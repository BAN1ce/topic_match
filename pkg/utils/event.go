package utils

func WithEventPrefix(prefix, s string) string {
	return prefix + "." + s
}
