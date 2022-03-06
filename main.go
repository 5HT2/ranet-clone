package main

import (
	"flag"
	"log"
	"os"
	"ranet-clone/cfg"
	"ranet-clone/dl"
	"ranet-clone/ocr"
	threads2 "ranet-clone/threads"
	"sync"
)

var (
	baseUrl              = "https://russianplanes.net/images/"
	baseDir              = flag.String("dir", "", "full dir path to download to, eg /raspberry/img/")
	mode                 = flag.String("mode", "all", "func to do")
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

	cfg.LoadConfig(dir)
	go cfg.SetupConfigSaving()

	switch *mode {
	case "all":
		modeDownload(dir)
		modeCleanup(dir)
		modeOcr(dir)
	case "download":
		modeDownload(dir)
	case "ocr":
		modeOcr(dir)
	case "cleanup":
		modeCleanup(dir)
	default:
		log.Fatalln(*mode + " is not a recognized mode")
	}
}

func modeCleanup(dir string) {
	for _, f := range threads2.GetFiles(dir) {
		if f.Size() == 0 {
			log.Printf("removing 0 bytes file: %s\n", f.Name())
			err := os.Remove(dir + f.Name())
			if err != nil {
				log.Printf("error removing file: %v\n", err)
			}
		}
	}
}

func modeOcr(dir string) {
	if err := ocr.InitClient(*tessDataPrefix, *threads); err != nil {
		log.Fatalln(err)
	}

	log.Println("computing chunked paths")
	paths, err := ocr.GeneratePaths(dir)
	log.Printf("finished computing chunked paths (%v)\n", len(paths))

	if err != nil {
		log.Fatalln(err)
	}

	ocr.ProcessImages(paths, dir)
	cfg.SaveConfig()
}

func modeDownload(dir string) {
	log.Println("computing chunked paths")
	paths, err := dl.GenerateChunkedPaths(dir, *threads, minImg, maxImg)
	log.Printf("finished computing chunked paths (%v)\n", len(paths))

	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		log.Printf("starting download thread %v\n", i)
		wg.Add(1)
		go dl.DownloadFiles(&wg, i, paths[i], dir, baseUrl)
	}

	log.Println("waiting for downloads to finish")
	wg.Wait()
	log.Println("finished downloading")
	cfg.SaveConfig()
}
