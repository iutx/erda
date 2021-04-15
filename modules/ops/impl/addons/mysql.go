// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package addons

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
)

// InitMysql execute init sql in addon mysql
func (a *Addons) InitMysql(reqs []apistructs.OpsAddonMysqlInitRequest) error {
	logrus.Infoln("mysql init:", reqs)
	for _, req := range reqs {
		dsn := GetDSN(req.URL, req.User, req.Password)
		if err := Exec(dsn, req.SQLs); err != nil {
			logrus.Errorf("exec sql error: %v", err)
			return err
		}
	}
	return nil
}

// CheckMysql check master and slave status of addon mysql
func (a *Addons) CheckMysql(req apistructs.OpsAddonMysqlCheckRequest) (map[string]interface{}, error) {
	const sqlSlaveStatus = "show slave status"

	logrus.Infoln("mysql check:", req)
	db, err := sql.Open("mysql", GetDSN(req.URL, req.User, req.Password))
	if err != nil {
		logrus.Errorf("open connect error: %v", err)
		return nil, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Errorf("close db connection error: %v", err)
		}
	}()

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	rows, err := db.QueryContext(ctx, sqlSlaveStatus)
	if err != nil {
		logrus.Errorf("query %s error: %v", sqlSlaveStatus, err)
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("no rows")
	}
	cols, err := rows.Columns()
	if err != nil {
		logrus.Errorf("get row columns error: %v", err)
		return nil, err
	}
	status := make([]interface{}, len(cols))
	for i := range status {
		if cols[i] == "Slave_IO_Running" {
			status[i] = new(string)
		} else if cols[i] == "Slave_SQL_Running" {
			status[i] = new(string)
		} else {
			status[i] = new(interface{})
		}
	}
	err = rows.Scan(status...)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	m := make(map[string]interface{}, 2)
	for i, col := range cols {
		if col == "Slave_IO_Running" {
			m["slaveIoRunning"] = *(status[i].(*string))
		} else if col == "Slave_SQL_Running" {
			m["slaveSqlRunning"] = *(status[i].(*string))
		}
	}
	return m, nil
}

// ExecMysql exec sql
func (a *Addons) ExecMysql(req apistructs.OpsAddonMysqlExecRequest) error {
	logrus.Infoln("mysql exec:", req)
	dsn := GetDSN(req.URL, req.User, req.Password)
	if err := Exec(dsn, req.SQLs); err != nil {
		logrus.Errorf("exec sql error: %v", err)
		return err
	}
	return nil
}

// ExecFileMysql exec sql file
func (a *Addons) ExecFileMysql(req apistructs.OpsAddonMysqlExecFileRequest) error {
	logrus.Infoln("mysql exec file:", req)
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	host, port := ParseAddr(req.URL)
	c := exec.CommandContext(ctx, "bash",
		"/app/sql.sh", host, port, req.User, req.Password, strings.Join(req.CreateDbs, "\n"), req.OssURL)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = []string{"NETDATA_SQLDATA_PATH=" + os.Getenv("NETDATA_SQLDATA_PATH")}
	if err := c.Run(); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

// Exec exec sql
func Exec(dsn string, scripts []string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Errorf("close db connection error: %v", err)
		}
	}()
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	for _, script := range scripts {
		_, err = db.ExecContext(ctx, script)
		logrus.Infof("exec sql: %s", script)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseAddr parse host and port from request url
func ParseAddr(url string) (string, string) {
	host, port, _ := net.SplitHostPort(strings.Replace(url, "jdbc:mysql://", "", -1))
	return host, port
}

// GetDSN get dns info
func GetDSN(url, user, password string) string {
	host, port := ParseAddr(url)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4,utf8&parseTime=true", user, password, host, port)
}
