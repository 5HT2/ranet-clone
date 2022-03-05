package dl

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"ranet-clone/cfg"
	"strconv"
)

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

func DownloadFile(p cfg.ImageInfo, dir, baseURL string) {
	out, err := os.Create(dir + p.Name)
	defer out.Close()
	if err != nil {
		// if you fail to make the file, you've run out of drive space or don't have perms
		panic(err)
	}

	resp, err := http.Get(baseURL + p.Path)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	n, err := io.Copy(out, resp.Body)
	total, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("downloaded: %v\n", total-n)

	if total-n == 0 {
		p.Size = n
		cfg.AddCompletedDownload(p)
		cfg.UpdateNumDownloaded(1, false)
	} else {
		panic("missing file chunks")
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
