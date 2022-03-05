package cfg

func AddCompletedDownload(imagePath ImageInfo) {
	config.run(func(c *Config) {
		c.Downloaded = append(c.Downloaded, imagePath)
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

func InQueue(i ImageInfo) bool {
	inQueue := false
	config.run(func(c *Config) {
		inQueue = contains(c.Queue, i.Name)
	})
	return inQueue
}

func AddToQueue(i ImageInfo) {
	config.run(func(c *Config) {
		c.Queue = append(c.Queue, i)
	})
}

func RemoveFromQueue(p ImageInfo) {
	config.run(func(c *Config) {
		for i, q := range c.Queue {
			if q.Name == p.Name {
				c.Queue[i] = c.Queue[len(c.Queue)-1]
				c.Queue = c.Queue[:len(c.Queue)-1]
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
