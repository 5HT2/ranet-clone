FROM golang:latest


RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
RUN apt-get install -y -qq tesseract-ocr-eng

RUN mkdir /ranet-data
ADD . /ranet-clone
WORKDIR /ranet-clone
ENV MODE all
ENV THREADS 4
RUN go build -o ranet .
CMD /ranet-clone/ranet -mode "$MODE" -dir /ranet-data -threads "$THREADS" -tessdata "$TESSDATA_PREFIX"
