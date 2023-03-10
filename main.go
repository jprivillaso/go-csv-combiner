package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"
)

type UserInfo struct {
	Id        string
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

// This structure can be a subset of the struct. The idea is to give a bit of control over
// what you really want from the data
var lookupFields = []string{"Id", "FirstName", "LastName", "Phone", "Email"}

var flagDir = flag.String("dir", "", "directory containing CSVs")

// Sets a property dynamically to the struct
func (i *UserInfo) setProperty(propName string, propValue string) *UserInfo {
	reflect.ValueOf(i).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))
	return i
}

func main() {
	fmt.Println("Process Started ...")
	var start = time.Now().UnixNano()
	ParseInputFlags()
	WriteOutput(ParseCSVFiles())
	var end = time.Now().UnixNano()
	fmt.Println(fmt.Sprintf("Finished! %d ns", end-start))
}
