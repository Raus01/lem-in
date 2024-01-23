package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Nodes = map[string][]string{}
var Paths [][]string
var count int
var UnrelatedEdges [][]int

type Ant struct {
	id           int
	path         []string
	step         int
	currentPlace int
}

// https://www.youtube.com/watch?v=bSZ57h7GN2w
// Graph sructure
// Graph represents an adjacency list graph
type Graph struct {
	vertices []*Colony //graph structure would need a list that holds vertices the	pointers of vertices in in a slice and for the vertex struct
}

// Colony structure
// vertex represents a graph vertex
type Colony struct {
	roomNames string
	Edges     []*Colony // vertices of the neighbours. that would also just look like a slice of pointers to vertices first
}

func main() {
	test := &Graph{}
	/*roomNames := getroomName()

	for _, i := range roomNames {
		test.AddVertex(i)
	}*/

	textlines := ReadFile()
	for i := range textlines {
		if match, _ := regexp.MatchString("^\\w+-\\w+$", textlines[i]); match {
			edge := textlines[i]
			Edges := strings.Split(edge, "-")
			u := Edges[0]
			v := Edges[1]
			test.AddEdge(u, v)

		}
	}

	//test.Print()
	//Edges := getEdge()

	Ant := getAnt()
	//Start := getStart()
	//end := getEnd()

	roomNames, err := getNodeNames(textlines)
	if !err {
		fmt.Println("bad format when get")
		return
	}

	nodeConnections, err := getNodeConnections(textlines, roomNames)
	if !err {
		fmt.Println("bad format when find connections")
		return
	}
	fillMap(roomNames, nodeConnections)
	path := []string{roomNames[0]}
	end := roomNames[len(roomNames)-1]
	findPath(path[0], end, path)

	//picture(nodeCoordinates, roomNames)
	unrelatedArrays := findArrayOfNotRealtedPaths()
	heights := findCountOfAuntMovesForEveryArray(unrelatedArrays, Ant)
	index := getMinIndex(heights)

	var array []int
	for i := range Paths {
		findUnrelated(array, i)
	}
	if index >= 0 {
		textlines := ReadFile()
		for _, line := range textlines {
			fmt.Println(line)
		}
		fmt.Println()
	}
	sendAnts(index, Ant, roomNames, unrelatedArrays)
	//            fmt.Println("unrelatedeedges3", UnrelatedEdges)
}

func getNodeNames(s []string) ([]string, bool) {
	var nodeNames []string
	var start, finish int
	for i, pp := range s {
		if pp == "##start" {
			start = i + 1
		}
		if pp == "##end" {
			finish = i + 1
		}

	}
	if start == 0 || finish == 0 {
		return nodeNames, false
	}

	for i := range s {
		pps := strings.Split(s[i], " ")
		if len(pps) == 3 && i != finish && i != start {
			nodeNames = append(nodeNames, pps[0])
		}
	}
	pps := strings.Split(s[finish], " ")
	if len(pps) == 3 {
		nodeNames = append(nodeNames, pps[0])
	} else {
		return nodeNames, false
	}

	copy := []string{}
	pps = strings.Split(s[start], " ")
	if len(pps) == 3 {
		copy = append(copy, pps[0])
	} else {
		return nodeNames, false
	}
	for _, v := range nodeNames {
		copy = append(copy, v)
	}
	nodeNames = copy

	return nodeNames, true
}

func getNodeConnections(s []string, roomNames []string) ([][]string, bool) {
	nodeConnections := make([][]string, len(roomNames))
	for i := range nodeConnections {
		nodeConnections[i] = make([]string, 0)
	}
	for i := range s {
		pps := strings.Split(s[i], "-")
		if len(pps) == 2 {
			//fmt.Println(pps)
			for j := 0; j < len(roomNames); j++ {
				if roomNames[j] == pps[0] {
					nodeConnections[j] = append(nodeConnections[j], pps[1])
				}
				if roomNames[j] == pps[1] {
					nodeConnections[j] = append(nodeConnections[j], pps[0])
				}
			}
		}

	}
	return nodeConnections, true
}

