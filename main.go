package main

import (
	"kee"
)

func main() {
	r := kee.Default()

	r.Static("/file", ".")
	r.Run(":9999")
}
