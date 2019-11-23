package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var addrs []string

var sendvar = 1
var cantLines = 748
var myip = ""

type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

var Blockchain []Block
var blocktemporal Block

var team0 = 0
var team1 = 0

var resultado string

var c = make(chan bool)

//Takes block data and create a SHA256 HASH of it

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//Generate block with all the elements needed

func generateBlock(oldBlock Block, Data string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

//Make sure the blocks havent been tampered with

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println("index")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println(oldBlock.Hash)
		fmt.Println(newBlock.PrevHash)
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		fmt.Println("calculatehash")
		return false
	}
	fmt.Println("true")
	return true
}

//Compare the lenght of the slices of blocks, the longest one is
//the most up to date, if we have the shortest, we replace it

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func handleNewBlock(data string) {
	newBlock, _ := generateBlock(Blockchain[len(Blockchain)-1], data)
	blocktemporal = newBlock
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}
}

func printBlock() {
	for i := range Blockchain {
		fmt.Println("Hash :")
		fmt.Println(Blockchain[i].Index)
		fmt.Println(Blockchain[i].Data)
		fmt.Println(Blockchain[i].Hash)
		fmt.Println(Blockchain[i].PrevHash)
	}
}

func blockServer(hostAddr string) {
	host := fmt.Sprintf("%s:8002", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleBlockchain(conn)
	}
}

func handleReceived(cad string) (num string, ip string) {
	res1 := strings.Split(cad, "*")
	index, _ := strconv.Atoi(res1[0])
	timestamp := res1[1]
	data := res1[2]
	hash := res1[3]
	prevhash := res1[4]
	ipFinal := res1[5]
	var newBlock Block
	newBlock.Index = index
	newBlock.Timestamp = timestamp
	newBlock.Data = data
	newBlock.Hash = hash
	newBlock.PrevHash = prevhash
	fmt.Println(prevhash)
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}

	return data, ipFinal
}
func handleBlockchain(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	strNum, _ := r.ReadString('\n')

	num := strings.TrimSpace(strNum)
	fmt.Println(num)
	tempdata, ipFinal := handleReceived(num)

	printBlock()
	if myip != ipFinal {
		blockSend(tempdata, ipFinal)
	}
	c <- true
}
func blockSend(num string, ip string) {
	idx := rand.Intn(len(addrs))
	fmt.Println(idx)
	blockTempString := strconv.Itoa(blocktemporal.Index) + "*" + blocktemporal.Timestamp + "*" + num + "*" + blocktemporal.Hash + "*" + blocktemporal.PrevHash + "*" + ip
	fmt.Println(blockTempString)
	remote := fmt.Sprintf("%s:8002", addrs[idx])
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprintln(conn, blockTempString)
	sendvar = -1
}
func registerSend(remoteAddr, hostAddr string) {
	remote := fmt.Sprintf("%s:8000", remoteAddr)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()

	// Enviar direccion
	fmt.Fprintln(conn, hostAddr)

	// Recibir lista de direcciones
	r := bufio.NewReader(conn)
	strAddrs, _ := r.ReadString('\n')
	var respAddrs []string
	json.Unmarshal([]byte(strAddrs), &respAddrs)

	// agregamos direcciones de nodos a propia libreta
	for _, addr := range respAddrs {
		if addr == remoteAddr {
			return
		}
	}
	addrs = append(respAddrs, remoteAddr)
	fmt.Println(addrs)
}
func registerServer(hostAddr string) {
	host := fmt.Sprintf("%s:8000", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleRegister(conn)
	}
}
func handleRegister(conn net.Conn) {
	defer conn.Close()

	// Recibimos addr del nuevo nodo
	r := bufio.NewReader(conn)
	remoteIp, _ := r.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	// respondemos enviando lista de direcciones de nodos actuales
	byteAddrs, _ := json.Marshal(addrs)
	fmt.Fprintf(conn, "%s\n", string(byteAddrs))

	// notificar a nodos actuales de llegada de nuevo nodo
	for _, addr := range addrs {
		notifySend(addr, remoteIp)
	}

	// Agregamos nuevo nodo a la lista de direcciones
	for _, addr := range addrs {
		if addr == remoteIp {
			return
		}
	}
	addrs = append(addrs, remoteIp)
	fmt.Println(addrs)
}
func notifySend(addr, remoteIp string) {
	remote := fmt.Sprintf("%s:8001", addr)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprintln(conn, remoteIp)
}
func notifyServer(hostAddr string) {
	host := fmt.Sprintf("%s:8001", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleNotify(conn)
	}
}
func handleNotify(conn net.Conn) {
	defer conn.Close()

	// Recibimos addr del nuevo nodo
	r := bufio.NewReader(conn)
	remoteIp, _ := r.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	// Agregamos nuevo nodo a la lista de direcciones
	for _, addr := range addrs {
		if addr == remoteIp {
			return
		}
	}
	addrs = append(addrs, remoteIp)
	fmt.Println(addrs)
}

