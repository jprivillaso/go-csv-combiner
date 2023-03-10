package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func bootstrap() {
	var header = []string{"Id", "FirstName", "LastName", "Email", "Phone"}

	for j := 1; j <= 4; j++ {
		f, err := os.Create(fmt.Sprintf("./data/generated/%d.csv", j))
		defer f.Close()

		if err != nil {
			log.Fatalln("failed to open file", err)
		}

		w := csv.NewWriter(f)
		defer w.Flush()

		if err := w.Write(header); err != nil {
			log.Fatalln("error writing record to file", err)
		}

		for i := 1; i <= 1000000; i++ {
			var row = []string{strconv.Itoa(i), "Juan", "Rivillas", "", ""}
			w.Write(row)
		}
	}
}
