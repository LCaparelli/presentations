package main

import "fmt"

func main() {
	pedidosMassa := make(chan Pedido)
	bolosCrus := make(chan BoloCru)
	pm := PreparadorDeMassa{
		pedidos:   pedidosMassa,
		bolosCrus: bolosCrus,
	}
	go pm.Trabalhar()

	bolosAssados := make(chan BoloAssado)
	ab := AssadorDeBolo{
		bolosCrus:    bolosCrus,
		bolosAssados: bolosAssados,
	}
	go ab.Trabalhar()

	pedidosCobertura := make(chan Pedido)
	coberturas := make(chan Cobertura)
	pc := PreparadorDeCobertura{
		pedidos:    pedidosCobertura,
		coberturas: coberturas,
	}
	go pc.Trabalhar()

	bolosComCobertura := make(chan BoloComCobertura)
	ac := AplicadorDeCobertura{
		coberturas:        coberturas,
		bolosAssados:      bolosAssados,
		bolosComCobertura: bolosComCobertura,
	}
	go ac.Trabalhar()

	pedidos := make(chan Pedido)
	go PedidoFanOut(pedidos, pedidosMassa, pedidosCobertura)

	pedidos <- Pedido{
		pesoEmGramas: 1000,
	}
	_ = <-bolosComCobertura
	fmt.Println("Bolo com cobertura assado!")
}
