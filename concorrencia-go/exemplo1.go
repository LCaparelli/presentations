package main

import "fmt"

func main() {
	pedidos := make(chan Pedido)
	bolosCrus := make(chan BoloCru)
	pm := PreparadorDeMassa{
		pedidos:   pedidos,
		bolosCrus: bolosCrus,
	}
	go pm.Trabalhar()

	bolosAssados := make(chan BoloAssado)
	ab := AssadorDeBolo{
		bolosCrus:    bolosCrus,
		bolosAssados: bolosAssados,
	}
	go ab.Trabalhar()

	pedidos <- Pedido{
		pesoEmGramas:   1000,
	}
	_ = <-bolosAssados
	fmt.Println("Bolo assado!")
}
