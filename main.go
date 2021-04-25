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