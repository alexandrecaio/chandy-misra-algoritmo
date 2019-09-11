package main

import (
	"time"
	"fmt"
)

type Message struct {
	node string
	dist int
	peso int
}

type Neighbour struct {
	Id     string
	From   chan Message
	To     chan Message
	Weight int
}

func sendMessage(ch chan Message, pid string, dist int, peso int) {
	ch <- Message{pid, dist, peso}
	fmt.Printf("Message sent from pid=%s. \n", pid)

}

func receiveMessage(ch chan Message) int {
	message := <-ch

	return message.dist + message.peso
}

func redirect(in chan Message, neigh Neighbour) {
	for{
		token := <-neigh.From
		in <- token
	}
	
}

func process(init int, pid string, neighs ...Neighbour) {
	nmap := make(map[string]Neighbour)
	dist := 999
	//	pai := -1
	//	sendMessage(ch13, "p", d, 3)

	if init == 1 {
		dist = 0

		for _, neigh := range neighs {
			nmap[neigh.Id] = neigh
			sendMessage(neigh.To, pid, dist, neigh.Weight)
		}
	}

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Message, 1)
	
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(in, neigh)
	}


	for{

		msg := <-in
		if msg.dist+msg.peso < dist {
			dist = msg.dist + msg.peso
			fmt.Printf("Msg received from d(%s)=%v  W(%s-%s)=%v \n",msg.node,msg.dist,msg.node,pid,msg.peso)
			for _, neigh := range neighs {
				nmap[neigh.Id] = neigh
				sendMessage(neigh.To, pid, dist, neigh.Weight)
			}
		}
		//time.Sleep(5 * time.Second)
		fmt.Printf("Dist to root from %s: %v\n",pid,dist)

	}

}

func main() {

	fmt.Printf("Running Chandy-Misra \n")

	//Mapeamento dos Canais do Grafo
	// p --- q
	oneTwo := make(chan Message, 100) // ch 12 (p->q)
	twoOne := make(chan Message, 100) // ch 21 (q->p)

	// p --- r
	oneThree := make(chan Message, 100) // ch 13 (p->r)
	threeOne := make(chan Message, 100) // ch 31 (r->p)

	// q --- r
	twoThree := make(chan Message, 100) // ch 23 (q->r)
	threeTwo := make(chan Message, 100) // ch 32 (r->q)

	// r --- t
	threeFour := make(chan Message, 100) // ch 34 (r->t)
	fourThree := make(chan Message, 100) // ch 43 (t->r)

	// r --- s
	threeFive := make(chan Message, 100) // ch 35 (r->s)
	fiveThree := make(chan Message, 100) // ch 53 (s->r)

	// t --- s
	fourFive := make(chan Message, 100) // ch 45 (t->s)
	fiveFour := make(chan Message, 100) // ch 54 (s->t)

	go process(1, "p", Neighbour{"q", oneTwo, twoOne, 1}, Neighbour{"r", oneThree, threeOne, 3})
	go process(0, "q", Neighbour{"p", twoOne, oneTwo, 1}, Neighbour{"r", twoThree, threeTwo, 1})
	go process(0, "r", Neighbour{"q", threeTwo, twoThree, 1}, Neighbour{"p", threeOne, oneThree, 3},Neighbour{"t", threeFour, fourThree, 4}, Neighbour{"s", threeFive, fiveThree, 1})
	go process(0, "t", Neighbour{"r", threeFour, fourThree, 4}, Neighbour{"s", fourFive, fiveFour, 1})
	go process(0, "s", Neighbour{"r", fiveThree, threeFive, 1}, Neighbour{"t", fiveFour, fourFive, 1})
	
	
	//processTwo(twoOne, oneTwo, twoThree, threeTwo)
	//go processThree(threeOne, oneThree, threeTwo, twoThree, threeFour, fourThree, threeFive, fiveThree)
	//go processFour(fourThree, threeFour, fourFive, fiveFour)
	//processFive(fiveThree, threeFive, fiveFour, fourFive)
	time.Sleep(10 * time.Second)
}
