package main

import (
	"flag"
	"log"
	"ranet-clone/cfg"
	"ranet-clone/dl"
	"ranet-clone/ocr"
	"sync"
)

var (
	baseUrl              = "https://russianplanes.net/images/"
	baseDir              = flag.String("dir", "", "full dir path to download to, eg /raspberry/img/")
	mode                 = flag.String("mode", "download", "func to do")
	tessDataPrefix       = flag.String("tessdata", "/usr/share/tessdata", "path to tessdata")
	threads              = flag.Int("threads", 4, "threads to use for downloading")
	minImg         int64 = 100000
	maxImg         int64 = 301605
)

func main() {
	flag.Parse()

	dir := *baseDir
	if len(dir) == 0 {
		log.Fatalln("-dir not set")
	}

	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	cfg.LoadConfig()
	go cfg.SetupConfigSaving()

	log.Println("computing chunked paths")
	paths, err := dl.GenerateChunkedPaths(dir, *threads, minImg, maxImg)
	log.Printf("finished computing chunked paths (%v)\n", len(paths))

	if err != nil {
		log.Fatalln(err)
	}

	switch *mode {
	case "download":
		modeDownload(dir, paths)
	case "ocr":
		modeOcr(dir, paths)
	default:
		log.Fatalln(*mode + " is not a recognized mode")
	}
}

func modeOcr(dir string, paths [][]cfg.ImageInfo) {
	if err := ocr.InitClient(*tessDataPrefix); err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		log.Printf("starting ocr thread %v\n", i)
		wg.Add(1)
		go ocr.ProcessImages(&wg, i, paths[i], dir)
	}

	log.Println("waiting for ocr to finish")
	wg.Wait()
	log.Println("finished processing ocr")
}

func modeDownload(dir string, paths [][]cfg.ImageInfo) {
	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		log.Printf("starting download thread %v\n", i)
		wg.Add(1)
		go dl.DownloadFiles(&wg, i, paths[i], dir, baseUrl)
	}

	log.Println("waiting for downloads to finish")
	wg.Wait()
	log.Println("finished downloading")
}
