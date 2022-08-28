// TODO: golint had WAY stricter checks, but it's deprecated...
package main

import "fmt"

type Failure bool

func Main() {
// TODO: e.g. right now this is the only check I can get to fail
if true == true {
fmt.Println("This file should FAIL gofmt & golangci-lint checks")
}
}
