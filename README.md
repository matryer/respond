# respond [![GoDoc](https://godoc.org/github.com/matryer/respond?status.svg)](https://godoc.org/github.com/matryer/respond)

Package respond provides low-touch idiomatic API responses for Go.

## Features

  * Idiomatic way of responding to data APIs using `respond.With`
  * Encoder abstraction lets you easily speak different formats
  * `Before` and `After` function fields allow you to envelope and mutate data, set common HTTP headers, log activity etc.
  * Protected against multiple responses
