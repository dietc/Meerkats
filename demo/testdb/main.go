package main

import (
	"github.com/boltdb/bolt"
	"log"
	"strings"
	"regexp"
)

const ROOT string = "serverdb/"

var db *bolt.DB

func main() {
	var root string = suffix(safe(ROOT))
    var err error
    db, err = bolt.Open(root + ".storage/my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
    if err := db.Update(func(tx *bolt.Tx) error{
		b := bucket([]byte("space"), tx)
		//b.Get b.Put b.Stats().KeyN b.Delete
        //////////////////////////////////////
		b.ForEach(func(k, v []byte) error {
			b.Delete(k)
			return nil
		})
 
		b1 := bucket([]byte("previous"), tx)
		//b.Get b.Put b.Stats().KeyN b.Delete
        //////////////////////////////////////
		b1.ForEach(func(k, v []byte) error {
			b1.Delete(k)
			return nil
		})

		b2 := bucket([]byte("files"), tx)
		//b.Get b.Put b.Stats().KeyN b.Delete
        //////////////////////////////////////
		b2.ForEach(func(k, v []byte) error {
			b2.Delete(k)
			return nil
		})
		return nil
    }); err != nil {
		log.Fatal(err)
	}
}

func suffix(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func safe(path string) string {
	path = regexp.MustCompile("\\\\+").ReplaceAllString(path, "/")
	return path
}

func bucket(key []byte, tx *bolt.Tx) *bolt.Bucket {
	b, err := tx.CreateBucketIfNotExists(key)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
