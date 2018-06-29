package main

import (
	jaipur "./jaipur"
)

func move() {
	jaipurGame := jaipur.NewJaipur()
	haha(jaipurGame)
}

func haha(class GameClass) {
	class.Init()
}
