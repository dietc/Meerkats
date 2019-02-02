//server
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
	"regexp"
	"strings"
	"github.com/boltdb/bolt"
    "strconv"
    "errors"
    "set"
    "net"
    "bufio"
    "io"
    "crypto/hmac"
	"crypto/md5"
    "os"
    "time"
)

const ROOT string = "client"
type digest string

type FileObj struct{
    Name string
    Typ  uint8
	Digest [16]byte
}

type FileInfo struct{
    Name   string
    Digest [16]byte
    Dir    string
}

type Previous struct{
    Device string
    P map[string]int64
    Cmp map[digest][]string
    Checked map[string]bool
}

type CmdObj struct{
    Name   string
    Digest [16]byte
    Cmd    byte
}

type task map[string]CmdObj

type Uploaded struct{
    Name    string
    Content []byte
    Digest  [16]byte
}

func NewPrevious(device string) *Previous {
    var p *Previous = new(Previous)
    var bytes []byte
    p.Device = device
    if err := db.View(func(tx *bolt.Tx) error{
        bytes = tx.Bucket([]byte("previous")).Get([]byte(device))
        return nil
	}); err != nil {
        log.Fatal(err)
	}
    json.Unmarshal(bytes, &p.P)
    p.Checked = make(map[string]bool, len(p.P))
    p.Cmp = make(map[digest][]string)
    var fi *FileInfo = new(FileInfo)
    for name, idx := range p.P{
        fi = fetchFileInfo(idx)
        //Cmp
        p.Cmp[digest(fi.Digest[:])] = append(p.Cmp[digest(fi.Digest[:])], name)
        //Checked
        p.Checked[name] = false
    }
    return p
}
func (pre Previous)checkIndex(name string, idx int64) (error, bool) {
    if val, ok := pre.P[name];!ok{
        return errors.New("previous info not found"), false
    } else {
        if val == idx {
            return nil, true
        } else {
            return nil, false
        }
    }
}
func (pre Previous)checkDigest(digest []byte) {
    return
}
func (pre *Previous)checkRename(name string, d [16]byte) bool {
    arr, _ := pre.Cmp[digest(d[:])]
    if len(arr) < 1 {
        //no digest matched    
    } else if len(arr) == 1 {
        //potentially renamed file identified
        if Space[arr[0]] == pre.P[arr[0]] {
            //rename identified
            //update
            return true
        }
    } else {
       // var sN int = strings.Count(name, "/")
       // for i,k := range arr {
       //     if strings.Count(k, "/") == sN {
       //         pre.Cmp[digest(d[:])] = eraseSliceEle(pre.Cmp[digest(d[:])], i)
       //     }
       // }
    }
    //fmt.Println("cmp:", pre.Cmp)
    return false
}

/** list received from client **/
var FileList map[string]FileObj

/** keys from list **/
var FileNames []string

/** current storage space **/
var Space map[string]int64

var db *bolt.DB

/** device token **/
var device string = "A"

/** intersection of list and current space **/
var InSpace []string

var NotInSpace []string

var matched *set.Set = set.New()

var rF []string

var rS []string

