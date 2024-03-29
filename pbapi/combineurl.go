package pbapi

import (
	pb_models "capitalbank/pbapi/models"
	"fmt"
	neturl "net/url"
	"strconv"
)



func (a PrivatBankAPI) CombineURL(us pb_models.PbURL) (url string, err error) {
	err = nil

	values := neturl.Values{}

	url = us.URL

	if us.Acc != "" {
		values.Add("acc", us.Acc)
	} else {
		return url, fmt.Errorf("Error raised while was combined URL: the Account is a required parameter!")
	}

	if !us.StartDate.IsZero() {
		values.Add("startDate", us.StartDate.Format("02-01-2006"))
	} else {
		return url, fmt.Errorf("Error raised while was combined URL: the StartDate is a required parameter!")
	}

	if !us.EndDate.IsZero() {
		values.Add("endDate", us.EndDate.Format("02-01-2006"))
		//url = url + "&endDate=" + us.EndDate.Format("02-01-2006")
	}

	if us.Limit > 0 {
		values.Add("limit", strconv.Itoa(us.Limit))
	}

	if us.FollowId != "" {
		values.Add("followId", us.FollowId)
	}

	if len(values) > 0 {
		url = url + "?" + values.Encode()
	}

	return
}
