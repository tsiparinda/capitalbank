package pb_models

import (
	"time"
)

type PbURL struct {
	URL       string
	Acc       string
	StartDate time.Time
	EndDate   time.Time
	FollowId  string
	Limit     int
}