func fillMap(roomNames []string, nodeConnections [][]string) {
	for i := range roomNames {
		Nodes[roomNames[i]] = nodeConnections[i]
	}
}

func findPath(current string, end string, path []string) {
	if current == end {
		count++
		copy := make([]string, len(path)-2)
		for i := 0; i < len(path)-2; i++ {
			copy[i] = path[i+1]
		}
		Paths = append(Paths, copy)
		return
	}
	for _, v := range Nodes[current] {
		flag := true
		for _, f := range path {
			if v == f {
				flag = false
			}
		}
		if flag {
			copy := path
			copy = append(copy, v)
			findPath(v, end, copy)
		}

	}

}

func findArrayOfNotRealtedPaths() [][]int {
	arrayOfNotRelatedPaths := make([][][]string, len(Paths))
	for i := range arrayOfNotRelatedPaths {
		arrayOfNotRelatedPaths[i] = make([][]string, 0)
	}

	arrayOfNotRelatedPathsNames := make([][]int, len(Paths))
	for i := range arrayOfNotRelatedPaths {
		arrayOfNotRelatedPaths[i] = append(arrayOfNotRelatedPaths[i], Paths[i])
		arrayOfNotRelatedPathsNames[i] = append(arrayOfNotRelatedPathsNames[i], i)
		for j := i; j < len(Paths); j++ {
			flag := true
			v := Paths[j]
			for _, f := range arrayOfNotRelatedPaths[i] {
				if !isPathsUnique(v, f) {
					flag = false
				}
			}
			if flag && len(v) != 0 {
				arrayOfNotRelatedPaths[i] = append(arrayOfNotRelatedPaths[i], v)
				arrayOfNotRelatedPathsNames[i] = append(arrayOfNotRelatedPathsNames[i], j)
			}
		}
	}

	return arrayOfNotRelatedPathsNames
}

func findUnrelated(array []int, index int) {
	if len(array) == 0 {
		array = append(array, index)
	}
	flag := true
	for i := index + 1; i < len(Paths); i++ {
		flag = true
		for _, v := range array {
			if !isPathsUniqueByIndex(i, v) {
				flag = false
			}
		}
		if flag {
			copy := array
			copy = append(copy, i)
			// fmt.Println(copy)
			findUnrelated(copy, i)
			continue
		}
	}
	flag = false
	for _, v := range UnrelatedEdges {
		if isSubclasses(v, array) {
			flag = true
		}
	}
	if !flag || len(UnrelatedEdges) == 0 {
		UnrelatedEdges = append(UnrelatedEdges, array)
	}

}

func getMinIndex(s []int) int {
	if len(s) == 0 {
		return -1
	}
	min := s[0]
	mini := 0
	for i, v := range s {
		if v < min {
			min = v
			mini = i
		}
	}
	return mini
}

func findCountOfAuntMovesForEveryArray(UnrelatedEdges [][]int, antCount int) []int {
	var heights []int
	flagForZero := false
	for _, v := range UnrelatedEdges {

		var heightx, newHeight []int

		for _, f := range v {
			heightx = append(heightx, height(f))
			newHeight = append(newHeight, height(f))
			if height(f) == 0 {
				flagForZero = true
			}
		}
		if flagForZero {
			heights = append(heights, 1)
			continue
		}

		n := antCount
		for n > 0 {
			index := getMinIndex(newHeight)
			newHeight[index]++
			n--
		}
		max := 0
		for i := range heightx {
			if heightx[i] != newHeight[i] {
				if max < newHeight[i] {
					max = newHeight[i]
				}
			}
		}

		heights = append(heights, max)
	}

	return heights
}

func height(index int) int {
	return len(Paths[index])
}

func isSubclasses(s1 []int, s2 []int) bool {
	min := len(s1)
	if len(s2) < min {
		min = len(s2)
	}
	count := 0
	for i := range s1 {
		for j := range s2 {
			if s1[i] == s2[j] {
				count++
			}
		}
	}
	if count == min {
		return true
	}
	return false
}

