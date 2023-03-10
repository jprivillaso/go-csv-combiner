package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// This function transforms the CSV lines into UserInfo structs that the Reducer can parse easily
func Mapper(user []string, header map[string]int) []UserInfo {
	list := []UserInfo{}

	var newUser = &UserInfo{}

	for key := range lookupFields {
		var propName = lookupFields[key]
		var propIndex = header[propName]

		if propIndex != -1 {
			newUser = newUser.setProperty(lookupFields[key], user[propIndex])
		}
	}

	list = append(list, *newUser)

	return list
}

// This function could be potentially used to filter out records. At the moment it's not doing much
// besides appending its content to the next channel.
func Reduce(mapList chan []UserInfo, sendFinalValue chan []UserInfo) {
	final := []UserInfo{}

	for list := range mapList {
		final = append(final, list...)
	}

	sendFinalValue <- final
}

// This Function combines the structs processed by the Reducer
func Combiner(users []UserInfo, finalMap map[string]UserInfo) map[string]UserInfo {
	for _, user := range users {
		_, exists := finalMap[user.Id]

		if !exists {
			finalMap[user.Id] = user
		} else {
			var existingUser = finalMap[user.Id]

			userMap, _ := StructToMap(user)
			existingUserMap, _ := StructToMap(existingUser)

			for k, v := range userMap {
				if v != "" {
					existingUserMap[k] = v
				}
			}

			var mergedUser UserInfo
			Transcode(existingUserMap, &mergedUser)
			finalMap[user.Id] = mergedUser
		}
	}

	return finalMap
}

func ParseCSVFiles() map[string]UserInfo {
	var csvs []*csv.Reader
	files, err := os.ReadDir(*flagDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range files {
		if filepath.Ext(fi.Name()) == ".csv" {
			f, err := os.Open(filepath.Join(*flagDir, fi.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			csvs = append(csvs, csv.NewReader(f))
		}
	}

	var flattenedUsers []UserInfo

	for _, r := range csvs {
		lines, err := r.ReadAll()

		if err != nil {
			panic(err)
		}

		data := lines[1:]
		header := lines[0]
		headerIndexes := GetHeaderIndexes(header)

		mappedUsers := make(chan []UserInfo)
		reducedUsers := make(chan []UserInfo)

		var wg sync.WaitGroup

		wg.Add(len(data))

		for _, line := range data {
			go func(user []string) {
				defer wg.Done()
				mappedUsers <- Mapper(user, headerIndexes)
			}(line)
		}

		go Reduce(mappedUsers, reducedUsers)

		wg.Wait()
		close(mappedUsers)

		flattenedUsers = append(flattenedUsers, <-reducedUsers...)
	}

	var usersMap = make(map[string]UserInfo)
	Combiner(flattenedUsers, usersMap)

	return usersMap
}

func WriteOutput(usersMap map[string]UserInfo) {
	csvwriter := csv.NewWriter(os.Stdout)
	_ = csvwriter.Write(lookupFields)
	var userIds = GetUniqUserIds(usersMap)

	for _, userId := range userIds {
		var outputRow []string

		for _, prop := range lookupFields {
			var user = usersMap[strconv.Itoa(userId)]
			outputRow = append(outputRow, GetField(&user, prop))
		}

		_ = csvwriter.Write(outputRow)
	}

	csvwriter.Flush()
}
