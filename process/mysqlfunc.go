package process

import (
	"fmt"
)

func mysqlLoadShares(query string, offset int) ([]Company, bool) {
	var sharePrices = make([]Company, 0)
	err := mysqlObj.Check()
	if err != nil {
		fmt.Println(err)
		err = mysqlObj.Reconnect()
		if err != nil {
			fmt.Println(err)
			return sharePrices, false
		}
	}
	fmt.Println(query, offset)

	stmt, err := mysqlObj.PrepareStatement(query)
	if err != nil {
		fmt.Println(err)
		return sharePrices, false
	}

	result, err := mysqlObj.SelectQuery(stmt, []any{offset})
	if err != nil {
		fmt.Println(err)
		return sharePrices, false
	}

	if len(result) != 0 {
		for _, values := range result {
			sharePrices = append(sharePrices, Company{
				Name:   values[0].(string),
				IShare: values[1].(float64),
				CShare: values[2].(float64),
			})
		}
	}

	return sharePrices, true
}
