package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

// This does not create the file if it doesn't exist
func GetTodoFromDataFile(dataFileFullPath string) ([]entity.TodoEntity, error) {
	file, errReadFile := os.ReadFile(dataFileFullPath)
	if errReadFile != nil {
		return []entity.TodoEntity{}, fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []entity.TodoEntity
	if errUnmarshal := yaml.Unmarshal(file, &items); errUnmarshal != nil {
		return []entity.TodoEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil
}

func WriteTodoToFile(fileFullPath string) error {
	return nil
}

func GetDataFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(constant.DATA_FOLDER_KEY))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + constant.DEFAULT_DATA_PATH_AFTER_HOME
	}

	if errMkDirAll := os.MkdirAll(dataFolder, 0770); errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return dataFolder, nil
}

func PrettyPrintDataFile(dataFileFullPath string, prefix string) error {
	todoList, errGetTodoFromDataFile := GetTodoFromDataFile(dataFileFullPath)
	if errGetTodoFromDataFile != nil {
		return fmt.Errorf("GetTodoFromDataFile: %w", errGetTodoFromDataFile)
	}

	for _, todo := range todoList {
		fmt.Println(prefix + todo.Text)
	}
	return nil
}
