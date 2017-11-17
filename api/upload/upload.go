package upload


import(
  "io"
  "os"
  "time"
  "path/filepath"
  "strconv"
  "net/http"
  "log"
)

func UploadHandler(r *http.Request,id string,foldername string) (map[string]interface{}){
  var isOk    bool
  var pathURL   string
  var dirURL    string
  var fileNames   []string

  var newFile string

pathURL = "../assets" // Your baseurl or whatever path you want.

 err := r.ParseMultipartForm(2097152) // grab the multipart form
 if err != nil {
   isOk = false
 }


 formdata := r.MultipartForm // ok, no problem so far, read the Form data

 //get the *fileheaders
 files := formdata.File["foto"] // grab the filenames
 status := r.FormValue("status")

 for i, _ := range files { // loop through the files one by one
   file, err := files[i].Open()
   defer file.Close()
   if err != nil {
     isOk = false
   }

   dirURL = filepath.Join("/",foldername,"/")

   gabung := filepath.Join(pathURL,dirURL)
   gabungabs,gaberr := filepath.Abs(gabung)
   if gaberr != nil{
     //log.Println(gaberr)
   }

   _,dirErr := os.Stat(gabungabs)
   if os.IsNotExist(dirErr){
       os.Mkdir(gabungabs,0777)
       //log.Println(dirErr)
   }

   dirURL = filepath.Join("/",id,"/")

   gabung = filepath.Join(gabung,dirURL)
   gabungabs,gaberr = filepath.Abs(gabung)
   if gaberr != nil{
     log.Println(gaberr)
   }

   _,dirErr = os.Stat(gabungabs)
   if os.IsNotExist(dirErr){
       os.Mkdir(gabungabs,0777)
       //log.Println(dirErr)
   }

  if(status == "bukti transfer"){
    newFile = strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + filepath.Ext(files[i].Filename)
  } else { 
    detik := time.Duration(int64(i*5))
    newFile = strconv.FormatInt(time.Now().Add(time.Second * detik).UTC().UnixNano(), 10) + filepath.Ext(files[i].Filename)
  }
  fileNames = append(fileNames,newFile)
  
  pathnya := filepath.Join(gabung,newFile)
  pathabs,abserr := filepath.Abs(pathnya)
  if abserr != nil{
     //log.Println(abserr)
   }

   dst, dstErr := os.Create(pathabs)
   if dstErr != nil {
     isOk = false
   }else {
     if _, ioErr := io.Copy(dst, file) ; ioErr != nil {
       isOk = false
     }else {
       isOk = true
     }
   }
   defer dst.Close()
 }
 if isOk {
   message := map[string]interface{} {
     "status"   :   201,
     "filename" :   fileNames, //customized new name of files comma separated
     "message"   :   "Upload sucess",
   }
   return message
 }
 errMsg := map[string]interface{} {
   "status" :   200,
   "filename":   nil,
   "message"  :   "Upload failed",
 }
 return errMsg
}

