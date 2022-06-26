# Media Storage Api Service

## Overview

This service is used for handle all endpoint and data about **Media/Storage**. [Golang](https://golang.org/) is the main weapon of this service. This service using [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). We are implement [this go clean architecture](https://github.com/bxcodec/go-clean-arch). Please read their article first for explanation of this architecture.

## How to run

1. Clone it
1. Copy paste `.env.example` and rename it into `.env`
1. Adjust the config in your `.env`
1. Run `make run` or manually `go run app/main.go`. It will download all dependencies and running your application
