package util

import (
    "crypto/md5"
    "strings"
    //"safe"
    "regexp"
    "os"
    //"log"
    "log"
)

func Md5(data []byte, key []byte) [16]byte{
    var tmp []byte
    tmp = append(tmp, data...)
    tmp = append(tmp, key...)
    return md5.Sum(tmp)
}

func Digest(data []byte) [16]byte{
    return md5.Sum(data)
}

func CheckMd5(rmd5 []byte, data []byte, key []byte) (ret bool) {
	var cmd5 [16]byte
    var tmp []byte
    tmp = append(tmp, data...)
    tmp = append(tmp, key...)
	cmd5 = md5.Sum(tmp)
    //log.Printf("%x %x", rmd5, cmd5)
    ret = Compare(rmd5, cmd5[:])
    //log.Printf("content: %x", data)
    log.Printf("%x", rmd5)
    log.Printf("%x", cmd5)
    log.Println("checking md5:", ret)
	return
}

func Compare(d1 []byte, d2 []byte) (ret bool) {
    ret = true
    if len(d1) != len(d2) {
        ret = false
        return 
    }
    for k,v := range d1 {   
        if d2[k]^v != 0 {
            ret = false
            break
        }
    }
    return
}

func Suffix(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func Safe(path string) string {
    path = regexp.MustCompile("\\\\+").ReplaceAllString(path, "/")
    return path
}

func CreateDirIfNotExisted(path string) error {
    var err error
    if _, err = os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, os.ModeDir)
    }
    return err
}

func GetBlocksHashList(data []byte) []byte{
    var res []byte
    var length int = len(data)
    var parts_num int = length/1024
    var hash [16]byte
    //we can ensure length > 1024 so parts_num > 1
    for i:=0;i<parts_num;i++ {
       hash = md5.Sum(data[1024*i:1024*i+1024])
       res = append(res, hash[:]...)
    }    
    return res
}


