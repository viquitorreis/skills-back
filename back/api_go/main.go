package main

func main() {
	server := NewApiServer(":3030")

	server.Run()
}

