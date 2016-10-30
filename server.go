package main

type Server struct {
	// Host, Server, Type, Status, Boot string
	Cores  string `json"smp:cores"`
	Ide00  string `json"ide:0:0"`
	Ide01  string `json"ide:0:1"`
	Ide10  string `json"ide:1:0"`
	Ide11  string `json"ide:1:1"`
	Block0 string `json"block:0"`
	Block1 string `json"block:1"`
	Block2 string `json"block:2"`
	Block3 string `json"block:3"`
}
