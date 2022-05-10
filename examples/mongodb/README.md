# MongoDB

this example shows how to use the CHC basic features including, routes, controllers, cookies, and redirects.

## Usage

    $ go mod init example/mongodb
    $ go get github.com/cookie-for-pres/chc
    $ go get go.mongodb.org/mongo-driver/mongo
    $ go run main.go

remove `.example` from `.env.example` and add your mongodb uri after the `=` sign.
then you can go to http://localhost:8080 to see all your databases.
