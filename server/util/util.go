package util

import (
    "crypto/md5"
)

func Md5(data []byte, key []byte) [16]byte{
    return md5.Sum(append(data, key...))
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

