package utilities

import (
	"fmt"
	"time"
)

func GetDate() {
	curr_date := time.Date(2024, 2, 29, 0, 0, 0, 0, time.Now().Location())
	//oneDayLater := curr_date.AddDate(0, 0, 1)
	oneMonthLater := curr_date.AddDate(0, 1, 0)
	threeMonthLater := curr_date.AddDate(0, 3, 0)
	halfyearly := curr_date.AddDate(0, 6, 0)
	oneYearLater := curr_date.AddDate(1, 0, 0)

	oneYearBack := curr_date.AddDate(-1, 0, 0)

	fmt.Println("Current date: ", curr_date)

	//fmt.Println("oneDayLater: ", oneDayLater)
	fmt.Println("Three Months Later", threeMonthLater)
	fmt.Println("Half yearly", halfyearly)
	fmt.Println("oneMonthLater: ", oneMonthLater)

	fmt.Println("oneYearLater: ", oneYearLater)

	fmt.Println("oneYearBack: ", oneYearBack)
}

func CalculateAgeold(idate1, idate2 string) (oage uint) {

	const (
		YYYYMMDD = "20200101"
	)
	return
}
