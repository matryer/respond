# respond

Package respond provides responding capabilities to data services.

## Features

  * Idiomatic way of responding to data APIs using `respond.With`
  * Encoder abstraction lets you speak different formats
  * `Transform` allows you to envelope data or handle types differently

## Usage

Step 1. Create and configure a Responder

```
r := respond.New()
```

Step 2. Wrap your `http.HandlerFunc` or `http.Handler` using the responder

```
// wrap a handler
handler = r.Handler(handler)

// wrap a HandlerFunc
fn = r.HandlerFunc(fn)
```

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