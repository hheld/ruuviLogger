package db

import "context"

func (db *sensorDb) AddRuuvitag(name, address string) error {
    query := `INSERT INTO ruuvitag
                (name, address)
              VALUES
                ($1, $2)
              ON CONFLICT DO NOTHING;`

    _, err := db.Pool.Exec(context.Background(), query, name, address)

    return err
}

func (db *sensorDb) GetRuuvitagID(address string) (int, error) {
    query := "SELECT id FROM ruuvitag WHERE address=$1"

    row := db.Pool.QueryRow(context.Background(), query, address)

    var sensorID int
    err := row.Scan(&sensorID)

    return sensorID, err
}
