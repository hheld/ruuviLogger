package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v4/pgxpool"
    "ruuviLogger"
)

type sensorDb struct {
    *pgxpool.Pool
}

func ConnectToDb() (*sensorDb, error) {
    cfg := ruuviLogger.Cfg

    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.DbUser,
        cfg.DbPwd,
        cfg.DbHost,
        cfg.DbPort,
        cfg.DbName)

    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, dbURL)
    if err != nil {
        return nil, err
    }

    return &sensorDb{pool}, nil
}
