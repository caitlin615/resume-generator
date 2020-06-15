FROM golang

# Based on https://github.com/mark-adams/docker-chromium-xvfb

RUN apt-get update && apt-get install -y \
  chromium \
  xvfb \
  unzip

RUN curl -SLO "https://chromedriver.storage.googleapis.com/$(curl -s https://chromedriver.storage.googleapis.com/LATEST_RELEASE_79)/chromedriver_linux64.zip" \
  && unzip "chromedriver_linux64.zip" -d /usr/local/bin \
  && rm "chromedriver_linux64.zip"

ENV CHROME_DRIVER_BINARY_PATH /usr/local/bin/chromedriver
ENV CHROME_BINARY_PATH /usr/bin/chromium
ENV TEMPLATE_ROOT /go/src/github.com/caitlin615/resume-generator/templates

WORKDIR /go/src/github.com/caitlin615/resume-generator
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go install

WORKDIR /

ENTRYPOINT ["resume-generator"]
