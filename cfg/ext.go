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
