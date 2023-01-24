package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
)

// コストの取得
func GetCost(svc costexploreriface.CostExplorerAPI) (result *costexplorer.GetCostAndUsageOutput) {

	granularity := aws.String("DAILY")

	metric := "UnblendedCost"
	metrics := []*string{&metric}

	// TimePeriod1
	// 現在時刻の取得
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().UTC().In(jst)
	dayBefore := now.AddDate(0, 0, -1)

	nowDate := now.Format("2006-01-02")
	dateBefore := dayBefore.Format("2006-01-02")

	// 昨日から今日まで
	timePeriod := costexplorer.DateInterval{
		Start: aws.String(dateBefore),
		End:   aws.String(nowDate),
	}

	group := costexplorer.GroupDefinition{
		Key:  aws.String("SERVICE"),
		Type: aws.String("DIMENSION"),
	}
	groups := []*costexplorer.GroupDefinition{&group}

	input := costexplorer.GetCostAndUsageInput{}
	input.Granularity = granularity
	input.Metrics = metrics
	input.TimePeriod = &timePeriod
	input.GroupBy = groups

	// 処理実行
	result, err := svc.GetCostAndUsage(&input)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}

//処理実行
func run() error {
	svc := costexplorer.New(session.Must(session.NewSession()))
	cost := GetCost(svc)

	fmt.Println("success")

	file, err := os.Create("AWS.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//jsonエンコード
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cost); err != nil {
		log.Fatal(err)
	}
	return nil
}

func main() {
	run()
}
