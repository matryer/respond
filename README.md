# respond [![GoDoc](https://godoc.org/github.com/matryer/respond?status.svg)](https://godoc.org/github.com/matryer/respond)

Package respond provides low-touch idiomatic API responses for Go.

## Features

  * Idiomatic way of responding to data APIs using `respond.With`
  * Encoder abstraction lets you easily speak different formats
  * `Transform` allows you to envelope data or handle data types differently
  * `Before` and `After` function fields allow you to envelope and mutate data, set common HTTP headers, log activity etc.
  * Protected against multiple responses
  * Helpers including `repsond.With`, `respond.WithStatus`, `respond.WithRedirect*` etc.

## Usage

Once you've installed respond (see below), responding with a data payload is as simple as calling `respond.With`:

```
// respond with data
respond.With(w, r, http.StatusOK, data)

// or respond with an error
if err := doSomething(); err != nil {
  respond.With(w, r, http.StatusInternalServerError, err)
  return
}
```

Additional helpers let you easily respond with an HTTP status code:

```
respond.WithStatus(w, r, http.StatusNotFound)
return
```

Or redirects:

```
respond.WithRedirect(w, r, http.StatusTemporaryRedirect, "/new/path")
```

## Getting started with respond

Get it:

```
go get github.com/matryer/respond
```

And import it:

```
import github.com/matryer/respond
```

#### Step 1. Create and configure a Responder

```
responder := respond.New()
```

  * Generally create one per app
  * Create it at the same time you setup your server (usually in `main.go`)
  * Look at the [fields inside the Responder struct](http://godoc.org/github.com/matryer/respond#Responder) for details on how you can customize respond

#### Step 2. Wrap your `http.HandlerFunc` or `http.Handler` using the responder

```
// wrap a handler
handler = responder.Handler(handler)

// wrap a HandlerFunc
fn = responder.HandlerFunc(fn)
```

  * Wrapping the handlers with the responder allows them to use the `respond.With` function
  * Calling `respond.With` inside unwrapped handlers will cause a panic.

#### Step 3. Use `respond.With` in your handlers

```
func handleSomething(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"something": true,
		"probably-from": "database",
	}

	respond.With(w, r, http.StatusOK, data)

}
```

## Panics

#### `respond: multiple responses`

Most requests need only one response, but a common bug in code is to call `respond.With` many times. Consider this example:

```
func handleSomething(w http.ResponseWriter, r *http.Request) {
  
  data, err := LoadDataFromDatabase()
  if err != nil {
    respond.With(w, r, http.StatusInternalServerError, err)
    // NOTE: missing return here
  }
  respond.With(w, r, http.StatusOK, data)

}
```

After the error case, the code should `return` to prevent future code from running.

The solution is to make sure `respond.With` is called only once per request.

  * Advanced: You can prevent this panic by setting `respond.ManyResponsesPanic = false`

#### `respond: must wrap with Handler or HandlerFunc`

In order to use `respond.With` you must wrap http.Handler and http.HandlerFunc objects. Wrapping the handlers allows respond to setup, and teardown things it needs in order to provide responding capabilities. It is also how `respond.With` knows which `respond.Responder` to use.

This:

```
http.HandleFunc("/hello", HelloServer)
```

Becomes:

```
responder := respond.New()
http.HandleFunc("/hello", responder.HandlerFunc(HelloServer))
```

or this:

```
handler := &MyHandler{}
http.ListenAndServe(":8080", handler)
```

becomes:

```
responder := respond.New()
handler := &MyHandler{}
http.ListenAndServe(":8080", responder.Handler(handler))
```
