package scheduling

import (
	"bufio"
	"bytes"
	"cfst/core"
	"encoding/csv"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func RunCftTask() {
	executor := fmt.Sprintf("%s%s", core.Config.Cft.Root, core.Config.Cft.Execute)
	cmd := exec.Command(executor)
	cmd.Dir = core.Config.Cft.Root
	// 获取脚本输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("获取脚本输出失败: %v", err)
		return
	}

	// 在 Goroutine 中执行脚本并读取输出
	go func() {
		err := cmd.Start()
		if err != nil {
			log.Printf("脚本启动失败: %v", err)
			return
		}

		// 读取并打印脚本输出
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Println(scanner.Text()) // 打印脚本输出
		}
		if err := scanner.Err(); err != nil {
			log.Printf("读取脚本输出失败: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("脚本执行失败: %v", err)
		} else {
			log.Println("脚本执行成功")
		}
	}()
	// 在 Goroutine 中执行脚本
	//go func() {
	//	err := cmd.Run()
	//	if err != nil {
	//		log.Printf("脚本执行失败: %v", err)
	//		// 在这里可以添加其他错误处理逻辑
	//	} else {
	//		log.Println("脚本执行成功")
	//	}
	//}()
}

type SpeedTestResult struct {
	IPAddress      string `json:"ip_address"`
	Sent           string `json:"sent"`
	Received       string `json:"received"`
	PacketLossRate string `json:"packet_loss_rate"`
	AverageLatency string `json:"average_latency"`
	DownloadSpeed  string `json:"download_speed"`
}

func ReadResult() string {
	resultCsv := fmt.Sprintf("%s%s", core.Config.Cft.Root, core.Config.Cft.Result)
	file, err := os.Open(resultCsv)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var results []SpeedTestResult
	// 遍历 CSV 记录，跳过第一行表头
	for _, record := range records[1:] {
		result := SpeedTestResult{
			IPAddress:      record[0],
			Sent:           record[1],
			Received:       record[2],
			PacketLossRate: record[3],
			AverageLatency: record[4],
			DownloadSpeed:  record[5],
		}
		results = append(results, result)
	}
	//将results[0]记录自动更新只cloudflare 域名dns记录
	updateCfDns(results[0])
	// 将结果转换为 JSON 格式
	jsonData, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

type RequestCloudflareBody struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func updateCfDns(result SpeedTestResult) {
	ip := result.IPAddress
	body := &RequestCloudflareBody{
		Name:    core.Config.Cft.AnalysisName,
		Type:    core.Config.Cft.AnalysisType,
		Content: ip,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling body:", err)
		return
	}
	req, err := http.NewRequest("PUT", core.Config.Cft.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("X-Auth-Email", core.Config.Cft.Email)
	req.Header.Set("X-Auth-Key", core.Config.Cft.Key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}
