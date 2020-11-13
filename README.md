# This is a fax machine

https://www.faxmachine.dev/

## Setup

### Download the code

```
git clone https://github.com/asta-li/fax-machine.git
```

### Install dependencies

This project requires:
- [npm, React](https://nodejs.org/en/)
- [Golang](https://golang.org/doc/install)
- [GCloud SDK](https://cloud.google.com/sdk/docs/install)

### Set up Google Cloud credentials

Run the following, selecting region `us-west2`:
```
gcloud init
gcloud auth application-default login
```

### Set up programmable fax credentials

Set up secure credentials in `fax.env` for local development.
Store the same credentials in `fax.yaml` for deployment.

### Set up environment variables

Load into the local shell:
```
source ./variables.env
```

## Run the code

### Run the Go server

Build the static React frontend and launch the Go server.

From project root:
```
cd client
npm run build
cd ..
go run server/*
```

### Run in React development mode

Run the client code in development mode.

From project root:
```
cd client
npm start
```

## Deploy

Deploy to the web via Google App Engine.

From project root:
```
cd client
npm run build
cd ..
gcloud app deploy
```



## Formatting
```bash
gofmt -w server/main.go
```

