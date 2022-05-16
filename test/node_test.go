package test

import (
	"bytes"
	"demoBlockChain/database"
	"demoBlockChain/node"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const url = "http://localhost:" + node.Port

func init() {
	go func() {
		//check if the node is running
		_, err := http.Get(url)
		if err == nil {
			log.Fatal("node is running")
		}

		//start the node
		err = node.Run("./_testdata")
		if err != nil {
			log.Fatal("ini node test, error: ", err)
		}
	}()

	time.Sleep(time.Second * 2)
}

func getBodyContent(resp *http.Response, t *testing.T) []byte {
	//get content of resp
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("close body error: ", err)
		}
	}(resp.Body)
	body := resp.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error("read body content error: ", err)
	}
	return content
}

func TestRun(t *testing.T) {

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Status code want 200, but got ", resp.StatusCode)
	}
}

func TestPostTxAdd(t *testing.T) {
	tx1 := database.NewTx("zhouyh", "li", 10, "")
	tx2 := database.NewTx("li", "zhouyh", 2, "")

	txList := []database.Tx{tx1, tx2}
	j, err := json.Marshal(txList)
	if err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(j)

	response, err := http.Post(url+"/tx/add", "application/json", reader)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Status code want 200, but got ", response.StatusCode)
	}
}

func TestGetBlocks(t *testing.T) {
	resp, err := http.Get(url + "/blocks")
	if err != nil {
		t.Error("get blocks error: ", err)
	}

	content := getBodyContent(resp, t)
	var blocks []*database.BlockFS
	err = json.Unmarshal(content, &blocks)
	if err != nil {
		t.Error("parse blocks error: ", err)
	}
}
