package main

import (
    /*
    */
    "fmt"
    "time"
    "log"
    tb "gopkg.in/tucnak/telebot.v2"
    "github.com/tkanos/gonfig"
    md "charon/src/Models"
)

func main(){
    cfg := loadConfig();
    fmt.Println(cfg.Token);
    runBot(cfg);
}

func runBot(cfg md.Configuration){
    b, err := tb.NewBot(tb.Settings{
        Token:  cfg.Token,
        Poller: &tb.LongPoller{Timeout: 10 * time.Second},
    })	
    if err != nil {
        log.Fatal(err)
        return
    }

    b.Handle("/hello", func(m *tb.Message) {
        b.Send(m.Sender, "hello world")
    })

    b.Start()

}

func loadConfig() md.Configuration{
    configuration := md.Configuration{}
    err := gonfig.GetConf("config.json", &configuration)
    if err != nil {
        panic(err)
    }
    return configuration;
}
