package s

import (
    "encoding/json"
)

const STO string= "storage"

type FileObj struct{
    Name string
    Typ  uint8
	Digest [16]byte
}

type Space map[string]int64

type CmdObj struct{
    Name       string
    Digest     [16]byte
    Cmd        byte
	Ext        []byte
}

type Task []*CmdObj

type Pre struct{
    Device byte
    P map[string]int64
    Cmp map[DigestStr][]string
    Checked map[string]bool
}
func (p *Pre)CheckDelStatus(name string, d[16]byte) byte {
	if idx, ok := p.P[name];!ok {
        //new file
	    return 1
	} else {
        var fi *FileInfo
		var info []byte = GetFileInfo(idx)
        json.Unmarshal(info, fi) 
		if d == fi.Digest {
            //deleted 
		    return 2
		} else {
            //deleted but with conflict
		    return 3
		}	  
	}
}


type DigestStr string

type FileInfo struct{
    Name   string
    Digest [16]byte
    Dir    string
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