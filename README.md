# Url Shortner

### A Simple Url transformer/Shortener . Given any Url it returns a Redirect object that contain a redirect code for the given URL

## Requirements
*** 

* ### PostMan
* ### MongoDb
* ### Redis
***
## Libraries

* ### go get gopkg.in/dealancer/validate.v2
* ### go get github.com/teris-io/shortid

* ### go get github.com/pkg/errors


## EndPoints

* ##   POST localhost:8000

>> Request Body 

```GO
    {
            "url" : "https://www.udemy.com"
    }
```

>> Response 

```GO
{
    "code": "hbS9RwPVR",
    "url": "https://www.udemy.com",
    "created_at": 1682226936
}
```

* ## GET localhost:8000/{code}

 Use the code in the response body returned from the **POST** request to make a call to the GET endpoint   **localhost:8000**

The Url redirects to the orignal Url .

Example 
 <span style="color:Aqua">localhost:8000/hbS9RwPVR</span>  redirects to ***https://www.udemy.com***