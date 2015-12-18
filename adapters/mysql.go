package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlAdapter struct {
	Adapter
	name             string
	serviceConfig    *config.ServiceConfig
	connectionString string
	logger           lager.Logger
	db               *sql.DB
	tableName        string
}

type mysqlConfig struct {
	Credentials credentialsConfig `json:"credentials"`
}

type credentialsConfig struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	DbName   string `json:"name"`
}

func NewMysqlAdapter(serviceConfig *config.ServiceConfig, logger lager.Logger) *MysqlAdapter {
	return &MysqlAdapter{
		name:          serviceConfig.Name,
		serviceConfig: serviceConfig,
		logger:        logger,
	}
}

func (a MysqlAdapter) Name() string {
	return a.name
}

func (a *MysqlAdapter) Setup() error {
	var err error
	a.logger.Info("Setting up mysql adapter...")

	var configProperties mysqlConfig
	if err := json.Unmarshal(a.serviceConfig.ConfigJSON, &configProperties); err != nil {
		return fmt.Errorf("Failed to parse config for %s: %s", a.Name(), err.Error())
	}

	creds := configProperties.Credentials
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		creds.Username,
		creds.Password,
		creds.Hostname,
		creds.Port,
		creds.DbName,
	)

	a.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return fmt.Errorf("Failed to open mysql connection: %s", err.Error())
	}

	a.tableName = fmt.Sprintf("service_canary_%s", a.Name())
	_, err = a.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`( id INT )", a.tableName))
	if err != nil {
		return fmt.Errorf("Failed to create mysql table: %s", err.Error())
	}

	return nil
}

func (a *MysqlAdapter) Run() error {
	a.logger.Info("Running mysql adapter...")

	id := time.Now().Unix()
	_, err := a.db.Exec(fmt.Sprintf("INSERT INTO `%s` (id) VALUES (?)", a.tableName), id)
	if err != nil {
		return fmt.Errorf("Failed to INSERT for %s: %s", a.Name(), err.Error())
	}

	var rowCount int
	err = a.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE id=?", a.tableName), id).Scan(&rowCount)
	if err != nil {
		return fmt.Errorf("Failed to SELECT COUNT for %s: %s", a.Name(), err.Error())
	}

	if rowCount < 1 {
		return fmt.Errorf("Could not find inserted row for %s", a.Name())
	}

	//cleanup
	_, err = a.db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", a.tableName))
	if err != nil {
		return fmt.Errorf("Failed to TRUNCATE for %s: %s", a.Name(), err.Error())
	}

	return nil
}
