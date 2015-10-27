package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"bitbucket.org/tshannon/config"
	"mutantmonkey.in/code/golinx/progress"
)

type RespOkJSON struct {
	Filename   string
	Url        string
	Delete_Key string
	Expiry     string
	Size       string
}

type RespErrJSON struct {
	Error string
}

var Config struct {
	siteurl string
	logfile string
	apikey  string
}

var keys map[string]string

func main() {
	parseConfig()
	getKeys()

	var del bool
	var randomize bool
	var expiry int64
	var deleteKey string

	flag.BoolVar(&del, "d", false,
		"Delete file at url (ex: -d https://linx.example.com/myphoto.jpg")
	flag.BoolVar(&randomize, "r", false,
		"Randomize filename")
	flag.Int64Var(&expiry, "e", 0,
		"Time in seconds until file expires (ex: -e 600)")
	flag.StringVar(&deleteKey, "deletekey", "",
		"Specify your own delete key for the upload(s) (ex: -deletekey mysecret)")
	flag.Parse()

	if del {
		for _, url := range flag.Args() {
			deleteUrl(url)
		}
	} else {
		for _, fileName := range flag.Args() {
			upload(fileName, deleteKey, randomize, expiry)
		}
	}
}

func upload(filePath string, deleteKey string, randomize bool, expiry int64) {
	fileInfo, err := os.Stat(filePath)
	checkErr(err)
	file, err := os.Open(filePath)
	checkErr(err)
	fileName := path.Base(file.Name())

	reader := progress.NewProgressReader(fileName, bufio.NewReader(file), fileInfo.Size())

	escapedFileName := url.QueryEscape(fileName)

	req, err := http.NewRequest("PUT", Config.siteurl+"upload/"+escapedFileName, reader)
	checkErr(err)

	req.Header.Set("User-Agent", "linx-client")
	req.Header.Set("Accept", "application/json")

	if Config.apikey != "" {
		req.Header.Set("Linx-Api-Key", Config.apikey)
	}
	if deleteKey != "" {
		req.Header.Set("Linx-Delete-Key", deleteKey)
	}
	if randomize {
		req.Header.Set("Linx-Randomize", "yes")
	}
	if expiry != 0 {
		req.Header.Set("Linx-Expiry", strconv.FormatInt(expiry, 10))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	if resp.StatusCode == 200 {
		var myResp RespOkJSON

		err := json.Unmarshal(body, &myResp)
		checkErr(err)

		fmt.Println(myResp.Url)

		addKey(myResp.Url, myResp.Delete_Key)

	} else if resp.StatusCode == 401 {

		checkErr(errors.New("Incorrect API key"))

	} else {
		var myResp RespErrJSON

		err := json.Unmarshal(body, &myResp)
		checkErr(err)

		fmt.Printf("Could not upload %s: %s\n", fileName, myResp.Error)
	}
}

func deleteUrl(url string) {
	deleteKey, exists := keys[url]
	if !exists {
		checkErr(errors.New("No delete key for " + url))
	}

	req, err := http.NewRequest("DELETE", url, nil)
	checkErr(err)

	req.Header.Set("User-Agent", "linx-client")
	req.Header.Set("Linx-Delete-Key", deleteKey)

	if Config.apikey != "" {
		req.Header.Set("Linx-Api-Key", Config.apikey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)

	if resp.StatusCode == 200 {
		fmt.Println("Deleted " + url)
		delete(keys, url)
		writeKeys()
	} else {
		checkErr(errors.New("Could not delete " + url))
	}

}

func addKey(url string, deleteKey string) {
	keys[url] = deleteKey
	writeKeys()
}

func getKeys() {
	keyFile, err := ioutil.ReadFile(Config.logfile)
	if os.IsNotExist(err) {
		keys = make(map[string]string)
		writeKeys()
		keyFile, err = ioutil.ReadFile(Config.logfile)
		checkErr(err)
	} else {
		checkErr(err)
	}

	err = json.Unmarshal(keyFile, &keys)
	checkErr(err)
}

func writeKeys() {
	byt, err := json.Marshal(keys)
	checkErr(err)

	err = ioutil.WriteFile(Config.logfile, byt, 0600)
	checkErr(err)
}

func parseConfig() {
	cfgFilePath := filepath.Join(getConfigDir(), "linx-client.conf")
	cfg, err := config.LoadOrCreate(cfgFilePath)
	checkErr(err)

	Config.siteurl = cfg.String("siteurl", "")
	Config.logfile = cfg.String("logfile", "")
	Config.apikey = cfg.String("apikey", "")

	if Config.siteurl == "" || Config.logfile == "" {
		fmt.Println("Configuring linx-client")
		fmt.Println()
		for Config.siteurl == "" {
			Config.siteurl = getInput("Site url (ex: https://linx.example.com/)", false)

			if lastChar := Config.siteurl[len(Config.siteurl)-1:]; lastChar != "/" {
				Config.siteurl = Config.siteurl + "/"
			}
		}
		cfg.SetValue("siteurl", Config.siteurl)

		for Config.logfile == "" {
			Config.logfile = getInput("Logfile path (ex: ~/.linxlog)", false)

			homeDir := getHomeDir()
			if lastChar := homeDir[len(homeDir)-1:]; lastChar != "/" {
				homeDir = homeDir + "/"
			}

			Config.logfile = strings.Replace(Config.logfile, "~/", homeDir, 1)

		}
		cfg.SetValue("logfile", Config.logfile)

		if Config.apikey == "" {
			Config.apikey = getInput("API key (leave blank if instance is public)", true)
		}
		cfg.SetValue("apikey", Config.apikey)

		cfg.Write()

		fmt.Printf("Configuration written at %s\n", cfgFilePath)
	}
}
