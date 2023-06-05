# Coding Exercise: Data Storage API

Implement a small HTTP service to store objects organized by repository.
Clients of this service should be able to GET, PUT, and DELETE objects.

## General Requirements

* The service should de-duplicate data objects by repository.
* The service should listen on port `8282`.
* The included tests should pass and not be modified. Adding additional tests is encouraged.
* The service must implement the API as described below.
* The data can be persisted in memory, on disk, or wherever you like.
* Do not include any extra dependencies. Anything in the Go standard library is fine.

## Suggestions

* Your code will be read by humans, so organize it sensibly.
* Use this repository to store your work. Committing just the final solution is *ok* but we'd love to see your incremental progress.
* Treat this pull request as if you’re at work submitting it to your colleagues, or to an open source project. The body of the pull request can be used to describe your reasoning and any assumptions, limitations or tradeoffs in your implementation, or anything you're really proud of in your submission 😄.
* Remember that this is a web application and concurrent requests could come in. If you have time, this is a good problem to address.

## API

### Upload an Object

```
PUT /data/{repository}
```

#### Response

```
Status: 201 Created
{
  "oid": "2845f5a412dbdfacf95193f296dd0f5b2a16920da5a7ffa4c5832f223b03de96",
  "size": 1234
}
```

### Download an Object

```
GET /data/{repository}/{objectID}
```

#### Response

```
Status: 200 OK
{object data}
```

Objects that are not on the server will return a `404 Not Found`.

### Delete an Object

```
DELETE /data/{repository}/{objectID}
```

#### Response

```
Status: 200 OK
```

## Getting started and Testing

In `testapp.go` you'll find a naive first draft of the answer to the exercise written for you. Please improve this draft so that it passes the test written in `coding_test.go`. You can test this by running:

```
go test -count=1
```

Or simply:

```
make
```