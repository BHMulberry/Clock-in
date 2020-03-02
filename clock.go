package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
)

// LoginForm ...
type LoginForm struct {
	Code     string
	Password string
}

// DailyInfo ...
type DailyInfo struct {
	IsVia    string
	DateTrip string
}

// TripinfoList ...
type TripinfoList struct {
	ATripDate          string `json:"aTripDate"`
	FromAdr            string
	ToAdr              string
	Number             string
	Trippersoninfolist []string `json:"trippersoninfolist"`
}

// PostForm ...
type PostForm struct {
	Temperature         string
	RealProvince        string
	RealCity            string
	RealCounty          string
	RealAddress         string
	UnusualInfo         string
	IsUnusual           string
	IsTouch             string
	IsInsulated         string
	IsSuspected         string
	IsDiagnosis         string
	Dailyinfo           DailyInfo      `json:"dailyinfo"`
	Toucherinfolist     []string       `json:"toucherinfolist"`
	Tripinfolist        []TripinfoList `json:"tripinfolist"`
	IsInCampus          string
	IsViaHuBei          string
	IsViaWuHan          string
	InsulatedAddress    string
	TouchInfo           string
	IsNormalTemperature string
	Longitude           *string
	Latitude            *string
}

type data struct {
	Time         			string
	Code         			string
	Password     			string
	RealProvince 			string
	RealCity     			string
	RealCounty   			string
	RealAddress  			string
	RandomTimeFluctuation 	string
}

const loginURL = "http://fangkong.hnu.edu.cn/api/v1/account/login"
const addURL = "http://fangkong.hnu.edu.cn/api/v1/clockinlog/add"

const timeLayout = "2006-01-02 15:04:05"

const timeFluctuationRange = 1800

func readJSON(filePath string) (result string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	buf := bufio.NewReader(file)
	for {
		s, err := buf.ReadString('\n')
		result += s
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return
			}
		}
	}
	return result
}

func main() {
	var config data
	fmt.Println("读取配置...")
	result := readJSON("./config.json")
	err := json.Unmarshal([]byte(result), &config)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	loginForm := LoginForm{
		Code:     config.Code,
		Password: config.Password,
	}

	postForm := PostForm{
		Temperature:         "",
		RealProvince:        config.RealProvince,
		RealCity:            config.RealCity,
		RealCounty:          config.RealCounty,
		RealAddress:         config.RealAddress,
		IsUnusual:           "0",
		UnusualInfo:         "",
		IsTouch:             "0",
		IsInsulated:         "0",
		IsSuspected:         "0",
		IsDiagnosis:         "0",
		IsInCampus:          "0",
		IsViaHuBei:          "0",
		IsViaWuHan:          "0",
		InsulatedAddress:    "",
		TouchInfo:           "",
		IsNormalTemperature: "1",
		Dailyinfo: DailyInfo{
			IsVia:    "0",
			DateTrip: "",
		},
		Tripinfolist: []TripinfoList{
			TripinfoList{
				ATripDate:          "",
				FromAdr:            "",
				ToAdr:              "",
				Number:             "",
				Trippersoninfolist: []string{},
			},
		},
		Longitude:       nil,
		Latitude:        nil,
		Toucherinfolist: []string{},
	}

	h, _ := strconv.Atoi(config.Time[0:2])
	m, _ := strconv.Atoi(config.Time[3:5])
	s, _ := strconv.Atoi(config.Time[6:8])

	var testFlag bool

	if len(os.Args) > 1 && os.Args[1] == "test" {
		testFlag = true
	} else {
		testFlag = false
	}

	for {
		if testFlag == false {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), h, m, s, 0, next.Location())
			if config.RandomTimeFluctuation == "true"{
				rand.Seed(time.Now().UnixNano())
				next = next.Add(time.Second * time.Duration(rand.Intn(2*timeFluctuationRange)-timeFluctuationRange))
			}
			fmt.Println("下一次打卡时间为: ", next.Format(timeLayout))
			nextTimer := time.NewTimer(next.Sub(now))
			<-nextTimer.C
		}

		bytes, _ := json.Marshal(loginForm)
		loginInfo := strings.NewReader(string(bytes))

		client := &http.Client{}
		loginReq, _ := http.NewRequest("POST", loginURL, loginInfo)
		loginReq.Header.Add("Content-Type", "application/json")
		loginResq, _ := client.Do(loginReq)
		header := loginResq.Header
		cookies := header["Set-Cookie"]

		if len(cookies) == 0 {
			fmt.Println("登陆失败...请检查门户账号和密码")
			return
		}

		index1 := strings.Index(cookies[0], "=")
		index2 := strings.Index(cookies[0], ";")
		index3 := strings.Index(cookies[3], "=")
		index4 := strings.Index(cookies[3], ";")

		cookie1 := &http.Cookie{Name: "TOKEN", Value: cookies[0][index1+1 : index2]}
		cookie2 := &http.Cookie{Name: ".ASPXAUTH", Value: cookies[3][index3+1 : index4]}

		bytes, _ = json.Marshal(postForm)
		addInfo := strings.NewReader(string(bytes))

		addReq, _ := http.NewRequest("POST", addURL, addInfo)
		addReq.AddCookie(cookie1)
		addReq.AddCookie(cookie2)

		addReq.Header.Add("Content-Type", "application/json")

		addResq, err := client.Do(addReq)

		if err != nil {
			panic(err.Error())
		}
		fmt.Println(addResq.Status)
		bytes, _ = ioutil.ReadAll(addResq.Body)
		fmt.Println(string(bytes))
		if testFlag == true {
			return
		}
	}
}
