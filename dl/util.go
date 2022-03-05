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

type ImagePath struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// GeneratePaths generates the images sub-path in a given range.
// For example, 100000-100002 will produce [to101000/100000.jpg, to101000/100001.jpg, to101000/100002.jpg].
// The first image is 100000, and the last image is 301605.
func GeneratePaths(from, to int64) (arr []ImagePath, err error) {
	if to < from {
		return arr, errors.New("to cannot be less than from")
	}

	if from < 100000 || to > 301605 {
		return arr, errors.New("to or from outside of range (100000-301605)")
	}

	for i := from; i <= to; i++ {
		dir := strconv.FormatInt(i+(1000-i%1000), 10)
		img := strconv.FormatInt(i, 10) + ".jpg"

		arr = append(arr, ImagePath{"to" + dir + "/" + img, img, 0})
	}

	return arr, err
}

func DownloadFile(p ImagePath, dir, baseURL string) {
	out, err := os.Create(dir + p.Name)
	defer out.Close()
	if err != nil {
		// if you fail to make the file, you've run out of drive space or don't have perms
		panic(err)
	}

	resp, err := http.Get(baseURL + p.Path)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	total, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return
	}
	log.Printf("downloaded: %v\n", total-n)

	if total-n == 0 {
		p.Size = n
		cfg.AddCompletedDownload(p)
	} else {
		panic("missing file chunks")
	}
}
