/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/DanWlker/remind/entity"
	"github.com/DanWlker/remind/helper"
	"github.com/spf13/cobra"
)

var globalFlag_add = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Adds the todos to a global todo list",
	},
	Value: false,
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a todo with directory context",
	Long:  "Adds a todo with directory context, by default it will associate the todo with the local directory. Use -g to bind it to the global $HOME todo list",
	Run: func(cmd *cobra.Command, args []string) {
		globalFlag, errGetBool := cmd.Flags().GetBool(globalFlag_add.Name)
		if errGetBool != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool))
		}
		errAddRun := addRun(globalFlag, args)
		if errAddRun != nil {
			cobra.CheckErr(fmt.Errorf("addRun: %w", errAddRun))
		}
	},
}

func addTodoAndAssociateTo(directory string, todoListString []string) error {
	// Find the record in the record file
	recordItems, errGetRecordFileContents := helper.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	idx := slices.IndexFunc(recordItems, func(item entity.ProjectRecordEntity) bool {
		return item.Path == directory
	})

	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	var currentDirectoryRecord *entity.ProjectRecordEntity
	if idx == -1 {
		tempCurrentDirectoryRecord, errCreateNewRecord := helper.CreateNewRecord(directory)
		if errCreateNewRecord != nil {
			return fmt.Errorf("helper.CreateNewRecord: %w", errCreateNewRecord)
		}
		currentDirectoryRecord = &tempCurrentDirectoryRecord
		recordItems = append(recordItems, *currentDirectoryRecord)
		helper.SetRecordFileContents(recordItems)
	} else {
		currentDirectoryRecord = &recordItems[idx]
	}

	// Read the file, it will exist if it reaches here
	dataFileFullPath := dataFolder + string(os.PathSeparator) + currentDirectoryRecord.DataFileName
	_, errStat := os.Stat(dataFileFullPath)

	var todoList []entity.TodoEntity
	if errStat == nil {
		var errReadFromFile error
		todoList, errReadFromFile = helper.GetTodoFromDataFile(dataFileFullPath)
		if errReadFromFile != nil {
			return fmt.Errorf("helper.ReadFromFile: %w", errReadFromFile)
		}
	} else if errors.Is(errStat, os.ErrNotExist) {
		return fmt.Errorf("You fcked up, os.Stat: %w", errStat) // This should never occur
	} else {
		return fmt.Errorf("os.Stat: %w", errStat)
	}

	for _, item := range todoListString {
		todoList = append(todoList, entity.TodoEntity{Text: item})
	}

	errWriteTodoToFile := helper.WriteTodoToFile(dataFileFullPath, todoList)
	if errWriteTodoToFile != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", errWriteTodoToFile)
	}
	return nil
}

func addRun(globalFlag bool, args []string) error {
	if globalFlag {
		errAddTodoAndAssociateTo := addTodoAndAssociateTo("", args)
		if errAddTodoAndAssociateTo != nil {
			return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
		}
		return nil
	}

	homeRemCurrProExDir, errHomeRemCurrProExDir := helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	if errHomeRemCurrProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errHomeRemCurrProExDir)
	}

	if errAddTodoAndAssociateTo := addTodoAndAssociateTo(homeRemCurrProExDir, args); errAddTodoAndAssociateTo != nil {
		return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(globalFlag_add.Name, globalFlag_add.Shorthand, globalFlag_add.Value, globalFlag_add.Usage)
}
