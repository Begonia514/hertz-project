package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

type Student struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Sex     string   `json:"sex"`
	College College  `json:"college"`
	Email   []string `json:"email"`
}

type College struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// 假设你的服务端地址和端口号为 "http://127.0.0.1:8888"
var serverURL = "http://127.0.0.1:8888"

func BenchmarkAddStudentInfo(b *testing.B) {
	// 准备测试数据
	studentData := Student{
		ID:   1,
		Name: "Emma",
		Sex:  "female",
		College: College{
			Name:    "software college",
			Address: "逸夫",
		},
		Email: []string{"emma@nju.com"},
	}

	// 将测试数据转换为JSON格式
	payload, err := json.Marshal(studentData)
	if err != nil {
		b.Fatalf("Error marshaling JSON: %v", err)
	}

	// 开始性能测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		// 发送 POST 请求
		req, err := http.NewRequest("POST", serverURL+"/add-student-info", bytes.NewBuffer(payload))
		if err != nil {
			b.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			b.Errorf("Handler returned wrong status code: got %v, want %v", resp.StatusCode, http.StatusOK)
		}
	}
}

func BenchmarkQueryStudentInfo(b *testing.B) {
	// 准备测试数据（ID为1，对应上面添加的数据）
	studentID := 1

	// 开始性能测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		// 发送 GET 请求
		req, err := http.NewRequest("GET", serverURL+"/query?id="+strconv.Itoa(studentID), nil)
		if err != nil {
			b.Fatalf("Error creating request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			b.Errorf("Handler returned wrong status code: got %v, want %v", resp.StatusCode, http.StatusOK)
		}
	}
}
