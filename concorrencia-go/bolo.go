package main

import (
	"time"
)

type Pedido struct {
	pesoEmGramas int
}
type BoloCru struct {
	pesoEmGramas int
}

type PreparadorDeMassa struct {
	pedidos   <-chan Pedido
	bolosCrus chan<- BoloCru
}

func (pm PreparadorDeMassa) Trabalhar() {
	for pedido := range pm.pedidos {
		pm.bolosCrus <- pm.PreparaBolo(pedido)
	}
}

func (pm PreparadorDeMassa) PreparaBolo(p Pedido) BoloCru {
	time.Sleep(time.Millisecond * 6 * time.Duration(p.pesoEmGramas)) // simulando trabalho
	return BoloCru{pesoEmGramas: p.pesoEmGramas}
}

type BoloAssado struct{}

type AssadorDeBolo struct {
	bolosCrus    <-chan BoloCru
	bolosAssados chan<- BoloAssado
}

func (ab AssadorDeBolo) Trabalhar() {
	for boloCru := range ab.bolosCrus {
		ab.bolosAssados <- ab.AssaBolo(boloCru)
	}
}

func (ab AssadorDeBolo) AssaBolo(b BoloCru) BoloAssado {
	time.Sleep(time.Millisecond * 6 * time.Duration(b.pesoEmGramas)) // simulando trabalho
	return BoloAssado{}
}

type Cobertura struct{}

type PreparadorDeCobertura struct {
	pedidos    <-chan Pedido
	coberturas chan<- Cobertura
}

func (pc PreparadorDeCobertura) Trabalhar() {
	for pedido := range pc.pedidos {
		pc.coberturas <- pc.PreparaCobertura(pedido)
	}
}

func (pc PreparadorDeCobertura) PreparaCobertura(p Pedido) Cobertura {
	time.Sleep(time.Second * 2) // simulando trabalho
	return Cobertura{}
}

type BoloComCobertura struct{}

type AplicadorDeCobertura struct {
	coberturas        <-chan Cobertura
	bolosAssados      <-chan BoloAssado
	bolosComCobertura chan<- BoloComCobertura
}

func (ac AplicadorDeCobertura) Trabalhar() {
	for cobertura := range ac.coberturas {
		boloAssado, ok := <-ac.bolosAssados
		if !ok {
			return
		}
		ac.bolosComCobertura <- ac.AplicaCobertura(cobertura, boloAssado)
	}
}

func (ac AplicadorDeCobertura) AplicaCobertura(cobertura Cobertura, boloAssado BoloAssado) BoloComCobertura {
	time.Sleep(time.Second * 2) // simulando trabalho
	return BoloComCobertura{}
}

func PedidoFanOut(fonte chan Pedido, destinos ...chan<- Pedido) {
	for p := range fonte {
		for _, dst := range destinos {
			dst <- p
		}
	}
}

type DivisorDePedido struct {
	pedidos       <-chan Pedido
	trabalhadores  []chan<- Pedido
}

func (d DivisorDePedido) Trabalhar() {
	for pedido := range d.pedidos {
		porcaoBolo := pedido.pesoEmGramas / len(d.trabalhadores)
		d.EnviarPedidoPorcao(Pedido{pesoEmGramas: porcaoBolo})
	}
}

func (d DivisorDePedido) EnviarPedidoPorcao(p Pedido) {
	entradaTrabalhadores := make(chan Pedido)
	go PedidoFanOut(entradaTrabalhadores, d.trabalhadores...)
	entradaTrabalhadores <- p
	close(entradaTrabalhadores)
}

type JuntadorDeBolo struct {
	divisoes     int
	pedacos      <-chan BoloAssado
	bolosAssados chan<- BoloAssado
}

func (j JuntadorDeBolo) Trabalhar() {
	i := 1
	for _ = range j.pedacos {
		if i < j.divisoes {
			i++
			continue
		}
		j.bolosAssados <- BoloAssado{}
		i = 1
	}
}
