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

### Set up Twilio credentials

Follow the [instructions](https://www.twilio.com/docs/usage/secure-credentials) to set up secure
credentials in `twilio.env` for local development.
Store the same credentials in `twilio.yaml` for deployment.
```
env_variables:
  TWILIO_ACCOUNT_SID: 'XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX'
  TWILIO_AUTH_TOKEN: 'XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX'
```

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
go run server/main.go
```

### Run in React development mode

After starting the Go server, run the following to run the client code in development mode.

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