func isPathsUnique(Path1 []string, Path2 []string) bool {
	for _, v := range Path1 {
		for _, f := range Path2 {
			if v == f {
				return false
			}
		}
	}
	return true
}

func isPathsUniqueByIndex(Index1 int, Index2 int) bool {
	for _, v := range Paths[Index1] {
		for _, f := range Paths[Index2] {
			if v == f {
				return false
			}
		}
	}
	return true
}

func sendAnts(routeIndex int, antCount int, roomNames []string, unrelatedArrays [][]int) {
	if routeIndex < 0 {
		fmt.Println("no path to end")
		return
	}

	end := roomNames[len(roomNames)-1]
	//fmt.Println("rouutindeks", routeIndex)
	route := unrelatedArrays[routeIndex]
	//fmt.Println("rouut", route, "rouutindeks", routeIndex)
	//tuleb sorteerida pikkuse järgi
	for i := 0; i < len(route)-1; i++ { // votsin siins i < len(route)-1
		for j := i + 1; j < len(route); j++ {

			if height(route[i]) > height(route[j]) {
				n := route[i]
				route[i] = route[j]
				route[j] = n
			}
		}
	}
	// sorteerimise lõpetanud

	var Routes [][]string
	for i := range route {
		var path []string
		path = append(Paths[route[i]], end)
		Routes = append(Routes, path)
	}
	// leidke rida selle kohta, kui palju sipelgaid peate igal rajal jooksma
	var heights []int
	var heightx, newHeight []int
	flagForZero := false
	for _, f := range route {
		// fmt.Println(f, height(f))
		heightx = append(heightx, height(f))
		newHeight = append(newHeight, height(f))
		if height(f) == 0 {
			flagForZero = true
		}
	}
	if flagForZero { // nulli läbimise korral

		ants := make([]Ant, antCount)
		var a []string
		a = append(a, end)
		for i := range ants {
			ants = append(ants, Ant{id: i, path: a, currentPlace: 0, step: 0})
		}
		sx := make([][]string, len(ants))

		for i := range ants {
			Ant1 := ants[i]
			for Ant1.currentPlace < len(Ant1.path) {
				s := "L" + strconv.Itoa(Ant1.id+1) + "-" + Ant1.path[Ant1.currentPlace]
				if len(s) != 0 {
					sx[i] = append(sx[i], s)
				}
				Ant1.currentPlace++
			}
		}

		sf := make([]string, 1)
		for i, v := range sx {
			for j := range v {
				sf[j+ants[i].step] += v[j] + " "
			}
		}
		for _, v := range sf {
			fmt.Println(v)
		}
		return
	}

	n := antCount
	for n > 0 {
		index := getMinIndex(newHeight)
		newHeight[index]++
		n--
	}
	max := newHeight[0]
	for i := range heightx {
		if newHeight[i] != heightx[i] && newHeight[i] > max {
			max = newHeight[i]
		}
		heights = append(heights, newHeight[i]-heightx[i])
	}

	var ants []Ant
	sum := 0
	for _, v := range heights {
		sum += v
	}
	count := 0
	steps := 0
	for sum > 0 {
		for i, v := range heights {
			if v != 0 {
				heights[i]--
				ants = append(ants, Ant{id: count, path: Routes[i], step: steps})
				count++
			}
		}
		steps++
		sum = 0
		for _, v := range heights {
			sum += v
		}
		//fmt.Println(heights, sum)
	}
	// for _, v := range ants {
	// 	fmt.Println(v)
	// }
	sx := make([][]string, len(ants))

	for i := range ants {
		Ant1 := ants[i]
		for Ant1.currentPlace < len(Ant1.path) {
			s := "L" + strconv.Itoa(Ant1.id+1) + "-" + Ant1.path[Ant1.currentPlace]
			sx[i] = append(sx[i], s)
			Ant1.currentPlace++
		}
	}

	sf := make([]string, max)
	for i, v := range sx {
		for j := range v {
			sf[j+ants[i].step] += v[j] + " "
		}
	}
	for _, v := range sf {
		fmt.Println(v)
	}

}

