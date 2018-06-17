FROM golang:1.8

COPY . "$GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM/"
RUN cd "$GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM" && go get -v && go install -v

WORKDIR $GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM

CMD ["go", "run", "main.go"]
