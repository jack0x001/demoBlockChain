package database

import (
	"os"
	"path"
)

func GetGenesisFilePath(dataDir string) string {
	return path.Join(dataDir, "genesis.json")
}

func GetBlockFilePath(dataDir string) string {
	return path.Join(dataDir, "block.db")
}

func initDataDirIfNotExists(dataDir string) error {

	//if dataDir does not exist, create it
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//if genesis.json does not exist, create it
	genesisFilePath := GetGenesisFilePath(dataDir)
	if _, err := os.Stat(genesisFilePath); os.IsNotExist(err) {
		_, err = os.Create(genesisFilePath)
		if err != nil {
			return err
		}

		err = writeGenesis(genesisFilePath)
		if err != nil {
			return err
		}
	}

	//if block.db does not exist, create it
	blockFilePath := GetBlockFilePath(dataDir)
	if _, err := os.Stat(blockFilePath); os.IsNotExist(err) {
		_, err = os.Create(blockFilePath)
		if err != nil {
			return err
		}
		err = os.WriteFile(blockFilePath, []byte(""), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
