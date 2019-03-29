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
    "crypto/md5"
)


var ServerDir string
var Index     int64
var tmpbuffer []byte


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
    
    //for _,obj := range fileList {
        //cmdUpload(&task, obj.Name, obj.Digest) 
        //log.Println(obj.Name, obj.Digest)
    //} 
    //bytes,_ := json.Marshal(task) 
    //return nil;
    
    //space initialization
    space = fetchSpaceInfo()
    //log.Println(space)   
    pre = FetchPreInfo(device)
    //log.Println(pre)
    /** compare fileList and pre fStatus=>status table fNew=>new file table**/
    /** 0 => no change     **/
    /** 3 => rename        **/
    /** 2 => modify        **/
    /** 4 => user delete   **/
    fStatus, fNew = CheckCurrentStatus(pre, fileList)
    //log.Println(fStatus, fNew)
    pre = FetchPreInfo(device)
    sStatus, sNew = CheckSpaceStatus(pre, space)
    pre = FetchPreInfo(device)
    Decide(fStatus, sStatus, fNew, sNew, fileList, space, pre, &task, device)
    bytes,_ := json.Marshal(task) 
    log.Printf("%s", bytes)
    return bytes
}

func Decide(f,s StatusTable, fn,sn NewTable, fileList map[string]FileObj, space Space, pre *Pre, task *Task, device byte) {
      log.Println("f", f)
      //log.Println("fNew", fn)
      log.Println("s", s)
      //log.Println("sNew", sn)
      //return
    /** check conflict situation in user created files **/
    for name, digest := range fn {
        //others created as well
        if odigest,ok := sn[name];ok {
            if !util.Compare(digest[:], odigest[:]) {
                log.Println("others created(conflict): ", name) 
                internalAddPreInfo(device, name, space[name].Idx)
                cmdBackup(task, name)
                cmdDownload(task, name, space[name].Idx)
            } else {
                internalAddPreInfo(device, name, space[name].Idx)
            }    
        } else {
            if _,ok := space[name];ok {
            //others have renamed to a same name
                log.Println("user back up: ", name)
                cmdBackup(task, name)
            } else {
            //user can upload
                log.Println("user created: ", name)
                cmdUpload(task, name, digest)
            }         
        }
    }
    
    /** check conflict situation in others created files **/
    for name, _ := range sn {
        if _,ok := fn[name];!ok{
            if _,flag := f[name];!flag {
                log.Println("others created: ", name)
                cmdDownload(task, name, space[name].Idx)
                internalAddPreInfo(device, name, space[name].Idx)
            }        
        }    
    }
    
    
    for name,status := range f {
        if status.Typ != 0 {
            continue
        }
        /** Current user has done nothing to last version **/
        if s[name].Typ > 0 {
            switch s[name].Typ{
                case 0x03:                
                log.Println("others renamed: ", s[name].FromName, " => ", s[name].ToName)
                cmdRename(task, name, s[name].ToName)
                internalUpdatePreName(space[s[name].ToName].Idx, s[name].ToName, s[name].FromName, device)
                case 0x02:
                log.Println("others modified: ", s[name].FromName)
                cmdDownloadModified(task, name, space[s[name].FromName].Idx)
                internalUpdatePreName(space[s[name].FromName].Idx, s[name].FromName, s[name].FromName, device)
                case 0x04:
                log.Println("others deleted: ", s[name].FromName)
                cmdDelete(task, s[name].FromName)
                internalDeletePreName(s[name].FromName, device)
            }
        }
        delete(f, name)
        delete(s, name)
    }
    
    for name,status := range s {
        if status.Typ != 0 {
            continue
        }
        if f[name].Typ > 0 {
            switch f[name].Typ {
                case 0x03:
                log.Println("user renamed: ", f[name].FromName, " => ", f[name].ToName)
                internalRename(f[name].ToName, f[name].FromName, pre.P[f[name].FromName].Idx, device)
                case 0x02:
                log.Println("user modified: ", f[name].FromName)
                if GetFileLength(space[name].Idx) > 1024 {
                   //rsync
                   cmdUploadModified(task, name, fileList[name].Digest, space[name].Idx)
                } else {
                   //without rsync
                   cmdUpload(task, name, fileList[name].Digest)
                }
                //internal function
                case 0x04:
                log.Println("user deleted: ", f[name].FromName)
                internalDelete(f[name].FromName, device)
            }
        }
        delete(f, name)
        delete(s, name)
    }
    
     
    for name,status := range s{
        if status.Typ == 0x03 && f[name].Typ == 0x03{
            //others renamed && user renamed     
            cmdRename(task, f[name].ToName, status.ToName)
            internalUpdatePreName(space[status.ToName].Idx, status.ToName, status.FromName, device)
            log.Println("others renamed(user renamed conflict): ", f[name].ToName, "=>", status.ToName)
        }
        if status.Typ == 0x03 && f[name].Typ == 0x02 {
            //others renamed && user modifed
            cmdRename(task, f[name].FromName, status.ToName)
            internalUpdatePreName(space[status.ToName].Idx, status.ToName, f[name].FromName, device)
            cmdUpload(task, status.ToName, fileList[f[name].FromName].Digest)
           // cmdUploadModified(task, status.ToName, fileList[name].Digest)
            log.Println("others renamed(user modifed conflict): ")
        }
        if status.Typ == 0x03 && f[name].Typ == 0x04 {
            //others renamed && user deleted
            internalDelete(status.ToName, device)
            internalDeletePreName(f[name].FromName, device)
            log.Println("others renamed(user deleted conflict): ")
        }
        if status.Typ == 0x02 && f[name].Typ == 0x03 {
            //others modified && user renamed
            internalRename(f[name].ToName, f[name].FromName, space[f[name].FromName].Idx, device)
            //cmdDownloadModified
            log.Println("others modifed(user renamed conflict): ")
        }
        if status.Typ == 0x02 && f[name].Typ == 0x02 {
            var tmp []byte 
            for _,b := range f[name].ToDigest {
                tmp = append(tmp, b)
            }
            //others modified && user modified
            if !util.Compare(status.ToDigest[:], tmp) {
                 cmdBackup(task, f[name].FromName)
                 cmdDownloadModified(task, f[name].FromName, space[f[name].FromName].Idx)
                 internalUpdatePreName(space[f[name].FromName].Idx, f[name].FromName, f[name].FromName, device)
                 log.Println("others modifed(user modified conflict): ")
             } else {
                 internalUpdatePreName(space[f[name].FromName].Idx, f[name].FromName, f[name].FromName, device)
             }
        }
        if status.Typ == 0x02 && f[name].Typ == 0x04 {
            //others modifed && user deleted
            internalUpdatePreName(space[f[name].FromName].Idx, f[name].FromName, f[name].FromName, device)
            //cmdDownload()
            log.Println("others modifed(user deleted conflict): ")
        }
        if status.Typ == 0x04 && f[name].Typ == 0x03 {
            //others deleted && user renamed
            cmdDelete(task, f[name].ToName)
            internalDeletePreName(f[name].FromName, device)
            log.Println("others deleted(user renamed conflict): ")
        }
        if status.Typ == 0x04 && f[name].Typ == 0x02 {
            //others deleted && user modified
            cmdBackup(task, f[name].FromName)
            internalDeletePreName(f[name].FromName, device)
            log.Println("others  deleted(user modifed conflict): ")
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
                stat.FromName = name
                stat.FromDigest = val.Digest
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
    
    log.Println("fileList")
    for _,v := range fileList{
        log.Println(v.Name, v.Digest)
    }
    
    log.Println("pre.P")
    for n,v := range pre.P{
        log.Println(n, v.Digest, v.Idx)
    }
    
    log.Println("pre.Cmp")
    for n,v := range pre.Cmp{
        log.Println(n, v.List())
    }
    
    log.Println("pre.Checked")
    for n,v := range pre.Checked{
        log.Println(n, v)
    }
    
    //initialize checkedFileList and single out none-changed file
    for name,obj := range fileList {
        //no change file
        set, ok := pre.Cmp[obj.Digest]
        if ok && set.Has(name) {
            var stat Status = Status{}
            stat.Typ = 0
            stat.FromName = name
            stat.FromDigest = obj.Digest
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
    part.Name = util.Safe(part.Name)
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

func ProcessFileDownload(index []byte, device byte) [][]byte {
    var res [][]byte
    var info []byte
    var fi *FileInfo = new(FileInfo)
    var idx int64 
    idx,err := strconv.ParseInt(string(index), 10, 64)
    if err == nil {
        info = GetFileInfo(idx)
        json.Unmarshal(info, fi)      
        data,_ := ioutil.ReadFile(fi.Dir)
        resdata := divideFile(data, fi.Name)
        res = append(res, resdata...)
        var num int = 0
        for _,i := range res{
            log.Println(len(i))
            num += len(i)
        }
        log.Println(num)
    }
    return res
}

func ProcessModifiedFileUpload(data []byte, device byte) [][]byte {
    var length int = len(data)
    var fi *FileInfo = new(FileInfo)
    if length < 300 {
        var res []byte
        idx,_ := strconv.ParseInt(string(data), 10, 64)
        info := GetFileInfo(idx)
        json.Unmarshal(info, fi)
        data,_ := ioutil.ReadFile(fi.Dir)
        var length int = len(data)
        var parts_num int = length/1024
        for i:=0;i<parts_num;i++ {
            hash := md5.Sum(data[1024*i:1024*i+1024])
            res = append(res, hash[:]...)
        }
        return divideFile(res, fi.Name)   
    } else {
        log.Printf("hihi%x", data)
        var list BlockList = BlockList{List:[]Block{}}
        var tmp []byte = data[0:1000]  
        var lastIdx int
        for i:=1; i <= 1000; i++{
            if tmp[1000-i] != 0x00 {
                lastIdx = 1000-i
                break
            } 
        }
        json.Unmarshal(data[0:lastIdx+1], &list)
        tmpbuffer = append(tmpbuffer, data[1000:]...)
        if list.Idx == list.Num - 1 {
            var tmpfile []byte
            var current int = 0
            var fi *FileInfo = new(FileInfo)
            space := fetchSpaceInfo()
            info := GetFileInfo(space[list.Name].Idx)
            json.Unmarshal(info, fi)
            data,_ := ioutil.ReadFile(fi.Dir)
            for _,item := range list.List {
                if item.Typ == 1 {
                    tmpfile = append(tmpfile, tmpbuffer[current:current + item.Len]...)
                    current += item.Len
                } else {
                    tmpfile = append(tmpfile, data[item.Idx*1024:item.Idx*1024+1024]...)
                }
            }
            log.Printf("new file hash:%x", md5.Sum(tmpfile))
            dir, _ := indexFile(list.Name,  md5.Sum(tmpfile), device)
            log.Printf("new file:%s", tmpfile)
            saveFile(dir, tmpfile)
            tmpbuffer = []byte{}
        }
        return nil
    }
   
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
        var src []byte
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
        src = append(src, data[:65234]...)
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

func internalRename(new_name string, old_name string, old_index int64, device byte) {
    var index int64 = time.Now().UnixNano()
    ReindexFile(old_index, new_name, index)
    UpdateName(index, device, new_name, old_name)
}

func internalUpdatePreName(new_index int64, new_name string, old_name string, device byte) {    
    UpdatePreName(new_index, new_name, old_name, device)
}

func internalDeletePreName(name string, device byte) {
    DeletePreName(name, device)
}

func internalDelete(name string, device byte)  {
    DeleteFile(name, device)
}

func internalAddPreInfo(device byte, name string, idx int64) {
    var tmp map[string]int64 = make(map[string]int64)
    var bytes []byte
    bytes = GetPre(device)
    json.Unmarshal(bytes, &tmp)
    tmp[name] = idx
    ret,_ := json.Marshal(tmp)
    SavePre(device, ret)
} 

func cmdUpload(task *Task, name string, digest [16]byte) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Digest = digest
    co.Cmd = 1
    co.Ext = ""
    *task = append(*task, co)
}

func cmdDownload(task *Task, name string, idx int64) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Cmd = 2
    co.Ext = strconv.FormatInt(idx, 10)
    *task = append(*task, co)
}

func cmdRename(task *Task, name string, new_name string) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Cmd = 3
    co.Ext = new_name
    *task = append(*task, co)
}

func cmdUploadModified(task *Task, name string, digest [16]byte, idx int64) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Digest = digest
    co.Cmd = 4
    co.Ext = strconv.FormatInt(idx, 10)
    *task = append(*task, co)
}


func cmdDownloadModified(task *Task, name string, idx int64) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Cmd = 5
    co.Ext = strconv.FormatInt(idx, 10)
    *task = append(*task, co)
}

func cmdDelete(task *Task, name string) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Cmd = 6
    *task = append(*task, co)
}

func cmdBackup(task *Task, name string) {
    var co *CmdObj = new(CmdObj)
    co.Name = name
    co.Cmd = 7 
    *task = append(*task, co)
}

func FetchFileList(fileList map[string]FileObj, data []byte) {
    var objList []FileObj
    json.Unmarshal(data, &objList)
    for _,obj := range objList{
        fileList[util.Safe(obj.Name)] = obj
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




