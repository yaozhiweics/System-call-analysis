package main

import "syscall/app/controls"

func main() {
	controls.InitDB()
	controls.InitServer()
}
