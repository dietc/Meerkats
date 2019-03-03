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
	if os.Args[1] == "show" {
	    s.Show()
	}
	if os.Args[1] == "clear" {
        s.Clear()
	}
}