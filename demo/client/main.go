//client
package main

import (
	//"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"io/ioutil"
	"encoding/json"
	"log"
    "net"
    "bufio"
    "io"
    "fmt"
    "crypto/hmac"
	"crypto/md5"
)

type CmdObj struct{
    Name   string
    Digest [16]byte
    Cmd    byte
}

type task map[string]CmdObj

const ROOT string = "client"

type FileObj struct{
    Name string
    Typ  uint8
	Digest [16]byte	
}

var FileList map[string]FileObj

type FilePacket struct{
    Name string
    Digest [16]byte
    Content []byte
}

func init(){
	FileList = make(map[string]FileObj)
}

func main(){
    log.Println("Generating filesystem info...")
    var bytes []byte = generateFileList()
    log.Println("Succeed!")
    log.Println("Connecting to remote host...")
	addr, _:= net.ResolveTCPAddr("tcp", strings.TrimSpace("localhost:4356"))
    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
		log.Fatal(err)
    }
    log.Println("Succeed!")
    log.Println("Sending pull request...")
    bufferWriter := bufio.NewWriter(conn)
    bufferWriter.Write(enPacket(bytes, 0x1))
    bufferWriter.Flush()
    log.Println("Succeed!")
    log.Println("Listening...")
    var finished chan bool = make(chan bool)
    go handleConnection(conn, finished)
    <-finished
    fmt.Println("Task Finished!")
}

func enPacket(data []byte, typ byte) []byte{
    var cl uint16 = uint16(len(data) + 1)
    var tl int = int(cl) + 28
    var res []byte = make([]byte, tl)
    res[0] = 0x11
    res[1] = 0xff
    res[2] = 0x6c
    res[3] = 0x6f
    res[4] = 0x6e
    res[5] = 0x64
    res[6] = 0x6f
    res[7] = 0x6e
    res[8] = byte(cl >> 8)
    res[9] = byte(cl & 0xFF)
    //packet type 
    res[10] = typ
    copy(res[11:], data)
    mac := hmac.New(md5.New, []byte("aaaa"))
    mac.Write(res[10:(10+cl)])
    checksum := mac.Sum(nil)
    for i:=0;i<16;i++ {
        res[tl-18+i] = checksum[i]
    }
    res[tl-2] = 0xff
    res[tl-1] = 0xee
    return res
}


func generateFileList() []byte{
    var root string = suffix(safe(ROOT))
	var buffer []byte
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error{
		if path != root && !strings.HasPrefix(path, root + ".storage") {
			var obj *FileObj = new(FileObj)
			obj.Name = path[len(root):]
			if info.IsDir() {
			    obj.Typ = 2
			} else {
			    obj.Typ = 1
				buffer, err = ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				obj.Digest = md5.Sum(buffer)
			}
			FileList[obj.Name] = *obj
		}
		return nil
	})
	bytes, err := json.Marshal(FileList)
	if err != nil {
		log.Fatal(err)
	}
	//ioutil.WriteFile(root + ".storage/.list", bytes, 0644)
    return bytes
}



func safe(path string) string {
	path = regexp.MustCompile("\\\\+").ReplaceAllString(path, "/")
	return path
}

func suffix(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func handleConnection(c net.Conn, f chan bool) {
    state := byte(0x00)
    length := uint16(0)
    reader := bufio.NewReader(c)
    var err error
    var rbyte byte
    var rbuffer []byte
    cbuffer := make([]byte, 16)
    checksum := make([]byte, 16)
    for {
        if state != 0x0a && state != 0x0b {
            rbyte, err = reader.ReadByte()
        }
        if err != nil {
            if err == io.EOF {
                fmt.Printf("Closed by %s\n", c.RemoteAddr().String())
                return
            }
            continue
        }
        switch state {
            case 0x00:
            state = changeState(rbyte, 0x11, state)
            case 0x01:
            state = changeState(rbyte, 0xff, state)
            case 0x02:
            state = changeState(rbyte, 0x6c, state)
            case 0x03:
            state = changeState(rbyte, 0x6f, state)
            case 0x04:
            state = changeState(rbyte, 0x6e, state)
            case 0x05:
            state = changeState(rbyte, 0x64, state)
            case 0x06:
            state = changeState(rbyte, 0x6f, state)
            case 0x07:
            state = changeState(rbyte, 0x6e, state)
            /** packet length **/
            case 0x08:
            length += uint16(rbyte) * 256
            state++
            case 0x09:
            length += uint16(rbyte)
            //fmt.Println(length)
            rbuffer = make([]byte, length)
            state++
            case 0x0a:
            rbuffer, err = reader.Peek(int(length))
            if err != nil {
                state, length, rbuffer = restore() 
                continue
            }
            discarded, err := reader.Discard(int(length))
            if err != nil || discarded != int(length) {
                state, length, rbuffer = restore() 
                continue
            }
            state++
            case 0x0b:
            cbuffer, err = reader.Peek(16)
            //fmt.Println("received checksum:", cbuffer)
            if err != nil {
                state, length, rbuffer = restore() 
                continue
            }
            mac := hmac.New(md5.New, []byte("aaaa"))
            mac.Write(rbuffer)
            checksum = mac.Sum(nil)
            //fmt.Println("calculated checksum:", checksum)
            flag := hmac.Equal(cbuffer, checksum)
            if flag {
                discarded, err := reader.Discard(16)
                if err != nil || discarded != 16 {
                    state, length, rbuffer = restore() 
                    continue
                }
                state++
            } else {
                state, length, rbuffer = restore()
                continue
            }
            case 0x0c:
                state = changeState(rbyte, 0xff, state)
            case 0x0d:
            if rbyte == 0xee{
                /** buffer is ready to be used **/
                handleData(c, rbuffer, f)
            }
            state, length, rbuffer = restore()             
        }
    }
    c.Close()
}

func handleData(c net.Conn, buf []byte, f chan bool) {
    var t task = make(task)      
    switch buf[0] {
        case 0x10:
        json.Unmarshal(buf[1:], &t)
        for k,v := range t{
            if v.Cmd == 1 {
                //upload
                uploadFile(c, k, v.Digest)                
            }
            if v.Cmd == 2 {
                notifyBackupDeleted(v.Name)
            }
        }
        case 0x20:
        //do nothing
        log.Println("Nothing to do:Local repository already updated!")    
    }
    f <- true
}

func notifyBackupDeleted(name string) {
    log.Println("[" + name + "] has been deleted from remote server.")
}



func uploadFile(c net.Conn, name string, expected [16]byte) {
    bufferWriter := bufio.NewWriter(c) 
    log.Println("Uploading [" + name + "] ...")
    bytes, _ := ioutil.ReadFile("client/" + name)
    fmt.Println("client/" + name + ":", bytes)
    var digest [16]byte = md5.Sum(bytes)
    if string(digest[:]) == string(expected[:]) {
        var fp FilePacket = FilePacket{}
        fp.Name = name
        fp.Content = bytes
        fp.Digest = digest
        res,_ := json.Marshal(fp)
        bufferWriter.Write(enPacket(res, 0x2))
        bufferWriter.Flush()
        log.Println("Succeed!")
    } else {
        log.Println("Cancelled!")
    }
       
}

func changeState(r byte, target byte, state byte) (retstate byte) {
    if r == target {
        state++
    } else {
        state = 0x00
    }
    retstate = state
    return
}

func restore() (a byte, b uint16, c []byte) {
    a = 0 
    b = 0
    c = nil
    return
}
