package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)
// Struct Player 
type Player struct{
	nome string
	pontos int
	games int
	sets int
	estado string
	venceuGame bool
	venceuSets bool
	venceuJogo bool
}

var PONTOSPARAVITORIA = 4
var GAMESPARAVITORIA = 6
var SETSPARAVITORIA = 3
var RODADA = 0

//JogarGame simula a execucao de um Game
func JogarGame(player *Player, playerTwo *Player, waitGroup *sync.WaitGroup, mutexObject *sync.Mutex){
	
	for {

		if player.estado == "rebatendo"{
			
			/*
				Gerando numero aleatorio
			*/
			seed := rand.NewSource(time.Now().UnixNano())
			random := rand.New(seed)
			r := random.Intn(100)
		
			mutexObject.Lock() 
			
			/*
				Verificando se algum dos player possui 
				os requisitos para vencer o game
			*/
			if (player.pontos >= PONTOSPARAVITORIA || playerTwo.pontos >= PONTOSPARAVITORIA) && (player.pontos - playerTwo.pontos >= 2 || playerTwo.pontos - player.pontos >= 2) {
				playerTwo.estado = "rebatendo"
				if player.pontos > playerTwo.pontos{
					player.venceuGame = true
				}
				break
			}

			fmt.Println(player.nome, "ESPERANDO RECEBER A BOLA")
			
			/*
				Verificando se o player conseguiu 
				ou nao rebater a bola
			*/
			if r % 2 == 0{
				fmt.Println(player.nome, "REBATEU A BOLA")
			} else {
				fmt.Println(player.nome, "NÃO REBATEU A BOLA, PONTO PARA O ", playerTwo.nome)
				playerTwo.pontos = playerTwo.pontos + 1
				
				//funcao para imprimir o placar das rodadas
				printPlacarRodada(player, playerTwo)
			}
			
			/*
				mudando os estados dos players
			*/
			player.estado = "esperando"
			playerTwo.estado = "rebatendo"
			
			mutexObject.Unlock()
		}
	}

	/*
		mutex.Unlock para liberar o uso para 
		o outro player apos break
	*/
	mutexObject.Unlock()
	waitGroup.Done()
}

//JogarSet simula a execucao de um Set
func JogarSet(playerOne *Player, playerTwo *Player, waitGroup *sync.WaitGroup, mutexObject *sync.Mutex){
	
	for{
		/*
			Verificando se algum dos player possui 
			os requisitos para vencer o set
		*/
		if (playerOne.games >= GAMESPARAVITORIA || playerTwo.games >= GAMESPARAVITORIA) && (playerOne.games - playerTwo.games >= 2 || playerTwo.games - playerOne.games >= 2){
			if playerOne.games >= GAMESPARAVITORIA{
				playerOne.venceuSets = true
			} else{
				playerTwo.venceuSets = true
			}
			break
		}
		
		waitGroup.Add(2)

		/*resetando atributos dos players*/
		playerOne.pontos = 0
		playerOne.estado = "rebatendo"
		playerOne.venceuGame = false
		playerTwo.pontos = 0
		playerTwo.estado = "esperando"
		playerTwo.venceuGame = false

		/*iniciando goroutines*/
		go JogarGame(playerOne, playerTwo, waitGroup, mutexObject)
		go JogarGame(playerTwo, playerOne, waitGroup, mutexObject)
		
		/* esperando as goroutines finalizarem*/
		waitGroup.Wait()
		
		//funcao para imprimir o resultado do game
		printResultadoGame(playerOne, playerTwo)

		RODADA = 0

		//funcao para imprimir o placar dos games
		printPlacarGame(playerOne, playerTwo)
	}
}

//JogarMatch simula a execucao de um Match
func JogarMatch(playerOne *Player, playerTwo *Player, waitGroup *sync.WaitGroup, mutexObject *sync.Mutex){

	for{
		/*
				Verificando se algum dos player possui 
				os requisitos para vencer o jogo
			*/
		if playerOne.sets == SETSPARAVITORIA || playerTwo.sets == SETSPARAVITORIA{
			if playerOne.sets > playerTwo.sets{
				playerOne.venceuJogo = true
			} else{
				playerTwo.venceuJogo = true
			}
			break
		}

		/*resetando atributos dos players*/
		playerOne.games = 0
		playerOne.venceuSets = false
		playerTwo.games = 0
		playerTwo.venceuSets = false

		JogarSet(playerOne, playerTwo, waitGroup, mutexObject)

		printResultadoSet(playerOne, playerTwo)

		printPlacarSet(playerOne, playerTwo)
		
	}
}

