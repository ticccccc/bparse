package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ticccccc/bparse/models"
)

func main() {
	mid := os.Args[1]
	if mid == "" {
		fmt.Println("mid is empty!")
		return
	}

	jpth := "data/" + mid + ".json"
	vpth := "data/" + mid + ".csv"

	if _, err := os.Stat(jpth); os.IsExist(err) {
		fmt.Println(jpth + " is exist!")
		fmt.Println("program do nothing")
		return
	}

	jf, err := os.Create(jpth)
	if err != nil {
		fmt.Println("create " + jpth + " failed")
		fmt.Println("program do nothing")
		return
	}
	defer jf.Close()

	if _, err := os.Stat(vpth); os.IsExist(err) {
		fmt.Println(jpth + " is exist!")
		fmt.Println("program do nothing")
		return
	}
	vf, err := os.Create(vpth)
	if err != nil {
		fmt.Println("create " + vpth + "failed")
		fmt.Println("program do nothing")
		return
	}
	defer vf.Close()

	example := models.NewAuthor(mid)
	example.GetInfo()
	j, _ := json.Marshal(example)

	jf.Write(j)
	for _, video := range example.Videos {
		vf.WriteString(video.Bvid + "\n")
	}
}
