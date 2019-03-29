package s

import(
    "github.com/boltdb/bolt"
	"util"
	"log"
    "strconv"
    "io/ioutil"
    "encoding/json"
)

const ROOT string = "serverdb"

const ENV string = "notesting"

var db *bolt.DB

func StartDB(){
    var root string = util.Suffix(util.Safe(ROOT))
	var err error
	
	err = util.CreateDirIfNotExisted(root + ".storage/")
    if err != nil {
		log.Fatal(err)
    }
	
    db, err = bolt.Open(root + ".storage/my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
    }
	
	//bucket establish
	if err = db.Update(func(tx *bolt.Tx) error{
		//bucket:space
	    bucket([]byte("space"), tx)
		//bucket:pre
        bucket([]byte("pre"), tx)
		//bucket:files
		bucket([]byte("files"), tx)
        //bucket:mapping
        bucket([]byte("mapping"), tx)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

//create
func bucket(key []byte, tx *bolt.Tx) *bolt.Bucket {
	b, err := tx.CreateBucketIfNotExists(key)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
	
func GetSpace() map[string]int64{
    var space map[string]int64 = make(map[string]int64)
    if err := db.View(func(tx *bolt.Tx) error{
        b := tx.Bucket([]byte("space"))
        c := b.Cursor()
        k, v := c.First()
        for k!=nil || v!=nil {
            i, _ := strconv.ParseInt(string(v), 10, 64)
            space[string(k)] = i
            k,v  = c.Next()            
        }
        return nil
	}); err != nil {
        log.Fatal(err)
	}
    return space
}
	
func GetPre(device byte) []byte{
    var bytes []byte
	if err := db.View(func(tx *bolt.Tx) error{
		bytes = tx.Bucket([]byte("pre")).Get([]byte{device})
        return nil
	}); err != nil {
        log.Fatal(err)
    }
    return bytes
}

func SavePre(device byte, data []byte) {
    if err := db.Update(func(tx *bolt.Tx) error{
        tx.Bucket([]byte("pre")).Put([]byte{device}, data)
        return nil
	}); err != nil {
        log.Fatal(err)
    }
}

func GetFileInfo(idx int64) []byte{
    var res []byte
    i := strconv.FormatInt(idx, 10)
    
    if err := db.View(func(tx *bolt.Tx) error{
        res = tx.Bucket([]byte("files")).Get([]byte(i))
        return nil
	}); err != nil {
        log.Fatal(err)
    }
    return res
}

func DeletePreName(name string, device byte) {
    if err := db.Update(func(tx *bolt.Tx) error{
        b:= bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b.Get([]byte{device}); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }
        delete(pre, name)
        preb, _ := json.Marshal(pre)
		b.Put([]byte{device}, preb) 
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}


func UpdatePreName(new_index int64, new_name string, old_name string, device byte) {
    if err := db.Update(func(tx *bolt.Tx) error{
        b:= bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b.Get([]byte{device}); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }
        delete(pre, old_name)
        pre[new_name] = new_index
        preb, _ := json.Marshal(pre)
		b.Put([]byte{device}, preb) 
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func UpdateName(new_index int64, device byte, name string, old_name string) {
    if err := db.Update(func(tx *bolt.Tx) error{
        b:= bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b.Get([]byte{device}); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }     
        delete(pre, old_name)
        pre[name] = new_index
        preb, _ := json.Marshal(pre)
		b.Put([]byte{device}, preb) 
        
        b1 := bucket([]byte("space"), tx)
        b1.Delete([]byte(old_name))
        b1.Put([]byte(name), []byte(strconv.FormatInt(new_index,10)))
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func ReindexFile(old_index int64, new_name string, new_index int64) {
    if err := db.Update(func(tx *bolt.Tx) error{
        var fi FileInfo= FileInfo{}
        b := bucket([]byte("files"), tx)
        info := b.Get([]byte(strconv.FormatInt(old_index,10)))
        json.Unmarshal(info, &fi)
        fi.Name = new_name
        fib, _ := json.Marshal(fi)
        b.Put([]byte(strconv.FormatInt(new_index,10)), fib)
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func IndexFile(name string, dir string, digest [16]byte, index int64, device byte) {
    if ENV == "testing" {
        return
    }
    if err := db.Update(func(tx *bolt.Tx) error{
		b1 := bucket([]byte("space"), tx)
        b1.Put([]byte(name), []byte(strconv.FormatInt(index,10)))
        
        b2 := bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b2.Get([]byte{device}); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }     
        pre[name] = index
        preb, _ := json.Marshal(pre)
		b2.Put([]byte{device}, preb)  
        
        b3 := bucket([]byte("files"), tx)
        var fi FileInfo= FileInfo{name, digest, dir}
        fib, _ := json.Marshal(fi)
        b3.Put([]byte(strconv.FormatInt(index,10)), fib)
        
        if !util.Compare(digest[:], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) {
            b4 := bucket([]byte("mapping"), tx)
            b4.Put(digest[:], []byte(strconv.FormatInt(index,10)))
        }        
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func DeleteFile(name string, device byte) {
    if err := db.Update(func(tx *bolt.Tx) error{
        b1 := bucket([]byte("space"), tx)
        b1.Delete([]byte(name))
        
        b2 := bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b2.Get([]byte{device}); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }
        delete(pre, name)
        preb, _ := json.Marshal(pre)
		b2.Put([]byte{device}, preb) 
        
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func SupplementDigest(index int64, dir string) {
    if ENV == "testing" {
        return
    }
    if err := db.Update(func(tx *bolt.Tx) error{
        b3 := bucket([]byte("files"), tx)
        b4 := bucket([]byte("mapping"), tx)
        var fi FileInfo= FileInfo{}
        if tmp := b3.Get([]byte(strconv.FormatInt(index,10))); len(tmp) != 0 {
            json.Unmarshal(tmp, &fi)            
        } 
        buffer, err := ioutil.ReadFile(dir)
		if err != nil {
			log.Fatal(err)
		}
        fi.Digest = util.Digest(buffer)
        fib, _ := json.Marshal(fi)
        b3.Put([]byte(strconv.FormatInt(index,10)), fib)
        b4.Put(fi.Digest[:], []byte(strconv.FormatInt(index,10)))
        return nil
    }); err != nil {
		log.Fatal(err)
    }
}

func GetFileLength(idx int64) int{
    var length int
    if err := db.Update(func(tx *bolt.Tx) error{
        b := bucket([]byte("files"), tx)
        var fi FileInfo= FileInfo{}
        if tmp := b.Get([]byte(strconv.FormatInt(idx,10))); len(tmp) != 0 {
            json.Unmarshal(tmp, &fi)            
        }
        data,_ := ioutil.ReadFile(fi.Dir)
        length = len(data)
        return nil
    }); err != nil {
		log.Fatal(err)
    }
    return length
}

//for testing:display data
func Show(key []byte) {
    if err := db.Update(func(tx *bolt.Tx) error{
        log.Println("Space:")
        b1 := bucket([]byte("space"), tx)
        c1 := b1.Cursor()
        k1, v1 := c1.First()
        for k1!=nil || v1!=nil {
            log.Println(string(k1), string(v1))
            k1,v1  = c1.Next()            
        }
        
        log.Println("Pre:")
        b2 := bucket([]byte("pre"), tx)
        var pre map[string]int64 = make(map[string]int64)
		if tmp := b2.Get(key); len(tmp) != 0 {
            json.Unmarshal(tmp, &pre)            
        }
        for k,v := range pre{
            log.Println(k, v)
        }
        
        log.Println("Files:")
        b3 := bucket([]byte("files"), tx)
        c3 := b3.Cursor()
        k3, v3 := c3.First()
        for k3!=nil || v3!=nil {
            var fi FileInfo= FileInfo{}
            log.Println(string(k3))
            json.Unmarshal(v3, &fi)
            log.Println(fi)
            k3,v3  = c3.Next()            
        }
    	return nil
    }); err != nil {
		log.Fatal(err)
    }
}

//for testing:remove all data from database
func Clear() {
    if err := db.Update(func(tx *bolt.Tx) error{
		b1 := bucket([]byte("space"), tx)
		b1.ForEach(func(k, v []byte) error {
			b1.Delete(k)
			return nil
		})
 
		b2 := bucket([]byte("pre"), tx)
		b2.ForEach(func(k, v []byte) error {
			b2.Delete(k)
			return nil
		})

		b3 := bucket([]byte("files"), tx)
		b3.ForEach(func(k, v []byte) error {
			b3.Delete(k)
			return nil
		})
		return nil
    }); err != nil {
		log.Fatal(err)
    }
}