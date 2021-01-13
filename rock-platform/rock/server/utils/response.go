package utils

//import (
//	"encoding/json"
//	"fmt"
//)
//
//func MarshalResponse(src, dest interface{}) error {
//	respByte, err := json.Marshal(src)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("respByte: %v\n", string(respByte))
//
//	//dest := new(interface{})
//	err = json.Unmarshal(respByte, dest) // map[]
//	if err != nil {
//		return err
//	}
//	fmt.Printf("dest: %v\n", dest)
//	return nil
//}
