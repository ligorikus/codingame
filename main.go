package main

import (
	"fmt"
	"os"
	"sort"
)

const MaxDeep = 21

type Owner int64

type GameInput struct {
	owner    Owner
	entityId int
	arg2     int
	arg3     int
	arg4     int
	arg5     int
}

const (
	Player  Owner = 1
	Enemy   Owner = -1
	Neutral Owner = 0
)

type EntityType string

type Matrix [][]int

const (
	FactoryEntity EntityType = "FACTORY"
	TroopEntity   EntityType = "TROOP"
	BombEntity    EntityType = "BOMB"
)

type Factory struct {
	id                    int
	owner                 Owner
	garrison              int
	production            int
	turnsToStartProducing int
}

type Troop struct {
	id          int
	owner       Owner
	source      int
	destination int
	size        int
	arrival     int
}

type Bomb struct {
	id          int
	owner       Owner
	source      int
	destination int
	arrival     int
}

type FactoryOwner struct {
	player  []int
	neutral []int
	enemy   []int
}

type ScenarioType string

const (
	MoveScenario ScenarioType = "MOVE"
	BombScenario ScenarioType = "BOMB"
	IncScenario  ScenarioType = "INC"
)

type Scenario struct {
	scenarioType ScenarioType
	source       int
	destination  int
	size         int
}

type PlayerState struct {
	scriptList      []Scenario
	availableTroops AvailableTroops
	enemyTroops     AvailableTroops
}

type GameState struct {
	graph      Matrix
	costMatrix Matrix
	moveMatrix Matrix
	factories  map[int]Factory
	troops     map[int]Troop
	bombs      map[int]Bomb

	factoryOwner FactoryOwner
	playerState  PlayerState
}

type FactoryAttackCost struct {
	source      int
	destination int
	cost        int
}

type FactoryAttackCostList []FactoryAttackCost

type AvailableTroops map[int]int

func (facl FactoryAttackCostList) Len() int {
	return len(facl)
}

func (facl FactoryAttackCostList) Less(i, j int) bool {
	return facl[i].cost < facl[j].cost
}

func (facl FactoryAttackCostList) Swap(i, j int) {
	facl[i], facl[j] = facl[j], facl[i]
}

func (state *GameState) InitializeGraph() {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	state.graph = make(Matrix, factoryCount)
	for i := range state.graph {
		state.graph[i] = make([]int, factoryCount)
	}

	for i := 0; i < linkCount; i++ {
		var factory1, factory2, distance int
		fmt.Scan(&factory1, &factory2, &distance)
		state.graph[factory1][factory2] = distance
		state.graph[factory2][factory1] = distance
	}

	state.FloydWarshallWithPathReconstruction()
}

func (state *GameState) ClearState() {
	state.factories = map[int]Factory{}
	state.troops = map[int]Troop{}
	state.bombs = map[int]Bomb{}

	state.playerState.scriptList = make([]Scenario, 0)
	state.playerState.availableTroops = map[int]int{}
	state.playerState.enemyTroops = map[int]int{}

	state.ClearOwner()
}
func (state *GameState) ClearOwner() {
	state.factoryOwner.player = make([]int, 0)
	state.factoryOwner.neutral = make([]int, 0)
	state.factoryOwner.enemy = make([]int, 0)
}

func (state *GameState) UpdateState() {
	state.ClearState()

	var entityCount int
	fmt.Scan(&entityCount)

	for i := 0; i < entityCount; i++ {
		var entityId int
		var entityType EntityType
		var owner Owner
		var arg2, arg3, arg4, arg5 int
		fmt.Scan(&entityId, &entityType, &owner, &arg2, &arg3, &arg4, &arg5)
		input := GameInput{entityId: entityId, owner: owner, arg2: arg2, arg3: arg3, arg4: arg4, arg5: arg5}

		switch entityType {
		case FactoryEntity:
			state.InputFactory(input)
			break
		case TroopEntity:
			state.InputTroop(input)
			break
		case BombEntity:
			state.InputBomb(input)
			break
		}
	}
}

func (state *GameState) InputFactory(input GameInput) {
	state.factories[input.entityId] = Factory{
		id:                    input.entityId,
		owner:                 input.owner,
		garrison:              input.arg2,
		production:            input.arg3,
		turnsToStartProducing: input.arg4,
	}

	if input.owner == Player {
		state.factoryOwner.player = append(state.factoryOwner.player, input.entityId)
		state.playerState.availableTroops[input.entityId] = input.arg2
	}

	if input.owner == Neutral {
		state.factoryOwner.neutral = append(state.factoryOwner.neutral, input.entityId)
		state.playerState.enemyTroops[input.entityId] = input.arg2
	}

	if input.owner == Enemy {
		state.factoryOwner.enemy = append(state.factoryOwner.enemy, input.entityId)
		state.playerState.enemyTroops[input.entityId] = input.arg2
	}
}

func (state *GameState) InputTroop(input GameInput) {
	state.troops[input.entityId] = Troop{
		id:          input.entityId,
		owner:       input.owner,
		source:      input.arg2,
		destination: input.arg3,
		size:        input.arg4,
		arrival:     input.arg5,
	}
}

