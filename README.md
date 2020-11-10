# This is a fax machine

It doesn't look like it but I swear it is.

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
gcloud app deploy
```

Live from San Francisco, it's a fax machine!
https://fax-machine-295219.wl.r.appspot.com/
