package main

import (
	"fmt"
	"github.com/juxuny/yc/cmd"
	"github.com/juxuny/yc/log"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

type server struct {
}

func (s *server) Prepare(cmd *cobra.Command) {
}

func (s *server) InitFlag(cmd *cobra.Command) {
}

func (s *server) Run() {
	log.Info("register")
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		//log.Info(r.RequestURI)
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadGateway)
			_, err = w.Write([]byte(fmt.Sprintf("%v", err)))
			if err != nil {
				log.Error(err)
			}
		}
		w.WriteHeader(http.StatusOK)
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 2)
			_, err = w.Write([]byte(fmt.Sprintf("%d\n", i)))
			if err != nil {
				log.Error(err)
			}
		}
	})
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Error(err)
	}
}

func main() {
	runCommand := cmd.NewCommandBuilder("", &server{})
	if err := runCommand.Build().Execute(); err != nil {
		log.Error(err)
	}
}
