package main

import (
	"bufio"
	"github.com/kataras/golog"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	//read , modify with spaces and separating file line by line
	readFile, err := os.Open("C:\\Users\\ASUS\\Desktop\\test.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	//temporary backup
	//tempLines := fileTextLines

	readFile.Close()

	// make and open output file for Basic2PL algorithm
	fo, err := os.Create("C:\\Users\\ASUS\\Desktop\\BasicTimeStamp.txt")
	if err != nil {
		panic(err)
	}


	// loop for iterating throw all schedules
	for i:=0 ; i<len(fileTextLines) ; i++ {

		//string variable for holding result
		result := ""

		transHolder := make([]string, 8)

		tsCounterArr := make([]int, 8)

		TSCounter := 1

		// a table which has r and w as columns and variables(v-z) as rows
		//writeTs and readTS in algorithm
		helpTable := make([][]string, 5)// making 5 rows of variables
		for i := range helpTable {
			helpTable[i] = make([]string, 2)//making 2 columns for each row : 0 for read and 1 for write lock
		}
		golog.Info("helpTable : ",helpTable)


		step:=0

		line := fileTextLines[i]
		// loop for iterating throw each schedule content
		for j:=0 ; j<len(line) ; j+=step {

			golog.Info("line : ",line)
			golog.Info("result : ",result)
			golog.Info("len(line) : ",len(line))
			if j+6 >= len(line) {
				break
			}

			golog.Info("line[j+2]-48 : ",line[j+2]-48)
			golog.Info("j : ",j)
			if tsCounterArr[line[j+2]-48] == 0 {
				tsCounterArr[line[j+2]-48] = TSCounter
				TSCounter++
			}
			golog.Info("tsCounterArr : ", tsCounterArr)
			golog.Info("TSCounter : ", TSCounter)

			flag1 := 0
			if line[j] == 'w' {

				step = 6

				var rTemp int = 0
				var wTemp int = 0

				golog.Info("line[j+4]-118 : ",line[j+4]-118)
				golog.Info("string(line[j+4]) : ",string(line[j+4]))
				golog.Info("jj : ",j)
				if helpTable[line[j+4]-118][0] != "" {
					rTemp, _ = strconv.Atoi(helpTable[line[j+4]-118][0])
				}

				if helpTable[line[j+4]-118][1] != "" {
					wTemp, _ = strconv.Atoi(helpTable[line[j+4]-118][1])
				}

				golog.Info("rTemp : ", rTemp)
				golog.Info("wTemp : ", wTemp)

				tsTemp := int(tsCounterArr[line[j+2]-48])
				golog.Info("tsTemp : ", tsTemp)

				golog.Info("rTemp <= tsTemp : ", rTemp <= tsTemp)
				golog.Info(" wTemp <= tsTemp : ", wTemp <= tsTemp)
				golog.Info("rTemp <= tsTemp  &&  wTemp <= tsTemp : ", rTemp <= tsTemp && wTemp <= tsTemp)

				if rTemp <= tsTemp && wTemp <= tsTemp {

					helpTable[line[j+4]-118][1] = strconv.Itoa(tsCounterArr[line[j+2]-48])
					result = result + line[j:j+6]

				} else {
					flag1 = 1 // cascade rollback should be handled
				}

			} else if line[j] == 'r' {

				step = 6

				var wTemp int = 0

				if helpTable[line[j+4]-118][1] != "" {
					wTemp, _ = strconv.Atoi(helpTable[line[j+4]-118][1])
				}

				golog.Info("wTemp : ", wTemp)

				tsTemp := int(tsCounterArr[line[j+2]-48])
				golog.Info("tsTemp : ", tsTemp)

				golog.Info(" wTemp <= tsTemp : ", wTemp <= tsTemp)

				if wTemp <= tsTemp {

					helpTable[line[j+4]-118][0] = strconv.Itoa(tsCounterArr[line[j+2]-48])
					result = result + line[j:j+6]

				} else {
					flag1 = 1 // cascade rollback should be handled
				}

			} else {

				step = 4
				result = result + line[j:j+4]

			}

			if flag1 == 1 { // handling cascade rollback

				temp := line[j+2]-48
				index := strings.Index(line, strconv.Itoa(int(temp))+"," )
				for index != -1 {
					golog.Info("1111")

					if index <= j+2 {
						golog.Info("****************************")
						step -= 6
					}
					transHolder[temp] += line[index : index+6]
					line = line[:index] + line[index+6:]
					index = strings.Index(line, strconv.Itoa(int(temp))+"," )

				}

				index = strings.Index(line, strconv.Itoa(int(temp)) )
				transHolder[temp] += line[index : index+4]
				line = line[:index] + line[index+4:]


				//modifying result string after cascading rollback
				resIndex := strings.Index(result, strconv.Itoa(int(temp))+"," )
				for resIndex != -1 {
					golog.Info("2222")
					result = result[:resIndex] + result[resIndex+6:]
					resIndex = strings.Index(result, strconv.Itoa(int(temp))+"," )
				}


			}

		}

		////////////////////

		if _, err := fo.Write([]byte(result+"\r\n")); err != nil {
			panic(err)
		}

	}



	//close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

}