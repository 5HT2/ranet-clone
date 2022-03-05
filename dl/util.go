package dl

import (
	"errors"
	"io"
	log "log"
	"net/http"
	"os"
	"ranet-clone/cfg"
	"runtime/debug"
	"strconv"
	"sync"
)

func GenerateChunkedPaths(dir string, threads int, from, to int64) ([][]cfg.ImageInfo, error) {
	arr, err := GeneratePaths(dir, from, to)
	if err != nil {
		return nil, err
	}

	return chunkSlice(arr, threads), nil
}

// GeneratePaths generates the images sub-path in a given range.
// For example, 100000-100002 will produce [to101000/100000.jpg, to101000/100001.jpg, to101000/100002.jpg].
// The first image is 100000, and the last image is 301605.
func GeneratePaths(dir string, from, to int64) (arr []cfg.ImageInfo, err error) {
	if to < from {
		return arr, errors.New("to cannot be less than from")
	}

	if from < 100000 || to > 301605 {
		return arr, errors.New("to or from outside of range (100000-301605)")
	}

	files := getFiles(dir)

	for i := from; i <= to; i++ {
		dir := strconv.FormatInt(i+(1000-i%1000), 10)
		img := strconv.FormatInt(i, 10) + ".jpg"

		if !contains(files, img) {
			arr = append(arr, cfg.ImageInfo{Path: "to" + dir + "/" + img, Name: img})
		}
	}
	cfg.UpdateNumDownloaded(int64(len(files)), true)

	return arr, err
}

func DownloadFiles(wg *sync.WaitGroup, p []cfg.ImageInfo, dir, baseUrl string) {
	defer logPanic()
	defer wg.Done()

	for _, i := range p {
		if cfg.InQueue(i) {
			continue
		}

		cfg.AddToQueue(i)
		DownloadFile(i, dir, baseUrl)
		cfg.RemoveFromQueue(i)
	}
}

func DownloadFile(p cfg.ImageInfo, dir, baseUrl string) {
	out, err := os.Create(dir + p.Name)
	defer out.Close()
	if err != nil {
		// if you fail to make the file, you've run out of drive space or don't have perms
		log.Println("failed to make file " + p.Name)
	}

	resp, err := http.Get(baseUrl + p.Path)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	n, err := io.Copy(out, resp.Body)
	total, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Println(err)
	}

	log.Printf("downloaded: %s (%vB missing)\n", p.Name, total-n)

	if total-n == 0 {
		p.Size = n
		cfg.AddCompletedDownload(p)
		cfg.UpdateNumDownloaded(1, false)
	} else {
		log.Printf("missing file chunks for " + p.Name)
	}
}

func getFiles(dir string) []os.FileInfo {
	f, err := os.Open(dir)
	if err != nil {
		log.Fatalln(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Fatalln(err)
	}

	return files
}

func contains(s []os.FileInfo, e string) bool {
	for _, a := range s {
		if a.Name() == e {
			return true
		}
	}
	return false
}

func chunkSlice(slice []cfg.ImageInfo, chunkSize int) [][]cfg.ImageInfo {
	var chunks [][]cfg.ImageInfo
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func logPanic() {
	if x := recover(); x != nil {
		// recovering from a panic; x contains whatever was passed to panic()
		log.Printf("panic: %s\n", debug.Stack())
	}
}
