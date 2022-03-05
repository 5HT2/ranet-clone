package main

import (
	"flag"
	"fmt"
	"log"
	"ranet-clone/cfg"
	"ranet-clone/dl"
	"sync"
)

var (
	baseUrl       = "https://russianplanes.net/images/"
	baseDir       = flag.String("dir", "", "full dir path to download to, eg /raspberry/img/")
	mode          = flag.String("mode", "download", "func to do")
	threads       = flag.Int("threads", 4, "threads to use for downloading")
	minImg  int64 = 100000
	maxImg  int64 = 301605
)

func main() {
	flag.Parse()

	dir := *baseDir
	if len(dir) == 0 {
		log.Fatalln("*baseDir is empty")
	}

	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	cfg.LoadConfig()
	go cfg.SetupConfigSaving()

	switch *mode {
	case "download":
		modeDownload(dir)
	default:
		log.Fatalln(*mode + " is not a recognized mode")
	}
}

func modeDownload(dir string) {
	log.Println("computing chunked paths")
	paths, err := dl.GenerateChunkedPaths(dir, *threads, minImg, maxImg)
	if err == nil {
		var wg sync.WaitGroup

		for i := 0; i < *threads; i++ {
			fmt.Println("starting download thread " + string(rune(i)))
			wg.Add(1)
			go dl.DownloadFiles(&wg, paths[i], dir, baseUrl)
		}

		fmt.Println("waiting for downloads to finish")
		wg.Wait()
		fmt.Println("finished downloading")
	} else {
		log.Fatalln(err.Error())
	}
}