func main(){
	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	/*
		inicializando player 1
	*/
	playerOne := Player{
		nome: "Player 1",
		pontos: 0,
		sets: 0,
		games: 0,
		estado: "rebatendo",
		venceuGame: false,
		venceuSets: false,
		venceuJogo: false,
	} 
	/*
		inicializando player 2
	*/
	playerTwo := Player{
		nome: "Player 2",
		pontos: 0,
		games: 0,
		sets: 0,
		estado: "esperando",
		venceuGame: false,
		venceuSets: false,
		venceuJogo: false,
	} 

	JogarMatch(&playerOne, &playerTwo, &waitGroup, &mutex)

	printResultadoMatch(&playerOne, &playerTwo)
}
//funcao para imprimir o resultado do Game
func printResultadoGame(playerOne *Player, playerTwo *Player){
	fmt.Println()
	fmt.Println("----------------------------------")
	fmt.Println("Pontos", playerOne.nome, "=>", playerOne.pontos)
	fmt.Println("Pontos", playerTwo.nome, "=>", playerTwo.pontos)
	if playerOne.venceuGame{
			fmt.Println(playerOne.nome, "VENCEU O GAME!! \nGanhou com ", playerOne.pontos, " pontos")
			playerOne.games++
	} else{
			fmt.Println(playerTwo.nome, "VENCEU O GAME!! \nGanhou com ", playerTwo.pontos, " pontos")
			playerTwo.games++
	}
	
	fmt.Println("O game teve ", RODADA, " Rodadas")
	fmt.Println("----------------------------------")
	fmt.Println()
}
//funcao para imprimir o resultado do Set
func printResultadoSet(playerOne *Player, playerTwo *Player){
	fmt.Println()
	fmt.Println("----------------------------------")
	fmt.Println("Games", playerOne.nome, "=>", playerOne.games)
	fmt.Println("Games", playerTwo.nome, "=>", playerTwo.games)
	if playerOne.venceuSets{
		fmt.Println(playerOne.nome, "VENCEU O SET!! \nGanhou ", playerOne.games, " games")
		playerOne.sets++
	} else{
		fmt.Println(playerTwo.nome, "VENCEU O SET!! \nGanhou ", playerTwo.games, " games")
		playerTwo.sets++
	}
	fmt.Println("----------------------------------")
}
//funcao para imprimir o resultado do Match
func printResultadoMatch(playerOne *Player, playerTwo *Player){
	fmt.Println()
	fmt.Println("----------------------------------")
	fmt.Println("Sets", playerOne.nome, "=>", playerOne.sets)
	fmt.Println("Sets", playerTwo.nome, "=>", playerTwo.sets)
	if playerOne.venceuSets{
		fmt.Println(playerOne.nome, "VENCEU O JOGO!! \nGanhou ", playerOne.sets, " sets")
		playerOne.sets++
	} else{
		fmt.Println(playerTwo.nome, "VENCEU O JOGO!! \nGanhou ", playerTwo.sets, " sets")
		playerTwo.sets++
	}
	fmt.Println("----------------------------------")
}
//funcao para imprimir o placar da Rodada
func printPlacarRodada(playerOne *Player, playerTwo *Player){
	fmt.Println("----------------------------------")
	fmt.Println(playerOne.nome, "pontos:", playerOne.pontos)
	fmt.Println(playerTwo.nome, "pontos:", playerTwo.pontos)
	fmt.Println("----------------------------------")
	RODADA++
	fmt.Println("Rodada", RODADA, "ponto do", playerTwo.nome)
	fmt.Println("==================================")
	fmt.Println()
}
//funcao para imprimir o resultado do Game
func printPlacarGame(playerOne *Player, playerTwo *Player){
	fmt.Println()
	fmt.Println("==================================")
	fmt.Println("Games", playerOne.nome, "=>", playerOne.games)
	fmt.Println("Games", playerTwo.nome, "=>", playerTwo.games)
	fmt.Println("==================================")
	fmt.Println()
}
//funcao para imprimir o resultado do Set
func printPlacarSet(playerOne *Player, playerTwo *Player){
	fmt.Println()
	fmt.Println("==================================")
	fmt.Println("Sets", playerOne.nome, "=>", playerOne.sets)
	fmt.Println("Sets", playerTwo.nome, "=>", playerTwo.sets)
	fmt.Println("==================================")
	fmt.Println()
}
