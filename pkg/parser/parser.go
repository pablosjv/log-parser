package parser

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const fileWaitTime = 1

type app struct {
	process              func()
	initTime             time.Time
	endTime              time.Time
	reportPeriod         int64
	targetHost           string
	originHost           string
	fileName             string
	connectedToTarget    map[string]int
	connectionFromOrigin map[string]int
	outboundConnections  map[string]int
}

// NewApp creates new app object
func NewApp(initTime int64, endTime int64, targetHost string, originHost string, fileName string, tail bool, reportPeriod int64) *app {
	a := app{}
	if tail {
		a.process = a.tailReport
	} else {
		a.process = a.getConnectedToTarget
	}
	a.initTime = time.Unix(initTime, 0)
	a.endTime = time.Unix(endTime, 0)
	a.targetHost = targetHost
	a.originHost = originHost
	a.fileName = fileName
	a.reportPeriod = reportPeriod
	a.connectedToTarget = map[string]int{}
	a.connectionFromOrigin = map[string]int{}
	a.outboundConnections = map[string]int{}

	return &a
}

// Run the parser
func (a *app) Run() {
	log.Println("Processing connected hosts")
	start := time.Now()
	a.process()
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("Processing time:", elapsed)
}

func (a *app) getConnectedToTarget() {
	// Compensate for the five minute disorder
	compensatedEndTime := a.endTime.Add(time.Minute * time.Duration(5))
	file, err := os.Open(a.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logTime, originHost, targetHost := parseLine(line)

		if logTime.Before(a.endTime) &&
			logTime.After(a.initTime) &&
			targetHost == a.targetHost {

			a.connectedToTarget[originHost]++

		} else if logTime.After(compensatedEndTime) {
			log.Println("No more host connection in the time")
			break
		}
	}
	log.Println("End of file")
	log.Println("Connected hosts to", a.targetHost, ":", a.connectedToTarget)
}

func (a *app) tailReport() {
	file, err := os.Open(a.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Suppose we only want to tail the file, not process any content before
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}

	a.endTime = time.Now().Add(time.Second * time.Duration(a.reportPeriod))
	for {
		if time.Now().After(a.endTime) {
			log.Println("Connected hosts to", a.targetHost, ":", a.connectedToTarget)
			log.Println("Connections made by host", a.originHost, ":", a.connectionFromOrigin)
			max := getMaxInMap(a.outboundConnections)
			log.Println("Host with maximun outboud connections:", max, "=", a.outboundConnections[max])
			a.connectedToTarget = map[string]int{}
			a.connectionFromOrigin = map[string]int{}
			a.outboundConnections = map[string]int{}
			a.endTime = time.Now().Add(time.Second * time.Duration(a.reportPeriod))
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// log.Println("Waiting")
				time.Sleep(fileWaitTime * time.Second)
				continue
			} else {
				break
			}
		}
		// Skip if we encouter a line break
		if line == "\n" {
			continue
		}
		_, originHost, targetHost := parseLine(line)
		if targetHost == a.targetHost {
			a.connectedToTarget[originHost]++
		}
		if originHost == a.originHost {
			a.connectionFromOrigin[targetHost]++
		}
		a.outboundConnections[originHost]++
	}
}

func getMaxInMap(m map[string]int) string {
	max := 0
	var name string
	for key, value := range m {
		if value > max {
			name = key
			max = value
		}
	}
	return name
}

func parseLine(line string) (time.Time, string, string) {
	elements := strings.Split(line, " ")
	t, err := strconv.ParseInt(elements[0], 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(t, 0), elements[1], strings.TrimSpace(elements[2])
}
