package ruuviLogger

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type dbConfig struct {
    DbHost               string
    DbPort               int
    DbUser               string
    DbPwd                string
    DbName               string
    DbWriteIntervalInSec int
}

type config struct {
    dbConfig
}

func defaultConfig() config {
    return config{
        dbConfig{
            DbHost:               "localhost",
            DbPort:               5432,
            DbUser:               "hheld",
            DbPwd:                "hheld-pwd",
            DbName:               "weatherdb",
            DbWriteIntervalInSec: 60,
        },
    }
}

var Cfg config

func init() {
    if err := godotenv.Load(); err == nil {
        log.Printf("loading settings from .env file")
    } else {
        Cfg = defaultConfig()
    }

    Cfg.DbPwd = os.Getenv("DB_PASSWORD")
    Cfg.DbUser = os.Getenv("DB_USER")
    Cfg.DbHost = os.Getenv("DB_HOST")
    Cfg.DbName = os.Getenv("DB_NAME")

    if writeInterval, err := strconv.Atoi(os.Getenv("DB_WRITE_INTERVAL_SEC")); err != nil {
        log.Print("could not determine the DB writing interval from environment variable 'DB_WRITE_INTERVAL_SEC', using the default one")
    } else {
        Cfg.DbWriteIntervalInSec = writeInterval
    }

    if port, err := strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
        log.Print("could not determine the port from enviroment variable 'DB_PORT', using the default one")
    } else {
        Cfg.DbPort = port
    }

    log.Printf("using database %s:%d/%s", Cfg.DbHost, Cfg.DbPort, Cfg.DbName)
}
