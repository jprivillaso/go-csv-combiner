package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
)

// Finds an element in an array and return its index. If the element is not present
// it returns -1
func findIndex(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}

	return -1
}

// This function will panic if the content is not a string
func GetField(v *UserInfo, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

// Precompute indexes for O(1) lookups in the Map function
func GetHeaderIndexes(header []string) map[string]int {
	headerIndexes := make(map[string]int)

	for key := range lookupFields {
		var propName = lookupFields[key]
		var propIndex = findIndex(header, propName)
		headerIndexes[propName] = propIndex
	}

	return headerIndexes
}

func GetUniqUserIds(usersMap map[string]UserInfo) []int {
	var userIds []int

	for k := range usersMap {
		var userId, _ = strconv.Atoi(k)
		userIds = append(userIds, userId)
	}

	sort.Ints(userIds)

	return userIds
}

func ParseInputFlags() {
	flag.Parse()

	if *flagDir == "" {
		fmt.Println("missing -dir flag")
		os.Exit(1)
	}
}

// Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

func Transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	decodeError := json.NewEncoder(buf).Encode(in)

	if decodeError != nil {
		panic(decodeError)
	}

	encodeError := json.NewDecoder(buf).Decode(out)

	if encodeError != nil {
		panic(encodeError)
	}
}
