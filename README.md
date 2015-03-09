# File Dupes [![Build Status](https://travis-ci.org/reem/go-file-dupes.svg?branch=master)](https://travis-ci.org/reem/go-file-dupes)

> Get all the duplicate files in a directory.

## Example

```go
package main

import (
    dupes "github.com/reem/go-file-dupes"
    "os"
    "fmt"
)

func main() {
	first, _ := os.Open("./file")
	other, _ := os.Open("./duplicate")
	duplicates, _ := dupes([]*os.File{first, other})
	fmt.Println("Dupes: ", duplicates)
}
```

## Author

[Jonathan Reem](https://medium.com/@jreem) is the primary author and maintainer of file-dupes

## License

MIT

