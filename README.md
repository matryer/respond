# respond
API responding for Go

### Examples

Simple responding:

```
respond.With{http.StatusOK, dataMap}.To(w, r)
```

```
respond.With{http.StatusInternalServerError, err}.To(w, r)
```

#### Global headers

```
respond.Headers.Set("X-App-Version", "1.0")

respond.With()
```