package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/karlseguin/ccache"

	"archive/zip"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/go-errors/errors"
	"go.uber.org/zap"
)

func Start(port int) {
	api := rest.NewApi()

	api.Use(rest.DefaultCommonStack...)

	router, err := rest.MakeRouter(
		//peer 提交Partner的BalanceProof,更新Partner的余额
		rest.Get("/cloud-server/api/assignid", AssignID),
		rest.Post("/cloud-server/api/log/:address/:id", Log),
		rest.Post("/cloud-server/api/download/:address", dbDownload),
		rest.Post("/cloud-server/api/upload", dbUpLoad),
	)
	if err != nil {
		log.Fatalf("maker router :%s", err)
	}
	api.SetApp(router)
	listen := fmt.Sprintf("0.0.0.0:%d", port)
	log.Fatalf("http listen and serve :%s", http.ListenAndServe(listen, api.MakeHandler()))
}

func dbDownload(w rest.ResponseWriter, r *rest.Request) {
	/*fn:=fmt.Sprintf("%s.tar.gz",r.PathParam("address"))
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> dbDownload\t,file = [%s]", fn))
	}()
	userReqFile:=filepath.Join(logdir,"upload",fn)
	b,_:=ioutil.ReadFile(userReqFile)
	w.WriteJson(b)*/
	/*filename:=r.PathParam("address")
	out,err:=os.Create(filename)
	if err!=nil{
		w.WriteJson(err)
	}
	defer out.Close()*/

	filename := fmt.Sprintf("%s.zip", r.PathParam("address"))
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> dbDownload\t,file = [%s]", filename))
	}()
	userReqFile := filepath.Join(logdir, "upload", r.PathParam("address"))
	zipFileName := userReqFile + ".zip"
	err := Zip(userReqFile, zipFileName)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"dbDownload err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	time.Sleep(time.Second)

	f, err := ioutil.ReadFile(zipFileName)
	if err != nil {
		w.WriteJson(err.Error())
		return
	}
	h := w.Header()
	h.Set("Content-type", "application/octet-stream")
	h.Set("Content-Disposition", "attachment;filename="+zipFileName)
	w.WriteJson(f)

	//w.WriteJson(userReqFile + ".zip")
}

func PrintTime() string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "] "
}

func AssignID(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(NextID())
}

func dbUpLoad(w rest.ResponseWriter, r *rest.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"dbUpLoad err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	defer file.Close()
	/*zr, err := gzip.NewReader(file)
	if err != nil {
		fmt.Fprintln(w.(http.ResponseWriter), err)
		return
	}
	defer zr.Close()
	handler.Filename=path.Base(handler.Filename)
	f, err := os.OpenFile(path.Join(logdir,"upload", handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, zr)
	fmt.Fprintln(w.(http.ResponseWriter), "success")
	fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> dbUpLoad\t,file = [%s]", handler.Filename))*/
	obj, err := NewZipHandle(file, handler.Size)
	obj.OnHandle(CreateFile)
	err = obj.Handle()
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"dbUpLoad Handle err= [%s]", err))
		w.WriteJson(err.Error())
		return

	}
	fmt.Fprintln(w.(http.ResponseWriter), "success")
	fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> dbUpLoad\t,file = [%s]", handler.Filename))
}

func CreateFile(z *Ziphandle) (err error) {
	if z == nil {
		return errors.New("Uninitialized Ziphandle")
	}
	if z.File == nil {
		return errors.New("no file")
	}
	for _, tf := range z.File {
		for _, fi := range tf {
			if !IsFileExist(fi.Path) {
				err = os.MkdirAll(fi.Path, os.ModePerm)
			}
			f := fi.File
			srcFile, err := f.Open()
			defer srcFile.Close()
			if err != nil {
				return err

			}

			newFile, err := os.Create(filepath.Join(logdir, "upload", f.Name))
			defer newFile.Close()
			if err != nil {
				return err
			}
			_, err = io.Copy(newFile, srcFile)
			if err != nil {
				return err
			}

		}
	}
	return
}

func Zip(src_dir string, zip_file_name string) (err error) {

	dir, err := ioutil.ReadDir(src_dir)
	if err != nil {
		return errors.New(fmt.Sprintf("ioutil.ReadDir err: %s", zap.Error(err)))
	}
	if len(dir) == 0 {
		return errors.New(src_dir + " is empty dir!")
	}
	os.RemoveAll(zip_file_name)
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()
	// zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {
		if path == src_dir {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}

		return err
	})
	return err
}

var cache = ccache.New(ccache.Configure().MaxSize(50).ItemsToPrune(5).OnDelete(func(item *ccache.Item) {
	item.Value().(*os.File).Close()
}))

func Log(w rest.ResponseWriter, r *rest.Request) {
	address := r.PathParam("address")
	id := r.PathParam("id")
	log.Printf("address=%s,id=%s\n", address, id)
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, "body err", http.StatusConflict)
		return
	}
	doLog(address, id, w, msg)
}

func doLog(address, id string, w rest.ResponseWriter, msg []byte) {
	if len(address) <= 0 || len(id) <= 0 {
		rest.Error(w, "arg error ", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("%s-%s", address, id)
	it := cache.Get(key)
	if it == nil {
		filename := fmt.Sprintf("%s-%s.log", address, id)
		idFile, err := os.OpenFile(filepath.Join(logdir, filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			rest.Error(w, fmt.Sprintf("OpenFile for file %s err %s", filename, err), http.StatusInternalServerError)
			return
		}
		cache.Set(key, idFile, time.Second)
		it = cache.Get(key)
	}
	it.Value().(*os.File).Write(msg)
}

//解析接受的文件
type Ziphandle struct {
	f         io.ReaderAt
	fileSize  int64
	File      map[string][]*ZipFileStruct
	callbacks []callback //注入解析函数
}

type callback func(handle *Ziphandle) (err error)

type ZipFileStruct struct {
	FileName string
	Path     string
	FileType string
	File     *zip.File
}

//zhuru
func (z *Ziphandle) OnHandle(cb ...callback) {
	if z.callbacks == nil {
		z.callbacks = make([]callback, 0)
	}
	if len(cb) > 0 {
		z.callbacks = append(z.callbacks, cb...)
	}
}

func (z *Ziphandle) Handle() (err error) {
	if z == nil || z.f == nil || z.File == nil {
		return errors.New("parameter error")
	}
	for _, cb := range z.callbacks {
		if err = cb(z); err != nil {
			return err
		}
	}
	return
}

func NewZipHandle(f io.ReaderAt, size int64) (*Ziphandle, error) {
	res := &Ziphandle{
		f:        f,
		fileSize: size,
		File:     make(map[string][]*ZipFileStruct),
	}
	err := res.ZipParse()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (z *Ziphandle) ZipParse() error {
	zipFile, err := zip.NewReader(z.f, z.fileSize)
	if err != nil {
		return err
	}
	for _, zf := range zipFile.File {
		info := zf.FileInfo()
		if info.IsDir() {
			continue
		} else {
			resFile := &ZipFileStruct{
				FileName: info.Name(),
				Path:     filepath.Join(logdir, "upload", strings.TrimRight(zf.FileHeader.Name, info.Name())),
				//Path:strings.TrimRight(zf.FileHeader.Name,info.Name()),
				File:     zf,
				FileType: "",
			}
			//fmt.Println(resFile.Path)
			ext := strings.TrimLeft(filepath.Ext(info.Name()), ".")
			if ext == "" {
				continue
			}
			resFile.FileType = ext
			z.File[ext] = append(z.File[ext], resFile)
		}
	}
	return nil
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
