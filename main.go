package main

import (
	"flag"
	"log"
	"ranet-clone/cfg"
	"ranet-clone/dl"
)

var (
	baseURL = "https://russianplanes.net/images/"
	baseDir = flag.String("dir", "", "full dir path to download to, eg /raspberry/img/")
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

	paths, err := dl.GeneratePaths(100000, 100002)
	if err == nil {
		for _, p := range paths {
			log.Println("downloading " + dir + p.Path)
			dl.DownloadFile(p, dir, baseURL)
		}
	} else {
		log.Println(err.Error())
	}
}
