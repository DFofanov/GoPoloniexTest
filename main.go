// The task of obtaining data on completed transactions from the Poloniex exchange using Websocket
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/droundy/goopt"
	"github.com/recws-org/recws"
)

type (
	Poloniex struct {
		Key           string
		Secret        string
		ws            recws.RecConn
		subscriptions map[string]string
	}
)

type (
	Subscription struct {
		Command string `json:"command"`
		Channel string `json:"channel"`
	}
)

type (
	RecentTrade struct {
		ID        int64
		Pair      string
		Side      string
		Price     float64
		Amount    float64
		Event     string
		Timestamp time.Time
	}
)

type Config struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

const (
	POLONIEX_URL = "wss://api2.poloniex.com/"
)

func readConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetConfig(path string) (*Config, error) {
	return readConfig(path)
}

func toString(i interface{}) string {
	switch i := i.(type) {
	case string:
		return i
	case float64:
		return fmt.Sprintf("%.8f", i)
	case int64:
		return fmt.Sprintf("%d", i)
	case json.Number:
		return i.String()
	}
	return ""
}

func toFloat(i interface{}) float64 {
	maxFloat := float64(math.MaxFloat64)
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return maxFloat
		}
		return a
	case float64:
		return i
	case int64:
		return float64(i)
	case json.Number:
		a, err := i.Float64()
		if err != nil {
			return maxFloat
		}
		return a
	}
	return maxFloat
}

func GoPoloniexTest(key, secret string) {
	p := &Poloniex{}
	p.Key = key
	p.Secret = secret
	p.ws = recws.RecConn{}
	p.ws.Dial(POLONIEX_URL, http.Header{})

	p.subscriptions = map[string]string{"121": "BTC_USDT", "256": "TRX_USDT", "149": "ETH_USDT"}

	defer func() {
		for id, _ := range p.subscriptions {
			unsubscribe := Subscription{Command: "unsubscribe", Channel: id}
			p.ws.WriteJSON(unsubscribe)
		}
		p.ws.Close()
	}()

	for id, _ := range p.subscriptions {
		subscribe := Subscription{Command: "subscribe", Channel: id}
		p.ws.WriteJSON(subscribe)
	}

	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			go p.ws.Close()
			log.Printf("Websocket closed %s", p.ws.GetURL())
			return
		default:
			message := []interface{}{}
			if err := p.ws.ReadJSON(&message); err != nil {
				log.Println(err)
				continue
			}
			chid := int64(message[0].(float64))
			chids := toString(chid)
			if _, ok := p.subscriptions[chids]; ok {
				marketID := int64(toFloat(message[0]))
				for _, _v := range message[2].([]interface{}) {
					v := _v.([]interface{})
					trade := RecentTrade{}
					trade.Pair = p.subscriptions[toString(marketID)]
					if v[0].(string) == "t" {
						trade.Event = "trade"
						trade.ID = int64(toFloat(message[1]))
						trade.Side = "sell"
						if t := toFloat(v[2]); t == 1.0 {
							trade.Side = "buy"
						}
						trade.Price = toFloat(v[3])
						trade.Amount = toFloat(v[4])
						t := time.Unix(int64(toFloat(v[5])), 0)
						trade.Timestamp = t

						log.Println("============================================================")
						log.Println(fmt.Sprintf("ID: %v", trade.ID))
						log.Println(fmt.Sprintf("Pair: %s", trade.Pair))
						log.Println(fmt.Sprintf("Price: %v", trade.Price))
						log.Println(fmt.Sprintf("Amount: %v", trade.Amount))
						log.Println(fmt.Sprintf("Side: %s", trade.Side))
						log.Println(fmt.Sprintf("Timestamp: %s", trade.Timestamp))
					}
				}
			}
		}
	}
}

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

func Version() error {
	fmt.Printf("GoPoloniexTest 0.1 %s\n\nCopyright (C) 2021 %s\n%s\n", goopt.Version, goopt.Author, License)
	os.Exit(0)
	return nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, goopt.Usage())
	os.Exit(1)
}

func main() {
	goopt.Author = "Dmitry Fofanov"
	goopt.Version = "0.1"
	goopt.Summary = "The task of obtaining data on completed transactions from the Poloniex exchange using Websocket"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s configFile\n or:\t%s OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", Version)
	goopt.Parse(nil)
	if len(goopt.Args) != 2 {
		PrintUsage()
	}
	config, _ := GetConfig(goopt.Args[0])
	GoPoloniexTest(config.Key, config.Secret)
}
