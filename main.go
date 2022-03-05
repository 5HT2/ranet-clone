package main

import (
	"flag"
	"log"
	"ranet-clone/cfg"
	"ranet-clone/dl"
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
	paths, err := dl.GeneratePaths(dir, minImg, maxImg)
	if err == nil {
		for _, p := range paths {
			log.Println("downloading " + dir + p.Path)
			dl.DownloadFile(p, dir, baseUrl)
		}
	} else {
		log.Fatalln(err.Error())
	}
}
