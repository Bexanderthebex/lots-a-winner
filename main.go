package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// valid value range for lottery picks
var minimumValidPick uint = 1
var maximumValidPick uint = 90

var lottoPickLength = 5

var maximumBettors = 10_000_000

func main() {
	// call file
	fileToIndex := os.Args[1:2]

	file, err := os.Open(fileToIndex[0])
	if err != nil {
		log.Fatal(err)
	}

	bitMap := NewBitMap(minimumValidPick, maximumValidPick, maximumBettors)
	lotteryBetsVisitor := NewLotteryBetsVisitor(bitMap, " ")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lottoBet := scanner.Text()
		lotteryBetsVisitor.Visit(lottoBet)
	}

	if fileCloseError := file.Close(); fileCloseError != nil {
		panic(fileCloseError)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("READY")

	// query engine
	for {
		//accumulator := make([]uint, bitMap.GetTotalRecords())
		text, readStringError := reader.ReadString('\n')
		if readStringError != nil {
			fmt.Println(readStringError)
			fmt.Println("Please enter the correct format")
		}

		text = strings.Replace(text, "\n", "", -1)
		winningPicks := make(map[uint]uint, 0)
		for _, bet := range strings.Split(text, " ") {
			if bet != "" {
				formattedBetString, _ := strconv.Atoi(bet)
				winningPicks[uint(formattedBetString)] += 1
			}
		}

		if len(winningPicks) < 5 {
			fmt.Printf("Please enter 5 winning picks, only %d were entered\n", len(winningPicks))
			continue
		}

		isPicksValid := picksValid(winningPicks)
		if !isPicksValid {
			fmt.Println("Please enter values that are only between 1 and 90")
			continue
		}

		fmt.Println("Winning picks parsed:")
		fmt.Println(winningPicks)

		queryPlan := NewSelectQueryPlan(bitMap)
		queryPlan.SetMinValue(2)
		queryPlan.SetMaxValue(5)
		queryPlan.SetColumnsToSelect(winningPicks)
		//queryPlan.AddAggregationStrategy(NewCountAggregation(&winningPicks, &accumulator, bitMap))
		//queryPlan.AddAggregationStrategy(NewGroupAggregation(&accumulator, []uint{5, 4, 3, 2}))
		queryPlan.AddAggregationStrategy(NewOptimizedAggregation(bitMap, winningPicks, 2, 5))
		queryEngine := LotteryBetsQueryEngine{
			bitmapIndex: bitMap,
		}

		answersMap := queryEngine.ExecuteQuery(queryPlan)

		for i := uint(5); i >= 2; i-- {
			fmt.Printf("%d: %d\n", i, answersMap[i])
		}
	}
}
