package cli

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	APP_NAME          string = "installer"
	APP_VERSION       string = "0.1.0"
	AUTHOR            string = "datnn"
	EMAIL             string = "datnn288@gmail.com"
	LICENSE           string = "Apache"
	SHORT_DESCRIPTION string = "IDM write in Go"
	LONG_DESCRIPTION  string = `GoIDM is a CLI application for handle download file multithread. That's all`
)

var rootCmd = &cobra.Command{
	Use:   APP_NAME,
	Short: SHORT_DESCRIPTION,
	Long:  LONG_DESCRIPTION,
}
var cfgFile string = ""

func Run() error {
	return rootCmd.Execute()
}
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		var homeDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatal("Can't find home dir")
		}
		if runtime.GOOS == "windows" {
			viper.AddConfigPath(fmt.Sprintf("%s\\.%s", homeDir, APP_NAME))
		} else {
			viper.AddConfigPath(fmt.Sprintf("%s/.config", homeDir))
		}
		viper.SetConfigFile("yaml")
		viper.SetConfigName(fmt.Sprintf(".%s", APP_NAME))
	}
	viper.AutomaticEnv()

}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", fmt.Sprintf("%s <%s>", AUTHOR, EMAIL))
	viper.SetDefault("license", "apache")
}