func (state *GameState) InputBomb(input GameInput) {
	state.bombs[input.entityId] = Bomb{
		id:          input.entityId,
		owner:       input.owner,
		source:      input.arg2,
		destination: input.arg3,
		arrival:     input.arg4,
	}
}

func (state *GameState) FloydWarshallWithPathReconstruction() {
	graph := state.graph
	dist := make([][]int, len(graph))
	next := make([][]int, len(graph))
	for i := 0; i < len(graph); i++ {
		dist[i] = make([]int, len(graph))
		next[i] = make([]int, len(graph))
		for j := 0; j < len(graph); j++ {
			if i == j {
				dist[i][j] = 0
				next[i][j] = -1
			} else {
				dist[i][j] = graph[i][j]
				next[i][j] = j
			}

		}
	}

	for k := 0; k < len(graph); k++ {
		for i := 0; i < len(graph); i++ {
			for j := 0; j < len(graph); j++ {
				if dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
	}

	state.costMatrix = dist
	state.moveMatrix = next
}

func (state *GameState) CalcMinDistanceByPlayer(playerFactories []int) []int {
	playerMinDistance := make([]int, len(state.factories))
	for i := range state.factories {
		playerMinDistance[i] = 21
	}

	for _, id := range playerFactories {
		for i, cost := range state.costMatrix[id] {
			if cost < playerMinDistance[i] {
				playerMinDistance[i] = cost
			}
		}
	}
	return playerMinDistance
}

func (state *GameState) GetInfluenceZone() []Owner {
	owner := make([]Owner, len(state.factories))
	playerMinDistance := state.CalcMinDistanceByPlayer(state.factoryOwner.player)
	enemyMinDistance := state.CalcMinDistanceByPlayer(state.factoryOwner.enemy)
	for i, cost := range playerMinDistance {
		if enemyMinDistance[i] > cost {
			owner[i] = Player
		} else if enemyMinDistance[i] < cost {
			owner[i] = Enemy
		} else {
			owner[i] = Neutral
		}
	}
	return owner
}

func GetInfluenceZoneByPlayer(influenceZone []Owner, player Owner) []int {
	result := make([]int, 0)
	for i, owner := range influenceZone {
		if owner == player || owner == Neutral {
			result = append(result, i)
		}
	}
	return result
}

func (state *GameState) CostOfAttack(source int, destination int) FactoryAttackCost {
	if state.factories[destination].production == 0 {
		return FactoryAttackCost{
			source:      source,
			destination: destination,
			cost:        255,
		}
	}
	cost := 0
	cost += state.costMatrix[source][destination]
	cost += (state.factories[destination].garrison / state.factories[destination].production) + 1

	return FactoryAttackCost{
		source:      source,
		destination: destination,
		cost:        cost,
	}
}

func (state *GameState) ProcessAttack() {
	influenceZone := state.GetInfluenceZone()
	playerInfluenceZone := GetInfluenceZoneByPlayer(influenceZone, Player)

	scenarioList := make(FactoryAttackCostList, 0)
	for _, playerFactory := range state.factoryOwner.player {
		for _, influenceFactory := range playerInfluenceZone {
			if state.factories[influenceFactory].owner == Player {
				continue
			}
			scenarioList = append(scenarioList, state.CostOfAttack(playerFactory, influenceFactory))
		}
	}
	sort.Sort(scenarioList)
	for _, scenario := range scenarioList {

		availableTroops := state.playerState.availableTroops[scenario.source]
		if availableTroops == 0 {
			continue
		}
		enemyTroops := state.playerState.enemyTroops[scenario.destination]
		delta := availableTroops - (enemyTroops + 1)

		var scenarioSize int
		if delta >= 0 {
			scenarioSize = (enemyTroops + 1)
			availableTroops -= (enemyTroops + 1)
			enemyTroops = 0
		} else {
			scenarioSize = availableTroops
			enemyTroops -= availableTroops
			availableTroops = 0
		}

		state.playerState.scriptList = append(state.playerState.scriptList, Scenario{
			scenarioType: MoveScenario,
			source:       scenario.source,
			destination:  scenario.destination,
			size:         scenarioSize,
		})
	}
}

func (state *GameState) ProcessScenario() {
	for _, scenario := range state.playerState.scriptList {
		switch scenario.scenarioType {
		case MoveScenario:
			fmt.Print("MOVE ", scenario.source, " ", scenario.destination, " ", scenario.size, ";")
			break
		}
	}
}

func (state *GameState) Process() {
	state.ProcessAttack()
	state.ProcessScenario()
}

func main() {
	state := GameState{}
	state.InitializeGraph()

	for _, item := range state.moveMatrix {
		fmt.Fprintln(os.Stderr, item)
	}
	fmt.Fprintln(os.Stderr)
	for _, item := range state.costMatrix {
		fmt.Fprintln(os.Stderr, item)
	}
	for {
		state.UpdateState()
		state.Process()

		fmt.Println("WAIT")
	}
}
