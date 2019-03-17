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
    "strings"
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
    var fStatus,sStatus StatusTable
    var fNew,sNew NewTable
    var task Task = Task{}
    //var fileNames []string
    //var matched *set.Set = set.New()
    //var fileListRemain []string
    //var spaceRemain []string
    //var bytes []byte
    //var err error
    
    FetchFileList(fileList, data)//
    
    //delete type folder
    for k,v := range fileList {
        if v.Typ == 2 {
            delete(fileList, k)
        }
    }
    
    for _,obj := range fileList {
        cmdUpload(&task, obj.Name, obj.Digest) 
    } 
    bytes,_ := json.Marshal(task) 
    return bytes;
    
    //space initialization
    space = fetchSpaceInfo()
    //log.Println(space)   
    pre = FetchPreInfo(device)
    //log.Println(pre)
    fStatus, fNew = CheckCurrentStatus(pre, fileList)
    
    pre = FetchPreInfo(device)
    sStatus, sNew = CheckSpaceStatus(pre, space)
    
    Decide(fStatus, sStatus, fNew, sNew)
    return nil
}

func Decide(f,s StatusTable, fn,sn NewTable) {
    //log.Println("fStatus", f)
    //log.Println("fNew", fn)
    //log.Println("sStatus", s)
    //log.Println("sNew", sn)
    for name,digest := range fn {
        v,ok := sn[name]
        if ok {
            if util.Compare(v[:], digest[:]) {
                log.Println("do nothing")    
            } else {
                log.Println("download(conflict)", name)
            }
        } else {
            log.Println("upload: ", name)
        }
    } 
    
    
    for name,status := range f {
        if status.Typ == 0 {
            if s[name].Typ > 0 {
                if s[name].Typ == 3 {
                    log.Println("others rename: ", s[name].FromName, " => ", s[name].ToName)
                }
                if s[name].Typ == 2 {
                    log.Println("others modified: ", s[name].FromName)
                } 
                if s[name].Typ == 4 {
                    log.Println("others deleted: ", s[name].FromName)
                }
            }   
        } else {
            if s[name].Typ == 0 {
                if status.Typ == 3 {
                    log.Println("user rename: ", status.FromName, " => ", status.ToName)
                }
                if status.Typ == 2 {
                    log.Println("user modified: ", status.FromName)
                }
                if status.Typ == 4 {
                    log.Println("user deleted: ", status.FromName)
                }
            }
            
        
        }
        
    }
}

func CheckSpaceStatus(pre *Pre, space Space) (StatusTable, NewTable) {
    var pretable StatusTable = make(StatusTable)
    var newtable NewTable = make(NewTable)
    var uncheckedSpace *set.Set = set.New()
    for name,u := range space {
        val,ok := pre.P[name]
        if ok {           
            if val.Idx == u.Idx{
                //no change   
                var stat Status = Status{}
                stat.Typ = 0
                pretable[name] = stat
                pre.Cmp[val.Digest].Remove(name)
                pre.Checked[name] = true               
                continue
            } 
        }        
        uncheckedSpace.Add(name)
    }
    for _,name := range uncheckedSpace.List() {
        //renamed
        set, ok := pre.Cmp[space[name].Digest]
        if ok && set.Len() > 0 {
            var stat Status = Status{}
            stat.Typ = 3
            stat.FromName = getSameDirName(set, name)
            stat.FromDigest, stat.ToDigest = space[name].Digest,space[name].Digest
            stat.ToName = name
            pretable[stat.FromName] = stat
            set.Remove(stat.FromName)
            pre.Checked[stat.FromName] = true
            uncheckedSpace.Remove(name)
        }
    }
    for _,name := range uncheckedSpace.List() {
        //modified
        flag, ok := pre.Checked[name]
        if ok && !flag {
            var stat Status = Status{}
            stat.Typ = 2
            stat.FromName,stat.ToName = name,name
            stat.FromDigest = pre.P[name].Digest
            stat.ToDigest = space[name].Digest
            pretable[name] = stat
            pre.Cmp[pre.P[name].Digest].Remove(name)
            pre.Checked[name] = true
            uncheckedSpace.Remove(name)
        }
    }
    //others added
    for _,name := range uncheckedSpace.List() {
        newtable[name] = space[name].Digest
    }
    //others deleted
    for name, flag := range pre.Checked {
        if flag == true {
            continue
        }
        var stat Status = Status{}
        stat.Typ = 4
        stat.FromName = name
        stat.FromDigest = pre.P[name].Digest
        pretable[name] = stat
    }
    
    return pretable, newtable
}

