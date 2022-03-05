package cfg

import "log"

func AddCompletedDownload(i ImageInfo) {
	config.run(func(c *Config) {
		c.Downloaded = append(c.Downloaded, i)
	})
}

func UpdateNumDownloaded(num int64, set bool) {
	config.run(func(c *Config) {
		if set {
			c.NumDone = num
		} else {
			c.NumDone += num
		}
	})
}

func UpdateOcrData(i ImageInfo, str string) {
	config.run(func(c *Config) {
		for _, n := range c.Downloaded {
			if n.Name == i.Name {
				log.Printf("processed %s (%s)\n", i.Name, str)
				n.OcrData = str
				break
			}
		}
	})
}

func InDlQueue(i ImageInfo) bool {
	inQueue := false
	config.run(func(c *Config) {
		inQueue = contains(c.DlQueue, i.Name)
	})
	return inQueue
}

func AddToDlQueue(i ImageInfo) {
	config.run(func(c *Config) {
		c.DlQueue = append(c.DlQueue, i)
	})
}

func RemoveFromDlQueue(p ImageInfo) {
	config.run(func(c *Config) {
		for i, q := range c.DlQueue {
			if q.Name == p.Name {
				c.DlQueue[i] = c.DlQueue[len(c.DlQueue)-1]
				c.DlQueue = c.DlQueue[:len(c.DlQueue)-1]
				break
			}
		}
	})
}

func InOcrQueue(i ImageInfo) bool {
	inQueue := false
	config.run(func(c *Config) {
		inQueue = contains(c.OcrQueue, i.Name)
	})
	return inQueue
}

func AddToOcrQueue(i ImageInfo) {
	config.run(func(c *Config) {
		c.OcrQueue = append(c.OcrQueue, i)
	})
}

func RemoveFromOcrQueue(p ImageInfo) {
	config.run(func(c *Config) {
		for i, q := range c.OcrQueue {
			if q.Name == p.Name {
				c.OcrQueue[i] = c.OcrQueue[len(c.OcrQueue)-1]
				c.OcrQueue = c.OcrQueue[:len(c.OcrQueue)-1]
				break
			}
		}
	})
}

func contains(s []ImageInfo, e string) bool {
	for _, a := range s {
		if a.Name == e {
			return true
		}
	}
	return false
}
