package lib

import (
	"fmt"
	"errors"
)

func GetOpt(msg string) (*string, error){
	var x string 
	fmt.Printf("%s : ",msg)
	_, _ = fmt.Scanf("%s",&x)
	if len(x) == 0 {
		return nil, errors.New(fmt.Sprintf("Not %s Args.",msg))
	}

	return &x, nil
}


func PointerStr(s string) *string {
	str := s 
	return &str
}