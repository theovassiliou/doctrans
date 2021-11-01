package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string `opts:"env=FOO"`
	Bar string `opts:"env"`
	Joe string `opts:"env"`
}

func main() {
	p1 := Config{}

	opts.New(&p1).
		UserConfigPath().
		ConfigPath("config.json").
		Repo("github.com/theovassiliou/doctrans").
		Parse()
	fmt.Println(p1.Foo)
	fmt.Println(p1.Bar)
	fmt.Println(p1.Joe)
}
