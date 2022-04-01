package cmd

import (
	"encoding/json"
	"net/http"

	"github.com/ledongthuc/aws_secrets_storage_sync/cache"
	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/sirupsen/logrus"
)

func startServer(dataSource *cache.SecretLastChanges) {
	mux := http.NewServeMux()
	mux.Handle("/meta/secrets", &server{dataSource: dataSource})
	serverAddr := configs.GetServerAddress()
	logrus.WithFields(logrus.Fields{"addr": serverAddr}).Info("Start API server")
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		panic(err)
	}
}

type server struct {
	dataSource *cache.SecretLastChanges
}

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("only support method GET"))
		return
	}

	data, err := json.Marshal(h.dataSource.All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("compose response: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
