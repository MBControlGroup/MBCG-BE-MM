FROM golang:1.8

COPY . "$GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM/"
RUN cd "$GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM" && go get -v && go install -v

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \  
    && echo 'Asia/Shanghai' >/etc/timezone   

WORKDIR $GOPATH/src/github.com/MBControlGroup/MBCG-BE-MM

CMD ["go", "run", "main.go"]
