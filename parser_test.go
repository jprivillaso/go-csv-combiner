package main

import (
	"reflect"
	"testing"
)

func TestValidParser(t *testing.T) {
	ParseInputFlags()
	var usersMap = ParseCSVFiles()

	if !reflect.DeepEqual(usersMap["1"], UserInfo{
		Id:        "1",
		FirstName: "Juan",
		LastName:  "Rivillas",
		Phone:     "310-111-1111",
		Email:     "contact@lalala.com",
	}) {
		t.Fatalf(`User 1 mismatch`)
	}

	if !reflect.DeepEqual(usersMap["2"], UserInfo{
		Id:        "2",
		FirstName: "John",
		LastName:  "Dow",
		Phone:     "213-222-2222",
		Email:     "john@boo.com",
	}) {
		t.Fatalf(`User 2 mismatch`)
	}

	if !reflect.DeepEqual(usersMap["3"], UserInfo{
		Id:        "3",
		FirstName: "Marina",
		LastName:  "Davis",
		Phone:     "",
		Email:     "",
	}) {
		t.Fatalf(`User 3 mismatch`)
	}

	if !reflect.DeepEqual(usersMap["4"], UserInfo{
		Id:        "4",
		FirstName: "Michelle",
		LastName:  "Obama",
		Phone:     "213-444-444",
		Email:     "lion@yue.net",
	}) {
		t.Fatalf(`User 4 mismatch`)
	}
}
