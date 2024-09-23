package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitializeViper
/*
InitializeViper inicializa viper para leer el archivo config.yml 
y las variables de entorno en la aplicación.
*/
func InitializeViper() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}