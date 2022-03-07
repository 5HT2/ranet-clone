package dl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"ranet-clone/cfg"
	"ranet-clone/threads"
	"strconv"
	"sync"
)

func GenerateChunkedPaths(dir string, nThreads int) ([][]cfg.ImageInfo, error) {
	arr, err := GeneratePaths(dir, threads.MinImg, threads.MaxImg)
	if err != nil {
		return nil, err
	}

	return threads.ChunkSlice(arr, len(arr)/nThreads), nil
}

// GeneratePaths generates the images sub-path in a given range.
// For example, 999-1001 will produce [to1000/000999.jpg, to1000/001000.jpg, to1000/001001.jpg].
// The first image is 1, and the last image is 301691.
func GeneratePaths(dir string, from, to int64) (arr []cfg.ImageInfo, err error) {
	if to < from {
		return arr, errors.New("to cannot be less than from")
	}

	if from < threads.MinImg || to > threads.MaxImg {
		return arr, errors.New("to or from outside of range (" + strconv.FormatInt(threads.MinImg, 10) + "-" + strconv.FormatInt(threads.MaxImg, 10) + ")")
	}

	files := threads.GetFiles(dir)

	for i := from; i <= to; i++ {
		dir := strconv.FormatInt(i+(1000-i%1000), 10)
		img := fmt.Sprintf("%06d", i) + ".jpg"

		if !contains(files, img) {
			arr = append(arr, cfg.ImageInfo{Path: "to" + dir + "/" + img, Name: img})
		}
	}
	cfg.UpdateNumDownloaded(int64(len(files)), true)

	return arr, err
}

func DownloadFiles(wg *sync.WaitGroup, thread int, p []cfg.ImageInfo, dir, baseUrl string) {
	defer threads.LogPanic()
	defer wg.Done()
	defer log.Printf("done thread %v\n", thread)

	log.Printf("thread %v will download %v files\n", thread, len(p))
	for _, i := range p {
		if cfg.InDlQueue(i) {
			continue
		}

		cfg.AddToDlQueue(i)
		DownloadFile(i, dir, baseUrl)
		cfg.RemoveFromDlQueue(i)
	}
}

func DownloadFile(p cfg.ImageInfo, dir, baseUrl string) {
	out, err := os.Create(dir + p.Name)
	defer out.Close()
	if err != nil {
		// if you fail to make the file, you've run out of drive space or don't have perms
		log.Println("failed to make file " + p.Name)
		return
	}

	resp, err := http.Get(baseUrl + p.Path)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("status %v returned for %s\n", resp.StatusCode, p.Name)
		return
	}

	total, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("downloaded: %s (%vB missing)\n", p.Name, n-total)

	if n-total == 0 {
		p.Size = n
		cfg.AddCompletedDownload(p)
		cfg.UpdateNumDownloaded(1, false)
	} else {
		log.Printf("missing file chunks for " + p.Name)
	}
}

func contains(s []os.FileInfo, e string) bool {
	for _, a := range s {
		if a.Name() == e {
			return true
		}
	}
	return false
}
