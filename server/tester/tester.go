package main

import (
	"s"
	"util"
	"log"
	"os"
)

func main(){
	s.StartDB()
    err := util.CreateDirIfNotExisted(s.STO)
    if err != nil {
		log.Fatal(err)
    }
	if os.Args[1] == "clear" {
        s.Clear()
		return
	}
	var key []byte
	if os.Args[2] == "A" {
		key = []byte{0x01}
	} else {
		key = []byte{0x02}
	}
	
	if os.Args[1] == "show" {
		log.Println(key)
	    s.Show(key)
	}
	
}