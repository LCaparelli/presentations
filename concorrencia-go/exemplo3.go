package main

import "fmt"

const (
	trabalhadoresBolo = 6000
	pesoBolo = 1000000
)

func main() {
	pedidosDivisor := make(chan Pedido)
	d := DivisorDePedido{
		pedidos:        pedidosDivisor,
	}

	pedacosAssados := make(chan BoloAssado)
	for i := 0; i < trabalhadoresBolo; i++ {
		pedidosMassa := make(chan Pedido)
		d.trabalhadores = append(d.trabalhadores, pedidosMassa)

		bolosCrus := make(chan BoloCru)
		pm := PreparadorDeMassa{
			pedidos:   pedidosMassa,
			bolosCrus: bolosCrus,
		}
		go pm.Trabalhar()

		ab := AssadorDeBolo{
			bolosCrus:    bolosCrus,
			bolosAssados: pedacosAssados,
		}
		go ab.Trabalhar()
	}
	go d.Trabalhar()

	bolosAssados := make(chan BoloAssado)
	j := JuntadorDeBolo{
		divisoes:     trabalhadoresBolo,
		pedacos:      pedacosAssados,
		bolosAssados: bolosAssados,
	}
	go j.Trabalhar()

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
	go PedidoFanOut(pedidos, pedidosDivisor, pedidosCobertura)

	pedidos <- Pedido{
		pesoEmGramas: pesoBolo,
	}
	_ = <-bolosComCobertura
	fmt.Println("Bolo com cobertura assado dividindo e conquistando!")
}
