# respond

Package respond provides low-touch idiomatic API responses for Go.

## Features

  * Idiomatic way of responding to data APIs using `respond.With`
  * Encoder abstraction lets you easily speak different formats
  * `Transform` allows you to envelope data or handle data types differently
  * `Before` and `After` function fields allow you to envelope and mutate data, set common HTTP headers, log activity etc.
  * Protected against multiple responses

## Usage

Step 1. Create and configure a Responder

```
responder := respond.New()
```

  * Generally create one per app
  * Create it at the same time you setup your server (usually in `main.go`)

Step 2. Wrap your `http.HandlerFunc` or `http.Handler` using the responder

```
// wrap a handler
handler = responder.Handler(handler)

// wrap a HandlerFunc
fn = responder.HandlerFunc(fn)
```

  * Wrapping the handlers with the responder allows them to use the `respond.With` function

Step 3. Use `respond.With` in your handlers

```
func handleSomething(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"something": true,
		"probably-from": "database",
	}

	respond.With(w, r, http.StatusOK, data)

}
```