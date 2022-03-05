FROM golang:latest

RUN mkdir /ranet-data
ADD . /ranet-clone
WORKDIR /ranet-clone
ENV MODE all
ENV THREADS 4
RUN go build -o ranet .
CMD /ranet-clone/ranet -mode "$MODE" -dir /ranet-data -threads "$THREADS"
