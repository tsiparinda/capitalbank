package utils

import (
	"fmt"
	"time"
)

func GetShortDate(t time.Time) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", "2023-04-20 12:17:41.0940092 +0300 EEST")
	if err != nil {
		fmt.Println(err)
		return t, err
	}

	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t, nil
}
