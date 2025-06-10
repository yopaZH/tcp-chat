package main

func main() {
	c, err := client.NewClient("localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	go c.HandleIncoming()
	c.HandleOutgoing()
}
