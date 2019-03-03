package main

import (
    "log"
    "net"
    "os"
    "syscall"
    "os/signal"
    "p"
    "s"
    "util"
    
)

const PORT string = "4356"

func main() {
    initialize()
    start()
}

func initialize() {
    s.StartDB()
    err := util.CreateDirIfNotExisted(s.STO)
    if err != nil {
		log.Fatal(err)
    }
}

func start() {
    var err error
    var channel chan os.Signal = make(chan os.Signal)
    
    l,err :=  net.Listen("tcp4", ":" + PORT)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    log.Println("Server is listening at PORT:" + PORT + "...")
    signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
    go func(){
        <-channel
        log.Println("Server closed.")
        os.Exit(1)
    }()
    
    for {
        c, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go p.Handle(c)
    }
}

