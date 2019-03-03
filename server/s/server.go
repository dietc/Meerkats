package s

import(
    "net"
    "encoding/json"
    "strconv"
    "io/ioutil"
	//"bufio"
    "log"
    "os"
    "set"
    "time"
    "util"
)


var ServerDir string
var Index     int64


func ProcessTestingPacket(c net.Conn, data []byte, device byte) []byte{
	if string(data) == "hello" && device == 0x01 {
        return []byte("world")
	}
    return nil
}

func ProcessLocalListPacket(data []byte, device byte) []byte{
    //received list
    var fileList map[string]FileObj = make(map[string]FileObj) 
    var space Space
    var pre *Pre
    var task Task = Task{}
    var fileNames []string
    var matched *set.Set = set.New()
    var fileListRemain []string
    
    json.Unmarshal(data, &fileList)
        //delete type folder
    for k,v := range fileList {
        if v.Typ == 2 {
            delete(fileList, k)
        }
    }

    //space initialization
    space = InitializeSpace()
    //log.Println(space) 
    
    pre = fetchPreInfo(device)
   // for name, index := range Space{
    
   // }
    log.Println(space)
    for k := range fileList {
        fileNames = append(fileNames, k)
        if !matched.Has(k) {
            fileListRemain = append(fileListRemain, k)
        }
    }
    if len(fileListRemain) > 0 {
        var dstatus byte
        for _,name := range fileListRemain {
            dstatus = pre.CheckDelStatus(name, fileList[name].Digest)
            //log.Println(name, dstatus)
            switch dstatus{
            case 1:cmdUpload(task, name, fileList[name].Digest)
            }
        }    
    }
    
    //spaceRemain
    
    if len(task) > 0 {
        var bytes []byte
        bytes,_ = json.Marshal(task)
        return bytes    
    }
    
    return nil
}

func ProcessFileAssembly(data []byte, device byte) bool {
    var part FilePart=FilePart{}
    var lastIdx int
    var bytes []byte = data[:300]
    var digest [16]byte
    var response = false
    
    for i:=1; i <= 300; i++{
        if bytes[300-i] != 0x00 {
            lastIdx = 300-i
            break
        } 
    }
    json.Unmarshal(data[:lastIdx+1], &part)
    //log.Println("l", len(data))
    if part.Num <= 1 {
        var serverdir string
        digest = util.Digest(data[300:])
        serverdir, _ = indexFile(part.Name, digest, device)
        saveFile(serverdir, data[300:])
        response = true
    } else {
        if part.Index == 0 {
            ServerDir, Index = indexFile(part.Name, digest, device)
        }
        appendFile(ServerDir, data[300:])
        //last one
        if part.Index == part.Num-1 {
            SupplementDigest(Index, ServerDir)
            Index = 0
            ServerDir = ""
            response = true
        }
    }
    return response
}

func cmdUpload(task Task, name string, digest [16]byte) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Digest = digest
    co.Cmd = 1
    task[name] = *co
}

func fetchPreInfo(device byte) *Pre{
    var p *Pre = new(Pre)
    var bytes []byte
    var info []byte
    var fi *FileInfo = new(FileInfo)
    
    p.Device = device
    bytes = GetPre(device)
    json.Unmarshal(bytes, &p.P)
    p.Checked = make(map[string]bool, len(p.P))
    p.Cmp = make(map[DigestStr][]string)
    for name, idx := range p.P{
        info = GetFileInfo(idx)
        json.Unmarshal(info, fi)
        //Cmp
        p.Cmp[DigestStr(fi.Digest[:])] = append(p.Cmp[DigestStr(fi.Digest[:])], name)
        //Checked
        p.Checked[name] = false
    }
    return p    
}

func indexFile(name string,  digest [16]byte, device byte) (string, int64){
    var index int64 = time.Now().UnixNano()
    var indexS string = strconv.FormatInt(index, 10)
    //dir in server
    var dir string = STO + "/" + indexS
    //save information into database
    IndexFile(name, dir, digest, index, device)
    return dir, index
}

func saveFile(dir string,data []byte){    
    var err error = ioutil.WriteFile(dir, data, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func appendFile(dir string,data []byte){
    f, err := os.OpenFile(dir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write(data); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}


