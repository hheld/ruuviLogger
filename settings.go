package ruuviLogger

import (
    "encoding/json"
    "io"
    "log"
    "os"
)

const settingsFileName = "ruuviLogger.json"

type settings struct {
    Ruuvitags []ruuvitagInfo
}

type ruuvitagInfo struct {
    Name    string `json:"name"`
    Address string `json:"address"`
}

var Settings settings

func init() {
    settingsFile, err := os.Open(settingsFileName)
    if err != nil {
        log.Fatalf("could not open file '%s'; without that there is no point of running this application: %s", settingsFileName, err)
    }
    defer settingsFile.Close()

    settingsData, err := io.ReadAll(settingsFile)
    if err != nil {
        log.Fatalf("error reading settings from '%s': %s", settingsFileName, err)
    }

    if err := json.Unmarshal(settingsData, &Settings); err != nil {
        log.Fatalf("could not interpret data from '%s': %s", settingsFileName, err)
    }

    log.Printf("using these settings: %+v", Settings)
}
