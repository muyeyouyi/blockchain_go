package utils

import "fmt"

func LogE(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func LogD(log string) {
	if log != "" {
		fmt.Println(log)
	}
}
