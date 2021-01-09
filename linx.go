package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/mutantmonkey/golinx/progress"
	"github.com/timshannon/config"
)

type RespOkJSON struct {
	Filename   string
	Url        string
	Delete_Key string
	Expiry     string
	Size       string
	Sha256sum  string
	Direct_Url string `json:",omitempty"`
}

type RespErrJSON struct {
	Error string
}

var Config struct {
	siteurl   string
	logfile   string
	apikey    string
	apikeycmd string
}

var keys map[string]string

func main() {
	var del bool
	var randomize bool
	var overwrite bool
	var expiry int64
	var deleteKey string
	var accessKey string
	var desiredFileName string
	var configPath string
	var noClipboard bool
	var useSelifURL bool

	flag.BoolVar(&del, "d", false,
		"Delete file at url (ex: -d https://linx.example.com/myphoto.jpg")
	flag.BoolVar(&randomize, "r", false,
		"Randomize filename")
	flag.Int64Var(&expiry, "e", 0,
		"Time in seconds until file expires (ex: -e 600)")
	flag.StringVar(&deleteKey, "deletekey", "",
		"Specify your own delete key for the upload(s) (ex: -deletekey mysecret)")
	flag.StringVar(&accessKey, "accesskey", "",
		"Specify an access key to limit access to the file with a password")
	flag.StringVar(&desiredFileName, "f", "",
		"Specify the desired filename if different from the actual filename or if file from stdin")
	flag.StringVar(&configPath, "c", "",
		"Specify a non-default config path")
	flag.BoolVar(&overwrite, "o", false,
		"Overwrite file (assuming you have its delete key")
	flag.BoolVar(&noClipboard, "no-cb", false,
		"Disable automatic insertion into clipboard")
	flag.BoolVar(&useSelifURL, "selif", false,
		"Return selif url")
	flag.Parse()

	parseConfig(configPath)
	getKeys()

	if del {
		for _, url := range flag.Args() {
			deleteUrl(url)
		}
	} else {
		for _, fileName := range flag.Args() {
			upload(fileName, deleteKey, accessKey, randomize, expiry, overwrite, desiredFileName, noClipboard, useSelifURL)
		}
	}
}

func upload(filePath string, deleteKey string, accessKey string, randomize bool, expiry int64, overwrite bool, desiredFileName string, noClipboard bool, useSelifURL bool) {
	var reader io.Reader
	var fileName string
	var ssum string

	// Need to fetch this before we setup the progress reader
	// as it can conflict with the apikeycmd output
	apikey := getApiKey()

	if filePath == "-" {
		byt, err := ioutil.ReadAll(os.Stdin)
		checkErr(err)

		fileName = desiredFileName

		br := bytes.NewReader(byt)

		ssum = sha256sum(br)
		br.Seek(0, 0)

		reader = progress.NewProgressReader(fileName, br, int64(len(byt)))

	} else {
		fileInfo, err := os.Stat(filePath)
		checkErr(err)
		file, err := os.Open(filePath)
		checkErr(err)

		if desiredFileName == "" {
			fileName = path.Base(file.Name())
		} else {
			fileName = desiredFileName
		}

		br := bufio.NewReader(file)
		ssum = sha256sum(br)
		file.Seek(0, 0)

		reader = progress.NewProgressReader(fileName, br, fileInfo.Size())
	}

	escapedFileName := url.QueryEscape(fileName)

	req, err := http.NewRequest("PUT", Config.siteurl+"upload/"+escapedFileName, reader)
	checkErr(err)

	req.Header.Set("User-Agent", "linx-client")
	req.Header.Set("Accept", "application/json")

	if apikey != "" {
		req.Header.Set("Linx-Api-Key", apikey)
	}
	if deleteKey != "" {
		req.Header.Set("Linx-Delete-Key", deleteKey)
	}
	if accessKey != "" {
		req.Header.Set("Linx-Access-Key", accessKey)
	}
	if randomize {
		req.Header.Set("Linx-Randomize", "yes")
	}
	if expiry != 0 {
		req.Header.Set("Linx-Expiry", strconv.FormatInt(expiry, 10))
	}
	if overwrite {
		fileUrl := Config.siteurl + fileName
		deleteKey, exists := keys[fileUrl]
		if !exists {
			checkErr(errors.New("No delete key for " + fileUrl))
		}

		req.Header.Set("Linx-Delete-Key", deleteKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	if resp.StatusCode == 200 {
		var myResp RespOkJSON
		var returnUrl string

		err := json.Unmarshal(body, &myResp)
		checkErr(err)

		if myResp.Sha256sum != ssum {
			fmt.Println("Warning: sha256sum does not match.")
		}

		if useSelifURL && len(myResp.Direct_Url) != 0 {
			returnUrl = myResp.Direct_Url
		} else {
			returnUrl = myResp.Url
		}

		if noClipboard {
			fmt.Println(returnUrl)
		} else {
			fmt.Printf("Copied %s into clipboard!\n", returnUrl)
			clipboard.WriteAll(returnUrl)
		}

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

	if apikey := getApiKey(); apikey != "" {
		req.Header.Set("Linx-Api-Key", apikey)
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

func getApiKey() string {
	if Config.apikey != "" {
		return Config.apikey
	}

	if Config.apikeycmd == "" {
		return ""
	}

	apikey, err := runCmdFirstLine(Config.apikeycmd)

	if err != nil {
		checkErr(fmt.Errorf("Failed to retrieve API key: %w", err))
	}

	if apikey == "" {
		checkErr(fmt.Errorf("Command did not produce an API key: %s", Config.apikeycmd))
	}

	return apikey
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

func parseConfig(configPath string) {
	var cfgFilePath string

	if configPath == "" {
		cfgFilePath = filepath.Join(getConfigDir(), "linx-client.conf")
	} else {
		cfgFilePath = configPath
	}

	cfg, err := config.LoadOrCreate(cfgFilePath)
	checkErr(err)

	Config.siteurl = cfg.String("siteurl", "")
	Config.logfile = cfg.String("logfile", "")
	Config.apikey = cfg.String("apikey", "")
	Config.apikeycmd = cfg.String("apikeycmd", "")

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

		if Config.apikeycmd == "" && Config.apikey == "" {
			for {
				Config.apikeycmd = getInput("API key retreival command (leave blank for plain token or if instance is public, ex: pass show linx-client)", true)

				if Config.apikeycmd == "" {
					break
				}

				err := validateCommand(Config.apikeycmd)

				if err != nil {
					fmt.Printf("Invalid API key retreival command: %v\n", err)
				} else {
					break
				}
			}
		}
		cfg.SetValue("apikeycmd", Config.apikeycmd)

		if Config.apikey == "" && Config.apikeycmd == "" {
			Config.apikey = getInput("API key (leave blank if instance is public, will be stored in plain text)", true)
		}
		cfg.SetValue("apikey", Config.apikey)

		cfg.Write()

		fmt.Printf("Configuration written at %s\n", cfgFilePath)
	}
}