func SendFuncHandle(gin *bufio.Reader) {
	if <-c == true {
		fmt.Print("Ingrese data : ")
		strNum, _ := gin.ReadString('\n')
		if strNum != "" {
			num := strings.TrimSpace(strNum)
			handleNewBlock(num)
			blockSend(num, myip)
			file, _ := os.Open("transfusion.csv")
			res1 := strings.Split(num, ",")
			x1, _ := strconv.Atoi(strings.Replace(res1[0], " ", "", -1))
			x2, _ := strconv.Atoi(strings.Replace(res1[1], " ", "", -1))
			x3, _ := strconv.Atoi(strings.Replace(res1[2], " ", "", -1))
			x4, _ := strconv.Atoi(strings.Replace(res1[3], " ", "", -1))
			distance, team := ReadFile2("transfusion.csv", 0, cantLines, x1, x2, x3, x4)
			quicksort(distance, team)
			distance = distance[:5]
			team = team[:5]
			for i := 0; i < 5; i++ {
				if team[i] == 0 {
					team0++
				} else {
					team1++
				}
			}
			if team0 > team1 {
				resultado = "0"
			} else {
				resultado = "1"
			}
			writer := csv.NewWriter(file)
			defer writer.Flush()
			num = num + "," + resultado + " "
			aux := strings.Split(num, " ")
			writer.Write(aux)
			cantLines++
			team0 = 0
			team1 = 0
			c <- false
		}
	}
}

func ReadFile2(file string, begin int, end int, y1 int, y2 int, y3 int, y4 int) ([]float64, []int) {
	distance := make([]float64, 0)
	team := make([]int, 0)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", file, err.Error())
	}
	defer f.Close()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
		if n < begin {
			continue
		}
		if n > end {
			break
		}
		res1 := strings.Split(scanner.Text(), ",")
		x1, _ := strconv.Atoi(strings.Replace(res1[0], " ", "", -1))
		x2, _ := strconv.Atoi(strings.Replace(res1[1], " ", "", -1))
		x3, _ := strconv.Atoi(strings.Replace(res1[2], " ", "", -1))
		x4, _ := strconv.Atoi(strings.Replace(res1[3], " ", "", -1))
		x5, _ := strconv.Atoi(strings.Replace(res1[4], " ", "", -1))
		distanciaP := distanciaEuclidiana(x1, x2, x3, x4, x5, y1, y2, y3, y4)
		distance = append(distance, distanciaP)
		team = append(team, x5)
	}
	fmt.Println(distance)
	fmt.Println(team)
	return (distance), team
}

func quicksort(a []float64, team []int) ([]float64, []int) {
	if len(a) < 2 {
		return a, team
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]
	team[pivot], team[right] = team[right], team[pivot]

	for i, _ := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			team[left], team[i] = team[i], team[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]
	team[left], team[right] = team[right], team[left]

	quicksort(a[:left], team[:left])
	quicksort(a[left+1:], team[left+1:])

	return a, team
}

func distanciaEuclidiana(x1 int, x2 int, x3 int, x4 int, x5 int, y1 int, y2 int, y3 int, y4 int) float64 {
	sum_cuadrados := math.Pow(float64(y1)-float64(x1), 2) + math.Pow(float64(y2)-float64(x2), 2) +
		math.Pow(float64(y3)-float64(x3), 2) + math.Pow(float64(y4)-float64(x4), 2)
	resultado := math.Sqrt(sum_cuadrados)
	return resultado
}

func main() {
	c <- true
	myip := "192.168.0.141"

	t := time.Now()
	genesisBlock := Block{0, t.String(), "", "", ""}
	Blockchain = append(Blockchain, genesisBlock)

	fmt.Printf("Soy %s\n", myip)
	go registerServer(myip)
	go blockServer(myip)

	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese direccion remota: ")
	remoteIp, _ := gin.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	if remoteIp != "" {
		registerSend(remoteIp, myip)
	}

	notifyServer(myip)

	for {
		go SendFuncHandle(gin)
	}

}
