package util

import (
    "crypto/md5"
    "strings"
    //"safe"
    "regexp"
    "os"
    //"log"
)

func Md5(data []byte, key []byte) [16]byte{
    return md5.Sum(append(data, key...))
}

func Digest(data []byte) [16]byte{
    return md5.Sum(data)
}

func CheckMd5(rmd5 []byte, data []byte, key []byte) (ret bool) {
	var cmd5 [16]byte
	cmd5 = md5.Sum(append(data, key...))
    ret = compare(rmd5, cmd5[:])
	return
}

func compare(d1 []byte, d2 []byte) (ret bool) {
    ret = true
    if len(d1) != len(d2) {
        ret = false
        return 
    }
    for k,v := range d1 {
        if d2[k]^v == 1 {
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



