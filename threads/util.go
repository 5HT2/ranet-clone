package threads

import (
	"log"
	"os"
	"ranet-clone/cfg"
	"runtime/debug"
)

func LogPanic() {
	if x := recover(); x != nil {
		// recovering from a panic; x contains whatever was passed to panic()
		log.Printf("panic: %s\n", debug.Stack())
	}
}

func GetFiles(dir string) []os.FileInfo {
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

func ChunkSlice(slice []cfg.ImageInfo, chunkSize int) [][]cfg.ImageInfo {
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