func init(){
    FileList = make(map[string]FileObj)
	Space = make(map[string]int64)
    
	var root string = suffix(safe(ROOT))
    
    //receive file list
	//bytes, err := ioutil.ReadFile(root + ".storage/.list")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//json.Unmarshal(bytes, &FileList)
    //TODO:remove empty dir
    //for k,v := range FileList {
    //    if v.Typ == 2 {
    //        delete(FileList, k)
    //    }
    //}
    var err error
    db, err = bolt.Open(root + ".storage/my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
    
    //database initialization
	if err := db.Update(func(tx *bolt.Tx) error{
        //*db:space
		b := bucket([]byte("space"), tx)
		//b.Get b.Put b.Stats().KeyN b.Delete
        //////////////////////////////////////
        fmt.Println("Current space:")
		b.ForEach(func(k, v []byte) error {
		    fmt.Printf("%s => %s\n", k, v)
			return nil
		})
        //////////////////////////////////////
        c := b.Cursor()
        k, v := c.First()
        for k!=nil || v!=nil {
            i, _ := strconv.ParseInt(string(v), 10, 64)
            Space[string(k)] = i
            k,v  = c.Next()            
        }
        
        
        //*db:previous
        b1 := bucket([]byte("previous"), tx)
        //////////////////////////////////////////
        fmt.Println("Current previous:")
		b1.ForEach(func(k, v []byte) error {
		    fmt.Printf("%s => %s\n", k, v)
			return nil
		})
        //////////////////////////////////////////
        /**var Previous map[string]int64 = make(map[string]int64)
        Previous["a"] = 2
        bytes, _ := json.Marshal(Previous)
        b.Put([]byte(device), bytes)*/
        
        
        //*db:files 
        b2:= bucket([]byte("files"), tx)
        //////////////////////////////////////////
        fmt.Println("Current files:")
		b2.ForEach(func(k, v []byte) error {
		    fmt.Printf("%s => %s\n", k, v)
			return nil
		})
        //////////////////////////////////////////
        //fi := FileInfo{"a", [16]byte{'f'}, "dir2a"}
        //bytes,_:= json.Marshal(fi)
        //b.Put([]byte("1"), bytes)
        //var fi *FileInfo = new(FileInfo)
        //res := b.Get([]byte("1"))
        //json.Unmarshal(res, &fi)
        //fmt.Println(fi.Digest, fi.Dir, fi.Name)
        
		return nil
	}); err != nil {
		log.Fatal(err)
	}

}

func main(){  
    establish()
}

func process(c net.Conn){
    //space initialization
    if err := db.View(func(tx *bolt.Tx) error{
        b := tx.Bucket([]byte("space"))
        c := b.Cursor()
        k, v := c.First()
        for k!=nil || v!=nil {
            i, _ := strconv.ParseInt(string(v), 10, 64)
            Space[string(k)] = i
            k,v  = c.Next()            
        }
        return nil
	}); err != nil {
        log.Fatal(err)
	}
    
    var t task = task{}
    log.Println("Checking uploaded list...")
    //var lN int = len(FileList)
    //var sN int = len(Space)
    //fmt.Println("Current Number:")
    //fmt.Println(lN, sN)
    var pre = NewPrevious(device)
    for name,index := range Space {
        if v, ok := FileList[name]; !ok{
            rS = append(rS, name)
        } else {
            matched.Add(name)
            //name matched 
            var fi *FileInfo = fetchFileInfo(index)
            //fmt.Println(string(v.Digest[:]))
            //fmt.Println(string(fi.Digest[:]))
            if string(v.Digest[:]) == string(fi.Digest[:]) {  
                err, flag := pre.checkIndex(name, index)
                if err == nil && flag {
                    //file not changed
                } else {
                    //update previous info
                }
            } else {
                err, flag := pre.checkIndex(name, index)
                if err != nil {
                    // lost previous infomation
                } else if flag {
                    // upload
                    fmt.Println("you have to upload " + name)
                } else {
                    fmt.Println(name + "is conflicted with others change")
                    // conflict download
                }
            }
            //fmt.Println(name, "=>", index, v)
        }
    }
    //array rF
    //fmt.Println("matched", matched.List())
    for k := range FileList {
        FileNames = append(FileNames, k)
        if !matched.Has(k) {
            rF = append(rF, k)
        }
    }
    
    //fmt.Println("rf", rF)
    var rflag bool
    for k,v := range rF {
        rflag = pre.checkRename(FileList[v].Name, FileList[v].Digest)
        if rflag {
            fmt.Println(k)
            //rename?            
        } else {
            cmdUploadFull(t, FileList[v].Name, FileList[v].Digest)
        }
    }   
    
    //rS
    for _,v := range rS {
        err, flag := pre.checkIndex(v, Space[v])
        if err == nil && flag {
            deleteFile(v)
            cmdDeleteBackup(t, v)  
        } else {
            //confilct delete
        }
    }
    
    
    if len(t) > 0 {
        var bytes []byte
        bytes,_ = json.Marshal(t)
        log.Println("Sending back commands to client...")
        send(c, bytes, 0x10)
        log.Println("Succeed!")
    } else {
        log.Println("Sending back commands to client...")
        send(c, []byte(""), 0x20)
        log.Println("Succeed!")
    }
    
    defer func(){
        FileList = map[string]FileObj{}
        rF = []string{}
        rS = []string{}
        matched.Clear()
        FileNames = []string{} 
        Space = map[string]int64{}
    
    }()
    //show()
}

func save(files []byte) {
    var up Uploaded = Uploaded{}
    json.Unmarshal(files, &up)
    log.Println("Received uploaded file [" + up.Name + "]")
    indexFile(up)
    log.Println("Successfully saved and indexed.")
    //ioutil.WriteFile("server/")
}

func deleteFile(name string) {
    if err := db.Update(func(tx *bolt.Tx) error{
        b := bucket([]byte("space"), tx)
        idx := b.Get([]byte(name))
        b.Delete([]byte(name))
        
        b1 := bucket([]byte("previous"), tx)
        var pre map[string]int64 = make(map[string]int64)
        if tmp := b1.Get([]byte(device)); len(tmp) != 0 {
             json.Unmarshal(tmp, &pre)
        }
        delete(pre, name)
        preb,_ := json.Marshal(pre)
        b1.Put([]byte(device), preb)
        
        b2 := bucket([]byte("files"), tx)
        info := b2.Get(idx)
        var fi *FileInfo = new(FileInfo) 
        json.Unmarshal(info, fi)
        os.Remove(fi.Dir)
        b2.Delete(idx)
        return nil
    }); err != nil {
		log.Fatal(err)
	}
}

func indexFile(up Uploaded){
    var index int64 = time.Now().UnixNano()
    var dir string = "server/" + strconv.FormatInt(index, 10)
    var err error = ioutil.WriteFile(dir, up.Content, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if err := db.Update(func(tx *bolt.Tx) error{
		b := bucket([]byte("space"), tx)
        b.Put([]byte(up.Name), []byte(strconv.FormatInt(index,10)))
        
        b1 := bucket([]byte("previous"), tx)
        var pre map[string]int64 = make(map[string]int64)
        if tmp := b1.Get([]byte(device)); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }     
        pre[up.Name] = index
        preb, _ := json.Marshal(pre)
        b1.Put([]byte(device), preb)  
        
        b2 := bucket([]byte("files"), tx)
        var fi FileInfo= FileInfo{up.Name, up.Digest, dir}
        fib, _ := json.Marshal(fi)
        b2.Put([]byte(strconv.FormatInt(index,10)), fib)
        return nil
    }); err != nil {
		log.Fatal(err)
	}
    
    
}

func send(c net.Conn, b []byte, by byte){
    bufferWriter := bufio.NewWriter(c)
    bufferWriter.Write(enPacket(b, by))
    bufferWriter.Flush()
}

func show(){
    fmt.Println("FN:", FileNames)
    fmt.Println("SP:", Space)
    fmt.Println("M:", matched.List())
    fmt.Println("rS:", rS)
    fmt.Println("rF:", rF)
}

func fetchFileInfo(idx int64) *FileInfo{
    var fi *FileInfo = new(FileInfo)
    var res []byte
    i := strconv.FormatInt(idx, 10)
    
    if err := db.View(func(tx *bolt.Tx) error{
        res = tx.Bucket([]byte("files")).Get([]byte(i))
        return nil
	}); err != nil {
        log.Fatal(err)
	}
    json.Unmarshal(res, fi)
    //if idx == 2 {
    //   fi.Name = "b"
    //   fi.Digest = [16]byte{12,32,43,1,4,4,4,5,6,7,8,9,0,87,65,34}
    //   fi.Dir = "dir"
    //}
    return fi
}

func eraseSliceIntEle(arr []int, idx int) []int{
    var newarr []int
    for k, v := range arr {
        if k == idx {
            continue
        }
        newarr = append(newarr, v)
    }
    return newarr
}

func eraseSliceStrEle(arr []string, idx int) []string{
    var newarr []string
    for k, v := range arr {
        if k == idx {
            continue
        }
        newarr = append(newarr, v)
    }
    return newarr
}


func checkValue(bn , key, val []byte) (error,bool) {
    return nil, true
}

func checkExist(bn []byte, key []byte) bool{
    var res []byte
    if err := db.View(func(tx *bolt.Tx) error{
        b:= bucket(bn, tx)
        res = b.Get(key)
        return nil
    }); err != nil {
        log.Fatal(err)
    }
    if len(res) ==  0 {
        return false
    } else {
        return true
    }
    
}

//let client upload files in full 
func cmdUploadFull(t task, f string, d [16]byte) {
    var co *CmdObj = new(CmdObj)
    co.Name = f
    co.Digest = d
    co.Cmd = 1
    t[f] = *co
}

//notify client a file backup has been deleted
func cmdDeleteBackup(t task, f string) {
    var co *CmdObj = new(CmdObj)
    co.Name = f
    co.Cmd = 2
    t[f] = *co
}

func bucket(key []byte, tx *bolt.Tx) *bolt.Bucket {
	b, err := tx.CreateBucketIfNotExists(key)
	if err != nil {
		log.Fatal(err)
	}
	return b
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

func establish() {
    l,err := net.Listen("tcp4", ":4356")
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    for {
        c, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go handleConnection(c)
    }
}

func handleConnection(c net.Conn) {
    log.Printf("Serving %s\n", c.RemoteAddr().String()) 
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
                handleData(c, rbuffer)
            }
            state, length, rbuffer = restore()             
        }
    }
}

func handleData(c net.Conn, buf []byte) {
    switch buf[0]{
        //pull request
        case 0x1:
        log.Println("Received pull request")
        json.Unmarshal(buf[1:], &FileList)
        for k,v := range FileList {
            if v.Typ == 2 {
                delete(FileList, k)
            }
        }
        process(c)
        //upload file request
        case 0x2:
        save(buf[1:])
    }
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

