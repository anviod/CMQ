package main

import (
    "github.com/spf13/viper"
    _ "github.com/heralight/logrus_mate/hooks/file"
    "github.com/sirupsen/logrus"

    "fmt"
    "github.com/tian-yuan/CMQ/util"
    "runtime"
	"github.com/tian-yuan/CMQ/message-dispatcher/commands"
)

func initLogger() {

	// ########## Init Viper
	var viper = viper.New()

	viper.SetConfigName("mate") // name of config file (without extension), here we use some logrus_mate sample
	viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath("./conf")               // optionally look for config in the working directory
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// ########### End Init Viper

	// Read configuration
	var c = util.UnmarshalConfiguration(viper) // Unmarshal configuration from Viper
	util.SetConfig(logrus.StandardLogger(), c) // for e.g. apply it to logrus default instance

	// ### End Read Configuration
}

func main() {
	initLogger()
	runtime.GOMAXPROCS(runtime.NumCPU())

	commands.Execute()

	stopCh := util.SetupSignalHandler()
	<-stopCh

	logrus.Infof("message dispatcher stop")
	commands.Stop()
}
