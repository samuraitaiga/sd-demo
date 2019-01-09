FROM amd64/debian:7.11 as builder

RUN apt-get update && apt-get -y upgrade && apt-get -y install git curl wget gcc
RUN git clone https://github.com/syndbg/goenv.git ~/.goenv
ENV GOENV_ROOT=/root/.goenv
ENV PATH /root/.goenv/bin:/root/.goenv/shims:$PATH
RUN echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile
RUN echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
RUN echo 'eval "$(goenv init -)"' >> ~/.bash_profile
RUN eval "$(goenv init -)"
RUN goenv install 1.11.2
RUN goenv global 1.11.2
RUN go get -u cloud.google.com/go/cmd/go-cloud-debug-agent github.com/astaxie/beego github.com/astaxie/beego/context github.com/beego/bee go.opencensus.io/trace contrib.go.opencensus.io/exporter/stackdriver cloud.google.com/go/profiler github.com/sirupsen/logrus
COPY . /root/go/src/devops-handson
RUN cd /root/go/src/devops-handson && go build -gcflags=all='-N -l' -ldflags=-compressdwarf=false ./main.go

FROM amd64/debian:7.11 as production
RUN apt-get update && apt-get -y upgrade && apt-get -y install ca-certificates
COPY --from=builder /root/go/src/devops-handson /root/go/src/devops-handson
EXPOSE 8080
CMD cd /root/go/src/devops-handson && GOOGLE_APPLICATION_CREDENTIALS=./gcp-credentials/auth.json ./main
