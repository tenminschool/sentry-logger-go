# Sentry Error Logger for Go

This repository provides a simple way to integrate Sentry error logging into your Go application using the `sentry-logger-go` package.

## Installation

To use `sentry-logger-go`, you need to download it first. You can do this by running:

```bash
go get github.com/tenminschool/sentry-logger-go
```

## Usages

Add a router middleware to capture all api's error

```bash
artifact.Router.Use(sentryLoggerGo.SentryMiddleware)
```

Sentry init

```bash
sentryLoggerGo.SentryInit(Config.GetString("App.SentryDsn"))
```

