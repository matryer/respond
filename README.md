# respond
API responding for Go

### Examples

Simple responding:

```
respond.With(http.StatusOK, data).To(w, r)
```

Or respond with errors:

```
respond.With(http.StatusInternalServerError, err).To(w, r)
```

## Advanced

  * Check out the [documentation](http://godoc.org/github.com/matryer/respond)

#### Public view

To control how your own types are exposed through repsond, they can implement the `Public` interface by providing a simple `Public() interface{}` method:

```
type Auth struct {
	Secret string
	Key string
}

func (a *Auth) Public() interface{} {
	return map[string]interface{}{"Key": a.Key}
}
```

  * When you respond with an object of type `Auth`, the map returned by the `Public()` method will be written to the response instead.
  * Returning another object that implements `Public` is OK up to a point. 

#### Headers

`respond.Headers()` allows you to setup headers for every response:

```
respond.Headers().Set("X-App-Version", "1.0")
```

Or use the `AddHeader`, `SetHeader` and `DelHeader` fluent functions on `With`:

```
respond.With(http.StatusOK, data).
	DelHeader("X-Global").
	AddHeader("X-RateLimit", rateLimitVal).
	SetHeader("X-Log", "Some item").
	To(w, r)
```

#### Transforming

The data can be transformed by setting a `TransformFunc` to use to modify any results before they are written.

By default, `error` types are wrapped into a map, and the `Error() string` is used as the value:

```
{"error":"error message"}
```

To add your own transformations, call `respond.Transform`:

```
respond.Transform(func(r *http.Request, data interface{}) interface{} {
	switch o := data.(type) {
		case error:
			return map[string]interface{}{"error": o.Error()}
		case YourType:
			return map[string]interface{}{"name": o.Name()}
	}
	return data // by default, don't transform it
})
```

#### Encoders

By default, repsond speaks JSON. But you can add other encoders:

```
respond.Encoders().Add("xml", YourXMLEncoder())
```

And you can stop respond from speaking JSON like this:

```
respond.Encoders().Del(respond.JSON)
```

By default, `respond.DefaultEncoder` will be used if no others match.
