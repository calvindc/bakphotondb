package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var id = 0
var idFileName = "id.dat"
var idFile *os.File
var lock sync.Mutex

func NextID() int {
	lock.Lock()
	defer lock.Unlock()
	id++
	_, err := idFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("seek err %s", err)
	}
	_, err = idFile.WriteString(strconv.Itoa(id))
	if err != nil {
		log.Printf("write err %s", err)
	}
	//err = idFile.Sync()
	//if err != nil {
	//	log.Printf("sync err %s", err)
	//}
	return id
}

// Exists returns true if directory||file exists
func Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}
func SetupDB(dir string) {
	var err error
	filename := filepath.Join(dir, idFileName)

	buf, err := ioutil.ReadFile(filename)
	if err == nil {
		id, err = strconv.Atoi(string(buf))
		if err != nil {
			panic(err)
		}
	}
	idFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	NextID()
}
