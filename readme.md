# Sicobo - Simple personal comic book collection management

Made for personal use - simple, light, open solution. Books are added by ISBN code, and data are identified in Azure bing cognitive services and google books public APIs.

**THIS IS A WORK IN PROGRESS**

Even the license may be modified, but for now you can already use is as a simple starter project. Take care of database switching : for now it's not threadsafe, it will be updated.

**Updated 2019-11-17** :
* Now use Go Modules
* Now use Gin https://github.com/gin-gonic/gin 
* Now integrated swagger gen with https://github.com/swaggo/gin-swagger

**Some tests ISBN :** (various french comics and books)

- 9782413008033
- 9780578175416
- 9782723437127
- 9782370731258
- 9782264039323
- 9782344021750
- 9782918645290
- 9782377540037

## Tech

### Dev

Go Modules

### Swagger spec

Update spec on API code and then run :
```
swag init
```

### For library update from ISBN code

Uses google book, fallover on bing search for unknown entries. Gather snapshot pictures from bing image search. All these services are **very** limited in amount of queries by day / by month. Don't expect to query them more than 1000x every month

### Database

Uses a mongoDb instance. One database = one book collection

Start a mongo instance with docker (for test only, no security) : 

```
docker run --name mongo -d -p 27017:27017 mongo:4.2
```

## Install sicobo

### Clone

Clone sources, then :

```
go mod tidy
```

### Edit config

Edit file `config.json` :

```
{
    "bingAPIkey": "you azure bing / cognitive-services api key",
    "googleBookAPIKey": "your google book api key",
    "logLevel": "INFO"{
    "mongoDb": "your mongodb instance , for example mongodb://my-server:27017"
}
```

### Start

(for now, dev is in progress)

```
go run biblio.go
```

## Use

Access http://127.0.0.1:8080/apidocs/?url=http://127.0.0.1:8080/apidocs.json and test with swagger-ui. A Angular frontend will be built later

## License

Do What The Fuck You Want To Public License (WTFPL) - http://www.wtfpl.net/about/
