package main

import (
	"fmt"

	"os"

	"github.com/hualet/moji"
)

func main() {
	c, err := moji.NewClient()
	if err != nil {
		fmt.Printf("error occured: %v", err)
		os.Exit(-1)
	}

	cond, err := c.ConditionByLatLong("30.6", "114.4")
	if err != nil {
		fmt.Printf("error occured: %v", err)
	}
	fmt.Printf("%v\n", cond)

	fore, err := c.ForecastByLatLong("30.6", "114.4")
	if err != nil {
		fmt.Printf("error occured: %v", err)
	}
	fmt.Printf("%v\n", fore)
}
