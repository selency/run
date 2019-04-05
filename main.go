package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	"github.com/bluele/slack"
)

var (
	slackClient *slack.Slack
)

func init() {
	viper.SetConfigName("run-config")
	viper.AddConfigPath("$HOME/.run")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %s", err))
	}

	slackClient = slack.New(viper.Get("slack.token").(string))
}

func sendMessage(message string) {
	slackClient.ChatPostMessage(viper.Get("slack.channel").(string), message, &slack.ChatPostMessageOpt{
		Username:  viper.Get("slack.username").(string),
		IconEmoji: viper.Get("slack.emoji").(string),
	})
}

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	sendMessage(fmt.Sprintf(
		":hourglass_flowing_sand: running command: `%s %s`",
		os.Args[1],
		strings.Join(os.Args[2:], " "),
	))

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	cmd.Run()

	if cmd.ProcessState.Success() {
		sendMessage(":heavy_check_mark: everything went well")
	} else {
		sendMessage(fmt.Sprintf(":x: command exited with non-zero status: `%s`", cmd.ProcessState))
		slackClient.FilesUpload(&slack.FilesUploadOpt{
			Channels: []string{viper.Get("slack.channel").(string)},
			Content:  stderr.String(),
			Filename: "stderr",
		})
	}
}
