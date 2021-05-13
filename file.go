package common

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx"
)

//ExtHTML ExtHTML
const (
	EXTHTML = ".html"
	EXTXML  = ".xml"
	EXTJSON = ".json"
	EXTZIP  = ".zip"
	EXTIDOC = ".idoc"
	EXTZED  = ".zed"
	EXTPDF  = ".pdf"
)

// GetFilelist 列出目录下所有文件路径
func GetFilelist(path string) []string {
	fileList := make([]string, 0)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	return fileList
}

// ForeachFile 列出目录下所有文件
func ForeachFile(path string, pattern string, onFindFile func(fileName string) error) error {
	return filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if ok, _ := filepath.Match(pattern, f.Name()); ok {
			err := onFindFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ReadFileContext 读文件内容
func ReadFileContext(path string) string {
	fi, err := os.Open(path)
	NoneError(err)
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	NoneError(err)
	return string(fd)
}

// SaveFile 保存文件
func SaveFile(file string, content string) {
	f := ForceCreateFile(file)
	defer f.Close()

	_, err := f.WriteString(content)
	NoneError(err)
}

// SaveUtf8BomFile 保存utf-8 bom文件
func SaveUtf8BomFile(file string, content string) {
	fs := ForceCreateFile(file)
	defer fs.Close()

	fs.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	fs.WriteString(content)
}

// ForceCreateFile 创建文件
func ForceCreateFile(filename string) *os.File {
	dir := filepath.Dir(filename)
	dir = UNCPath(dir)
	err := os.MkdirAll(dir, 0755)
	ShowError(err)
	filename = UNCPath(filename)
	fs, err := os.Create(filename)
	ShowError(err)
	return fs
}

// SetLastModified 最后修改时间
func SetLastModified(file string, time time.Time) {
	os.Chtimes(file, time, time)
}

// GetlastModified 最后修改时间
func GetlastModified(file string) time.Time {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}
	return fileInfo.ModTime()
}

// CopyFile 复制
func CopyFile(dst, src string) error {
	if !FileExists(src) {
		return errors.New("file not exist `" + src + "`")
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out := ForceCreateFile(dst)
	if err != nil {
		fmt.Println("copy ", src)
		fmt.Println("to   ", dst)
		fmt.Println("Failed")
		return nil
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

// CopyDir 复制文件夹
func CopyDir(dst, src string) error {
	dst = UNCPath(dst)
	src = UNCPath(src)
	return filepath.Walk(src,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				NoneError(err)
			}
			if !f.IsDir() {
				rel, err := filepath.Rel(src, path)
				NoneError(err)
				return CopyFile(filepath.Join(dst, rel), path)
			}
			return nil
		})
}

// FileExists 文件是否存在
func FileExists(filename string) bool {
	finfo, err := os.Stat(filename)
	if err == nil && !finfo.IsDir() {
		return true
	}
	return false
}

// WriteFile 文件写
func WriteFile(file string, data []byte) {
	fs := ForceCreateFile(file)
	defer fs.Close()
	fs.Write(data)
}

// ReadJSON 文件转json
func ReadJSON(obj interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

// SaveJSON json转文件
func SaveJSON(filename string, obj interface{}) error {
	fs := ForceCreateFile(filename)
	defer fs.Close()

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = fs.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// ReadXML ReadXML
func ReadXML(obj interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

// SaveXML SaveXML
func SaveXML(filename string, obj interface{}) error {
	fs := ForceCreateFile(filename)
	defer fs.Close()

	data, err := xml.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = fs.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// GetFileInfo GetFileInfo
func GetFileInfo(filename string) (name, ext string, en bool) {
	file := filepath.Base(filename)
	ext = filepath.Ext(file)

	name = strings.TrimSuffix(file, ext)

	if strings.ToLower(filepath.Ext(name)) == ".en" {
		en = true
		ext = filepath.Ext(name) + ext
		name = strings.TrimSuffix(file, ext)
	}
	return name, ext, en
}

// FileToURLPath FileToURLPath
func FileToURLPath(file string) string {
	xURL, err := url.Parse(file)
	NoneError(err)
	return xURL.EscapedPath()
}

// URLToLocalFile URLToLocalFile
func URLToLocalFile(path string) string {
	xURL, err := url.Parse(path)
	NoneError(err)
	return SHA1EncodeFileName(xURL.Path)
}

// SHA1EncodeFileName SHA1EncodeFileName
func SHA1EncodeFileName(filename string) string {
	//for en picture: xxx.en.png

	dir := filepath.Dir(filename)

	filenameOnly, ext, _ := GetFileInfo(filename)

	t := sha1.New()
	_, err := io.WriteString(t, filenameOnly)
	NoneError(err)

	return dir + "/" + fmt.Sprintf("%x", t.Sum(nil)) + ext
}

// RemoveFile  RemoveFile
func RemoveFile(file string) {
	if !FileExists(file) {
		return
	}
	err := os.Remove(file)
	if err != nil {
		log.Print(file, " file remove Error!\n")
		log.Printf("%s", err)
	} else {
		log.Print(file, " file remove OK!\n")
	}
}

// FileOpen FileOpen
func FileOpen(filename string) ([]byte, error) {
	var (
		open             *os.File
		filedata         []byte
		openerr, readerr error
	)
	open, openerr = os.Open(filename)
	if openerr != nil {
		return nil, openerr
	}
	defer open.Close()
	filedata, readerr = ioutil.ReadAll(open)
	if readerr != nil {
		return nil, readerr
	}
	return filedata, nil
}

// CurDir CurDir
func CurDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

// IsFileSame IsFileSame
func IsFileSame(left, right string) bool {
	sFile, err := os.Open(left)
	if err != nil {
		return false
	}
	defer sFile.Close()

	dFile, err := os.Open(right)
	if err != nil {
		return false
	}
	defer dFile.Close()

	sBuffer := make([]byte, 512)
	dBuffer := make([]byte, 512)
	for {
		_, sErr := sFile.Read(sBuffer)
		_, dErr := dFile.Read(dBuffer)
		if sErr != nil || dErr != nil {
			if sErr != dErr {
				return false
			}
			if sErr == io.EOF {
				break
			}
		}
		if bytes.Equal(sBuffer, dBuffer) {
			continue
		}
		return false
	}
	return true
}

// DiffFile DiffFile
func DiffFile(left, right string) (diff string, same bool, err error) {
	_, err = exec.LookPath("diff")
	if err != nil {
		return "", false, err
	}

	cmdName, _ := exec.LookPath("diff")
	cmd := exec.Command(cmdName, left, right)
	output, errOfCmd := cmd.CombinedOutput()
	if errOfCmd != nil {
		return string(output), false, nil
	}
	return "", true, nil
}

// NoneError NoneError
func NoneError(err error) {
	fmt.Println(err)
}

// ShowError ShowError
func ShowError(err error) {
	fmt.Println(err)
}

// UNCPath UNCPath
func UNCPath(path string) string {
	return ""
}

// PathExists 文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// ScanFiles ScanFiles
func ScanFiles(dir string, fileType []string) []string {
	var allFile []string
	finfo, _ := ioutil.ReadDir(dir)
	for _, x := range finfo {
		realPath := dir + "/" + x.Name()
		if x.IsDir() {
			allFile = append(allFile, ScanFiles(realPath, fileType)...)
		} else {
			ext := path.Ext(x.Name())
			if Contains(ext, fileType) {
				allFile = append(allFile, realPath)
			}
		}
	}
	return allFile
}

// ListDir ListDir
func ListDir(path string) []string {
	var allFile []string
	finfo, _ := ioutil.ReadDir(path)
	for _, x := range finfo {
		realPath := path + "/" + x.Name()
		// realPath := x.Name()
		//fmt.Println(x.Name()," ",x.Size())
		if x.IsDir() {
			allFile = append(allFile, ListDir(realPath)...)
		} else {
			allFile = append(allFile, realPath)
		}
	}
	return allFile
}

// ListDirFilenane ListDirFilenane
func ListDirFilenane(path string) []string {
	var allFile []string
	finfo, _ := ioutil.ReadDir(path)
	for _, x := range finfo {
		realPath := x.Name()
		// realPath := x.Name()
		//fmt.Println(x.Name()," ",x.Size())
		if x.IsDir() {
			allFile = append(allFile, ListDir(realPath)...)
		} else {
			allFile = append(allFile, realPath)
		}
	}
	return allFile
}

// GetHTTPFile 从http读取文件
func GetHTTPFile(URL, filePath string) error {
	resp, err := http.Get(URL)
	body1, err := ioutil.ReadAll(resp.Body)
	out, err := os.Create(filePath)
	io.Copy(out, bytes.NewReader(body1))
	out.Close()
	return err
}

// ReadFile ReadFile
func ReadFile(filename string) ([]string, error) {
	lines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return lines, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, err
}

// ReadFile2String ReadFile2String
func ReadFile2String(filename string) (string, error) {
	file, err := os.Open(filename)
	fileString := ""
	const bufferSize = 16 * 1024
	buffer := make([]byte, bufferSize)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		fileString += string(buffer)
	}
	fileString += string(buffer)
	return fileString, err
}

// getFileExt getFileExt
func getFileExt(filename string) string {
	return path.Ext(filename)
}

//ReadExcelFile ReadExcelFile
func ReadExcelFile(filename string) (allRows [][]string, err error) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	sheets := xlsx.GetSheetList()
	for _, sheet := range sheets {
		rows, err := xlsx.GetRows(sheet)
		if err == nil {
			allRows = append(allRows, rows...)
		}
	}
	return allRows, err
}

// PrintExcel PrintExcel
func PrintExcel(xlsxFile string) {
	xlFile, err := xlsx.OpenFile(xlsxFile)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("Sheet Name: %s\n", sheet.Name)
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf(text)
			}
			fmt.Printf("\n")
		}
	}
}