/*func getroomName() []string {
	var roomNames []string
	textlines := ReadFile()
	for i := range textlines {
		if match, _ := regexp.MatchString("^\\w+ \\w+ \\w+$", textlines[i]); match {
			rooms := textlines[i]
			roomLine := strings.Split(rooms, " ")
			roomName := roomLine[0]
			roomNames = append(roomNames, roomName)

		}
	}
	return roomNames
}*/

/*func getEdge() [][]string {
	var Edges [][]string
	textlines := ReadFile()
	for i := range textlines {
		if match, _ := regexp.MatchString("^\\w+-\\w+$", textlines[i]); match {

			edge := textlines[i]

			roomline := strings.Split(edge, "-")

			//Edge := roomline[0]

			Edges = append(Edges, roomline)
			//fmt.Println(Edges)

		}
	}
	return Edges

}
*/
// Add Vertex. adds a vertex to the graph
//we will write it as a method to the graph with a pointer receiver and it's
// going to take in an integer which is the key to add and inside the body you
//first create a vertex that has k as the key and then we're going
// to append this to the vertices list in the graph

// conteins. function to prevent duplicate key names
// we can also panic, but we do not program just stop here
func contains(s []*Colony, k string) bool {
	for _, v := range s {
		if k == v.roomNames {
			return true
		}
	}
	return false
}

// print test tegi sellise asja, trükkis "list of addresses"
// &{[0xc000022040 0xc000022060 0xc000022080 0xc0000220c0 0xc0000220e0]}
// kui tahame et trükitaks key, tuleb teha väike lisafunktsioon
// print will print the adjacent list of each vertex on the graph

// Reads file and returns []
func ReadFile() []string {

	readfile, _ := os.Open(os.Args[1])
	defer readfile.Close()

	var textlines []string

	scanner := bufio.NewScanner(readfile)
	for scanner.Scan() {
		textlines = append(textlines, scanner.Text())
	}

	return textlines
}
func getAnt() int {
	textlines := ReadFile()
	Ant, _ := strconv.Atoi(textlines[0])
	if Ant == 0 {
		fmt.Println("Please add ants")
	}
	return Ant
}

/*func (g *Graph) AddVertex(k string) {
	//see kontrollib, kas sellise nimega key on juba olemas
	if contains(g.vertices, k) {
		err := fmt.Errorf("Vertex %v not added because it is an existing key", k)
		fmt.Println(err.Error())

	} else {
		g.vertices = append(g.vertices, &Colony{roomNames: k})
	}
}*/

/*func (g *Graph) Print() {

	for _, v := range g.vertices {
		fmt.Printf("\nVertex %v : ", v.roomNames) // here we print out the key for each vertex
		for _, v := range v.Edges {
			fmt.Printf(" %v ", v.roomNames)
		}
	}
	fmt.Println()
}*/

// Add Edge (for directed graph)
func (g *Graph) AddEdge(from, to string) {
	// get vertex
	fromVertex := g.getVertex(from)
	toVertex := g.getVertex(to)
	//check errors
	if fromVertex == nil || toVertex == nil {
		//	err := fmt.Errorf("Invalid edge (%v-->%v)", from, to)
		//	fmt.Println(err.Error())
	} else if contains(fromVertex.Edges, to) {
		//err := fmt.Errorf("Existing edge (%v-->%v)", from, to)
		//fmt.Println(err.Error())
	} else {
		// add edge
		fromVertex.Edges = append(fromVertex.Edges, toVertex)
		toVertex.Edges = append(toVertex.Edges, fromVertex)
		//right now i just went to the adjacency list of the from vertex and added a
		// pointer of the two vertex so just with this we can start justadding edges to our graph
	}
}

// simple function that can get vertex
// getVertex returns a pointer to the Vertex with a key integer
func (g *Graph) getVertex(k string) *Colony {
	for i, v := range g.vertices {
		if v.roomNames == k {
			return g.vertices[i]
		}
	}
	return nil
}
