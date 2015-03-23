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
