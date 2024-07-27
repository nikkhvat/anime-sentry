package date

import (
	"fmt"
	"regexp"
	"strconv"
)

func ConvertDate(input string) (string, error) {
	monthMap := map[string]string{
		"января":   "01",
		"февраля":  "02",
		"марта":    "03",
		"апреля":   "04",
		"мая":      "05",
		"июня":     "06",
		"июля":     "07",
		"августа":  "08",
		"сентября": "09",
		"октября":  "10",
		"ноября":   "11",
		"декабря":  "12",
	}

	r := regexp.MustCompile(`(\d+)\s+(\p{Cyrillic}+)`)
	matches := r.FindStringSubmatch(input)
	if matches == nil || len(matches) < 3 {
		return "", fmt.Errorf("invalid date format")
	}

	day, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	month, ok := monthMap[matches[2]]
	if !ok {
		return "", fmt.Errorf("invalid month")
	}

	return fmt.Sprintf("%02d/%s", day, month), nil
}
