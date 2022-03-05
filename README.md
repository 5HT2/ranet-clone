# ranet-clone

A searchable clone of russianplanes.net, for transparency and ease of identifying planes.

## What is this?

This is a tool for
1. Downloading the entirety of the images hosted on russianplanes.net
2. Re-hosting them and allowing people to search them
3. Re-creating the metadata that was once attached to these images via OCR

## Why?

The website russianplanes.net was told to take down all their military aircraft listings by the Russian government.
This project aims to archive all the images hosted on their CDN in order to make identification of aircraft easier.

## Usage

```bash
git clone https://github.com/l1ving/ranet-clone
cd ranet-clone

echo "{}" > config.json
go build -o ranet .

# Make the dir first
./ranet -dir /path/to/images/dir -threads 4
```

## TODO

- [x] Async downloading
- [ ] Distributed hosting
- [ ] Searching
- [x] OCR
