package soldier

import (
	"fmt"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/modules/pkg/colonyutil"
	"github.com/erda-project/erda/modules/soldier/command"
	"github.com/erda-project/erda/pkg/httpclient"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	logrus.Infof("start----------------->")
	logrus.SetLevel(logrus.DebugLevel)
	router := mux.NewRouter()
	router.Use(colonyutil.LoggingMiddleware)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/terminal", command.Terminal)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func TestServer2(t *testing.T) {
	b := bundle.New(
		bundle.WithHTTPClient(httpclient.New(httpclient.WithTimeout(time.Second, time.Second))),
	)
	err := b.MySQLCheck(&apistructs.MysqlExec{
		URL:      "jdbc:mysql://mysql-slave.group-addon-mysql--xb4a3f249fe7247ea9ca817fc2ee83e19.svc.cluster.local:3306",
		User:     "root",
		Password: "fXHYLkgm6S0k21mG",
	}, "ops.default.svc.cluster.local:9027")
	fmt.Println(err)
}
