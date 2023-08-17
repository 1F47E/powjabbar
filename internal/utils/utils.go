package utils

// convert difficulty to criteria string
func DifficultyToCriteria(difficulty int) string {
	criteria := make([]rune, difficulty)
	for i := 0; i < difficulty; i++ {
		criteria[i] = '0'
	}
	return string(criteria)
}
