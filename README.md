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

Additonally, store the private key for service account `telnyx@fax-machine-295219.iam.gserviceaccount.com`
in the file `gcs_credentials.json`.
See the [Creating service account keys](https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating_service_account_keys) for how to create private keys.
These credentials are used for [signing URLs](https://cloud.google.com/storage/docs/access-control/signing-urls-manually).

### Set up programmable fax credentials

Set up secure credentials in `fax.env` for local development.
Store the same credentials in `fax.yaml` for deployment.

### Set up paypal credentials

Set up secure credentials in `payment.env` for local development, use sandbox credentials here
Store the same credentials in `payment.yaml` for deployment. use prod credentials here (if deploying to prod server)


### Set up environment variables

Load into the local shell:
```
source ./variables.env
```

## Run the code

### Run both react and go server
```bash
source ./variables.env && npm run build --prefix ./client && go run ./server
```

### Run the Go server

Build the static React frontend and launch the Go server.

From project root:
```
source ./variables.env
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

Golang formatter
```bash
gofmt -w server/*.go
```

Javascript formatter
```bash
cd client
prettier src/*.js --write
```

Javascript linter
```bash
eslint src/*.js
```

## Test the API locally

### Test `/api/fax-status`

Query the fax status endpoint for the given Fax ID.
```curl -X GET http://localhost:3000/api/fax-status?id=${FAX_ID}

```

606bcbe5-a33b-4747-aa05-e691fa09710d
### Test `/fax-webhook`

This webhook handles fax status updates from Telnyx:

See https://developers.telnyx.com/docs/api/v2/programmable-fax/Programmable-Fax-Commands
```
export TEST_RESPONSE='{ "data": { "event_type": "fax.queued", "id": "3691d047-d22a-424d-80ed-fe871981aa6d", "occurred_at": "2020-04-22T19:32:12.538002Z", "record_type": "event", "payload": { "connection_id": "7267xxxxxxxxxxxxxx", "fax_id": "b679398e-8b4c-46bd-8630-6797f1ab5228", "from": "+35319605860", "original_media_url": "http://www.telnyx.com/telnyx-fax/1.pdf", "status": "queued", "to": "+13129457420", "user_id": "a5b1dfa3-fd2e-4e02-8ea4-xxxxxxxxxxxx" } }, "meta": { "attempt": 1, "delivered_to": "http://example.com/webhooks" } }'
curl -X POST -H "Content-Type: application/json"  -d "${TEST_RESPONSE}"  http://localhost:3000/fax-complete
```
