package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const NAME = "Vitor"
const APPLICATION = "Monitor de sites."
const VERSION = 1.1

const MONITORING = 2
const DELAY = 5

type Options int

const (
	Finish Options = iota
	MonitorSites
	AddSites
	ShowSites
	ShowLogs
)

type HttpStatus int

const (
	Ok          HttpStatus = 200
	NotFound               = 404
	ServerError            = 500
)

func intro() {
	fmt.Println("Nome:", NAME, "\nAplicação:", APPLICATION, "\nVersão:", VERSION)
}

func getCommand() int {
	var command int

	fmt.Println("\n" + strconv.Itoa(int(MonitorSites)) + "- Iniciar Monitoramento")
	fmt.Println(strconv.Itoa(int(AddSites)) + "- Adicionar Sites")
	fmt.Println(strconv.Itoa(int(ShowSites)) + "- Lista de sites")
	fmt.Println(strconv.Itoa(int(ShowLogs)) + "- Logs")
	fmt.Println(strconv.Itoa(int(Finish)) + "- Sair")

	fmt.Scan(&command)

	return command
}

func testSite(site string) {
	response, err := http.Get(site)

	if err != nil {
		fmt.Println("Erro ao tentar ler o arquivo: ", err)
		return
	}

	switch response.StatusCode {
	case int(Ok):
		saveLogs(site, response.StatusCode)
		fmt.Println("Site funcionando: ", Ok)
	default:
		saveLogs(site, response.StatusCode)
		fmt.Println("Erro, status do site:", response.StatusCode)
	}
}

func getSitesFromFile() []string {
	var sites []string

	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Erro ao tentar ler o arquivo: ", err)
		file.Close()
		return nil
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}

	}

	file.Close()
	return sites
}

func startMonitor() {
	fmt.Println("\nIniciando monitoramento... ")
	sites := getSitesFromFile()

	if sites == nil {
		fmt.Println("\nNenhum site encontrado.")
		return
	}

	for i := 0; i < MONITORING; i++ {
		for _, site := range sites {
			fmt.Println("\nMonitorando site: ", site)
			testSite(site)
		}

		fmt.Println("\nAguardando ", DELAY, " segundos para o novo monitoramento...")
		time.Sleep(DELAY * time.Second)
	}

}

func addSites() {
	var site string
	fmt.Println("Digite a url completa do site: ")
	fmt.Scan(&site)

	file, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro ao abrir arquivo de sites: ", err)
		file.Close()
		return
	}

	file.WriteString("\n" + site)
	file.Close()
}

func saveLogs(site string, status int) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro ao criar/abrir arquivo de log: ", err)
		file.Close()
		return
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - Site: " + site + " - Status: " + strconv.Itoa(status) + "\n")
	file.Close()
}

func printSites() {
	file, err := ioutil.ReadFile("sites.txt")

	if err != nil {
		fmt.Println("Erro ao abrir aquivo de sites: ", err)
		return
	}

	fmt.Println(string(file))
}

func printLogs() {
	fmt.Println("Exibindo logs...")

	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Erro ao abrir aquivo de log: ", err)
		return
	}

	fmt.Println(string(file))
}

func main() {

	intro()

	for {
		command := getCommand()

		switch command {

		case int(MonitorSites):
			startMonitor()

		case int(AddSites):
			addSites()

		case int(ShowSites):
			printSites()

		case int(ShowLogs):
			printLogs()

		case int(Finish):
			fmt.Println("Finalizando aplicação")
			os.Exit(0)

		default:
			fmt.Println("Opção não encontrada, tente novamente.")
		}
	}
}
