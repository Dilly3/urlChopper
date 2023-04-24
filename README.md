# Url Shortener

### A Simple Url transformer/Shortener . Given any Url it returns a Redirect object that contain a redirect code for the given URL

## Requirements
*** 

* ### PostMan
* ### MongoDb
* ### Postgres
* ### Redis
* ### Chi Router
***
## Libraries

* ### go get gopkg.in/dealancer/validate.v2
* ### go get github.com/teris-io/shortid

* ### go get github.com/pkg/errors

* ### go get -u github.com/go-chi/chi/v5


## EndPoints

* ###  POST localhost:8000
* ### GET localhost:8000/{code}

<br>

## POST 
* ### localhost:8000

>> Request Body 

```GO
     {
        "url" : "https://www.udemy.com",
        "name" : "Udemy"
    }
```

>> Response 

```GO
{
    "name": "Udemy",
    "code": "910HGwEVR",
    "url": "https://www.udemy.com",
    "created_at": 1682230484
}
```
## GET
*  ### <span style="color:Orange">localhost:8000/{code}</span>

 Use the code in the response body returned from the **POST** request to make a call to the GET endpoint   **localhost:8000**

The Url redirects to the orignal Url .

Example 
 <span style="color:Aqua">localhost:8000/910HGwEVR</span>  redirects to ***https://www.udemy.com***

 <br>

 ## Start Program

### setup .env file
 

 ```GO

 PORT= 8000 || {any choice of port}
 choose db 
URL_DB= { any of mongo || postgres } // URL_DB=mongo or URL_DB=postgres
MONGO_URL={mongodb url}
MONGO_TIMEOUT=30

POSTGRES_PORT=****
POSTGRES_DBNAME=Redirects // or any Database name of choice
POSTGRES_PASSWORD=****
POSTGRES_USER=****

```

Run Program

```GO
terminal> 
>> go mod tidy

>> make run
```
