# respond [![GoDoc](https://godoc.org/github.com/matryer/respond?status.svg)](https://godoc.org/github.com/matryer/respond)

```
respond.With(w, r, http.StatusOK, data)
```

Package respond provides low-touch API responses for Go data services.

  * Idiomatic way of responding to data APIs using `respond.With`
  * Use `respond.With` to respond with default options, or make a `respond.Options` for [advanced features](https://godoc.org/github.com/matryer/respond#Options)
  * Encoder abstraction lets you easily speak different formats
  * `Before` and `After` function fields allow you to envelope and mutate data, set common HTTP headers, log activity etc.
  * Protected against multiple responses

## Usage

The simplest use of `respond` is to just call `respond.With` inside your handlers:

```
func handleSomething(w http.ResponseWriter, r *http.Request) {
	
	data, err := loadFromDB()
	if err != nil {

		// respond with an error
		respond.With(w, r, http.StatusInternalServerError, err)
		return // always return after responding

	}

	// respond with OK, and the data
	respond.With(w, r, http.StatusOK, data)

}
```

To tweak the behaviour of `respond.With` you can wrap the handler with a `respond.Options`:

```
func main() {
	
	opts := &respond.Options{
		// options go here
	}
	http.Handle("/foo", opts.Handler(fooHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
```

Use `respond.With` as normal.

For a complete list of options, see the [API documentation for respond.Options](https://godoc.org/github.com/matryer/respond#Options).