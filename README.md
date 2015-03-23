# respond
API responding for Go

### Examples

Simple responding:

```
respond.With(http.StatusOK, data).To(w, r)
```

```
respond.With(http.StatusInternalServerError, err).To(w, r)
```

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
