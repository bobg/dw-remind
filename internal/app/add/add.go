package add

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/DanWlker/remind/internal/pkg/data"
	"github.com/DanWlker/remind/internal/pkg/record"
	"github.com/DanWlker/remind/internal/pkg/shared"
)

func AddRun(globalFlag bool, args []string) error {
	if globalFlag {
		errAddTodoAndAssociateTo := addTodoAndAssociateTo("", args)
		if errAddTodoAndAssociateTo != nil {
			return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
		}
		return nil
	}

	homeRemCurrProExDir, errHomeRemCurrProExDir := shared.GetHomeRemovedWorkingDir()
	if errHomeRemCurrProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errHomeRemCurrProExDir)
	}

	if errAddTodoAndAssociateTo := addTodoAndAssociateTo(homeRemCurrProExDir, args); errAddTodoAndAssociateTo != nil {
		return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
	}

	return nil
}

func addTodoAndAssociateTo(directory string, todoListString []string) error {
	// Find the record in the record file
	recordItems, errGetRecordFileContents := record.GetFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	idx := slices.IndexFunc(recordItems, func(item record.RecordEntity) bool {
		return item.Path == directory
	})

	dataFolder, errGetDataFolder := data.GetFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	var currentDirectoryRecord *record.RecordEntity
	if idx == -1 {
		tempCurrentDirectoryRecord, errCreateNewRecord := record.CreateNewRecord(directory)
		if errCreateNewRecord != nil {
			return fmt.Errorf("helper.CreateNewRecord: %w", errCreateNewRecord)
		}
		currentDirectoryRecord = &tempCurrentDirectoryRecord
		recordItems = append(recordItems, *currentDirectoryRecord)
		if err := record.SetFileContents(recordItems); err != nil {
			return fmt.Errorf("helper.SetRecordFileContents: %w", err)
		}
	} else {
		currentDirectoryRecord = &recordItems[idx]
	}

	// Read the file, it will exist if it reaches here
	dataFileFullPath := dataFolder + string(os.PathSeparator) + currentDirectoryRecord.DataFileName
	_, errStat := os.Stat(dataFileFullPath)

	if errors.Is(errStat, os.ErrNotExist) {
		return fmt.Errorf("You fcked up, os.Stat: %w", errStat) // This should never occur
	} else if errStat != nil {
		return fmt.Errorf("os.Stat: %w", errStat)
	}

	todoList, errReadFromFile := data.GetTodoFromFile(dataFileFullPath)
	if errReadFromFile != nil {
		return fmt.Errorf("helper.ReadFromFile: %w", errReadFromFile)
	}

	for _, item := range todoListString {
		todoList = append(todoList, data.TodoEntity{Text: item})
	}

	errWriteTodoToFile := data.WriteTodoToFile(dataFileFullPath, todoList)
	if errWriteTodoToFile != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", errWriteTodoToFile)
	}
	return nil
}
