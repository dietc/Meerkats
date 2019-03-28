package s

import (
    "encoding/json"
    "log"
    "errors"
    "set"
    "fmt"
)

const STO string= "storage"

var ERRNF error = errors.New("previous info not found")

type FileObj struct{
    Name string
    Typ  uint8
	Digest [16]byte
}

type CmdObj struct{
    Name       string
    Digest     [16]byte
    Cmd        byte
	Ext        string
}

type Task []*CmdObj

type Space map[string]*Unit

type Unit struct{
    Digest [16]byte
    Idx    int64
}

type Pre struct{
    Device byte
    P map[string]*Unit
    Cmp map[[16]byte]*set.Set //checksum -> filenames
    Checked map[string]bool
}
func (p *Pre)CheckDelStatus(name string, d[16]byte) byte {
	if unit, ok := p.P[name];!ok {
        //new file
	    return 1
	} else {
		if d == unit.Digest {
            //deleted 
		    return 2
		} else {
            //deleted but with conflict
		    return 3
		}	  
	}
}
func (p *Pre)GetPreIdx(name string) (error, int64){
    if val, ok := p.P[name];!ok {
        return ERRNF, 0
    } else {
        return nil, val.Idx
    }
}
func (p *Pre)GetPreDigest(idx int64) [16]byte{
    var fi *FileInfo
    FetchFileInfo(fi, idx)
    return (*fi).Digest
}
func (p *Pre)String() string{
    var res string
    for name, unit := range p.P{
        res = res + name + " "
        res = res + fmt.Sprintf("%x", unit.Digest) + " "
        res = res + fmt.Sprintf("%d", unit.Idx) + " "
        res = res + "\n"
    }
    return res
}


type DigestStr string

type FileInfo struct{
    Name   string
    Digest [16]byte
    Dir    string
}
func FetchFileInfo(fi *FileInfo, idx int64) {
    err := json.Unmarshal(GetFileInfo(idx), fi)
    if err != nil {
        log.Fatal(err)
    }
}

type FilePart struct{
    Name  string
    Num   uint16
    Index uint16
}

type Uploaded struct{
    Name    string
    Content []byte
    Digest  [16]byte
}

type Status struct{
    Typ        byte
    FromName   string
    FromDigest [16]byte
    ToName     string
    ToDigest   [16]byte
}

type StatusTable map[string]Status
func (st StatusTable)String() (str string){
    for _, s:= range st{
        str = str + fmt.Sprintf("%d, %s => %s, %x => %x\n\n", s.Typ, s.FromName, s.ToName, s.FromDigest, s.ToDigest )
    }
    return str
}

type NewTable map[string][16]byte

type Block struct {
    Typ uint8
    Idx  int
    Len  int
}