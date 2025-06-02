package utils

func Contains(s, substr string) (bool, int) {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true, i
		}
	}
	return false, -1
}
