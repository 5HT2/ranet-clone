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

# Make the dir first
RANET_DATA=/path/to/images/dir
echo "{}" > "$RANET_DATA/config.json"

#
# Run directly
go build -o ranet .
./ranet -dir $RANET_DATA -threads 4

#
# Or, run via Docker
docker build -t ranet .
docker run --name ranet --mount type=bind,source="$RANET_DATA",target=/ranet-data --network host -d -e MODE=all -e THREADS=4 ranet
```

## TODO

- [x] Async downloading
- [ ] Distributed hosting
- [ ] Searching
- [x] OCR