func CheckCurrentStatus(pre *Pre, fileList map[string]FileObj) (StatusTable, NewTable) {
    var pretable StatusTable = make(StatusTable)
    var newtable NewTable = make(NewTable)
    var uncheckedFileSet *set.Set = set.New()   
    //initialize checkedFileList and single out none-changed file
    for name,obj := range fileList {
        //no change file
        set, ok := pre.Cmp[obj.Digest]
        if ok && set.Has(name) {
            var stat Status = Status{}
            stat.Typ = 0
            pretable[name] = stat
            set.Remove(name)
            pre.Checked[name] = true
            continue
        }
        uncheckedFileSet.Add(name)
    }
    
    //renamed file
    for _,name := range uncheckedFileSet.List() {
        set, ok := pre.Cmp[fileList[name].Digest]
        if ok && set.Len() > 0 {
            var stat Status = Status{}
            stat.Typ = 3
            stat.FromName = getSameDirName(set, name)
            stat.FromDigest,stat.ToDigest = pre.P[stat.FromName].Digest,pre.P[stat.FromName].Digest
            stat.ToName = name
            pretable[stat.FromName] = stat
            set.Remove(stat.FromName)
            pre.Checked[stat.FromName] = true
            uncheckedFileSet.Remove(name)
        }            
    }
    
    //modified file
    for _,name := range uncheckedFileSet.List() {
        flag, ok := pre.Checked[name]
        if ok && !flag {
            var stat Status = Status{}
            stat.Typ = 2
            stat.FromName,stat.ToName = name,name
            stat.FromDigest = pre.P[name].Digest
            stat.ToDigest = fileList[name].Digest
            pretable[name] = stat
            pre.Cmp[pre.P[name].Digest].Remove(name)
            pre.Checked[name] = true
            uncheckedFileSet.Remove(name)
        }    
    }
    
    //user add
    for _,name := range uncheckedFileSet.List() {
        obj := fileList[name]
        newtable[name] = obj.Digest
    }
    
    //user delete
    for name, flag := range pre.Checked {
        if flag == true {
            continue
        }
        var stat Status = Status{}
        stat.Typ = 4
        stat.FromName = name
        stat.FromDigest = pre.P[name].Digest
        pretable[name] = stat
    }

    return pretable, newtable
}

func getSameDirName(set *set.Set, name string) string{
    index := strings.LastIndex(name, "/")
    if index == -1{
        for _,target := range set.List() {
            if strings.LastIndex(target, "/") == -1 {
                 return target
            }
        }
    } else {
        for _,target := range set.List() {
            if strings.LastIndex(target, "/") == index && strings.HasPrefix(name, target[:index]) {
                return target
            }
        }    
    }
    return set.GetFirst()
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

func ProcessFileDownload(device byte) [][]byte {
    var res [][]byte
    data1,_ := ioutil.ReadFile("clientA/large.jpeg")
    resdata1 := divideFile(data1, "large.jpeg")
    res = append(res, resdata1...)
    data2,_ := ioutil.ReadFile("clientA/small")
    resdata2 := divideFile(data2, "small")
    res = append(res, resdata2...)
    return res
}

func divideFile(data []byte, name string) [][]byte{
    //65234   //65234 + 65534 + ..
    var res [][]byte
    var length int = len(data)
    var tmp []byte
    var part FilePart = FilePart{}
    if length <= 65234 {
    //single packet
        part.Name = name
        part.Num = uint16(1)
        part.Index = uint16(0)
        tmp,_ = json.Marshal(part)
        var stuff []byte = make([]byte, 300-len(tmp))
        tmp = append(tmp, stuff...)
        tmp = append(tmp, data...)
        res = append(res, tmp)
    } else {
    //multiple packets
        var src []byte = data[:65234]
        length = len(data[65234:])
        var num int = length/65534
        var re  int = length%65534
        if re != 0 {
            num++
        }
        part.Name = name
        part.Num = uint16(num+1)
        part.Index = uint16(0)
        tmp,_ = json.Marshal(part)
        var stuff []byte = make([]byte, 300-len(tmp))
        src = append(src, tmp...)
        src = append(src, stuff...)
        res = append(res, src)
        for i:=0;i<num;i++ {
            if i == num -1 {
                res = append(res, data[65234+65534*i:])
            } else {
                res = append(res, data[65234+65534*i:65234+65534*i+65534])
            }
        }
    }
    return res
}

func cmdUpload(task *Task, name string, digest [16]byte) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Digest = digest
    co.Cmd = 1
    co.Ext = []byte{}
    *task = append(*task, co)
}

func cmdUploadModified() {

}

func cmdDownloadModified() {

}

func cmdDownload() {


}

func FetchFileList(fileList map[string]FileObj, data []byte) {
    var objList []FileObj
    json.Unmarshal(data, &objList)
    for _,obj := range objList{
        fileList[obj.Name] = obj
    }
} 

func fetchSpaceInfo() Space{
    var s Space = make(Space)
    var info []byte
    var fi *FileInfo = new(FileInfo)
    var tmp map[string]int64
    
    tmp = GetSpace()
    for name, idx:= range tmp{
        info = GetFileInfo(idx)
        json.Unmarshal(info, fi)
        s[name] = new(Unit)
        s[name].Digest = fi.Digest
        s[name].Idx = idx
    }
    return s
}

func FetchPreInfo(device byte) *Pre{
    var p *Pre = new(Pre)
    var tmp map[string]int64
    var bytes []byte
    var info []byte
    var fi *FileInfo = new(FileInfo)
    
    p.Device = device
    bytes = GetPre(device)
    json.Unmarshal(bytes, &tmp)
    p.Checked = make(map[string]bool, len(p.P))
    p.Cmp = make(map[[16]byte]*set.Set)
    p.P = make(map[string]*Unit)
    for name, idx := range tmp{
        info = GetFileInfo(idx)
        json.Unmarshal(info, fi)
        p.P[name] = new(Unit)
        p.P[name].Digest = fi.Digest
        p.P[name].Idx = idx
        //Cmp
        if p.Cmp[fi.Digest] == nil {
            p.Cmp[fi.Digest] = set.New()
        }
        p.Cmp[fi.Digest].Add(name)
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

func searchDigestInSpace() {

}

func searchDigestInPre() {

}




