/*
Copyright © 2024 DanWlker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DanWlker/remind/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "remind",
	Short: "Remind is a project aware todo app",
	Long: `Remind is a project aware todo app that will show relevant todos
	depending on the project folder. It also stores all notes in
	$HOME/remind/ to allow for easy syncing`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFolder, errGetConfigFolder := config.GetConfigFolder()
	if errGetConfigFolder != nil {
		cobra.CheckErr(fmt.Errorf("roomCmd: initConfig: helper.GetConfigFolder: %w", errGetConfigFolder))
	}

	viper.AddConfigPath(configFolder)
	viper.SetConfigType(config.DEFAULT_CONFIG_FILE_TYPE)
	viper.SetConfigName(config.DEFAULT_CONFIG_FILE_NAME)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
