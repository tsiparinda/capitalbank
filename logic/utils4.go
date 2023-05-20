package logic

import (
	"capitalbank/store"
	"time"
)

func EarliestFBADTrue(balances *[]store.DataCheckBalance) (time.Time, error) {
	//var earliest *store.DataCheckBalance
	 earliestDate := time.Now().AddDate(0, 0, 2)

	for _, balance := range *balances {
		if balance.FBAD {
			if balance.Dpd.Before(earliestDate) {
				//earliest = &balance
				earliestDate = balance.Dpd
			}
		}
	}
	return earliestDate, nil
}

// func main() {
// 	balances := []DataCheckBalance{
// 		{"2023-05-18", "bank1", true},
// 		{"2023-05-19", "bank2", false},
// 		{"2023-05-20", "bank3", true},
// 		{"2023-05-17", "bank4", true},
// 		{"2023-05-21", "bank5", false},
// 	}

// 	earliest, err := EarliestFBADTrue(balances)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}

// 	if earliest != nil {
// 		fmt.Printf("The earliest date with FBAD=true is %v for %v\n", earliest.Dpd, earliest.BankRegistr)
// 	} else {
// 		fmt.Println("No dates with FBAD=true found")
// 	}
// }
