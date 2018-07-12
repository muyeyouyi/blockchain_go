package utils

import "fmt"

func LogE(err error)  {
	if	err != nil{
		fmt.Println(err)
	}
}