# This is a fax machine

https://www.faxmachine.dev/

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

