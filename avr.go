package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

// Setup various global variables (should be moved to config JSON later)
var logForward = true
var logToFile = false

// Serverconfig construct to parse data from a local JSON file.
type Serverconfig struct {
	Avrport                     string `json:"avrport"`
	Avrsecureport               string `json:"avrsecureport"`
	Avrtransport                string `json:"avrtransport"` //options are: mqttforwarder or file
	Avrexportfolder             string `json:"avrexportfolder"`
	Avrlogfile                  string `json:"avrlogfile"`
	Mqtttestlocal               bool   `json:"mqtttestlocal"`
	Mqttcredentialsneeded       bool   `json:"mqttcredentialsneeded"`
	Mqttremoteservercredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mqttservercredentials"`
	Mqttclientname     string `json:"mqttclientname"`
	Mqttlocalhost      string `json:"mqttlocalhost"` // when connecting to mosquitto on mac use docker.for.mac.localhost
	Mqttlocalport      string `json:"mqttlocalport"`
	Mqtthost           string `json:"mqtthost"`
	Mqttport           string `json:"mqttport"`
	Mqttpublishtopic   string `json:"mqttpublishtopic"`
	Mqttsubscribetopic string `json:"mqttsubscribetopic"`
}

// define general server config and variables
var serverconfiguration = Serverconfig{}
var avrTransport = serverconfiguration.Avrtransport
var testWithLocalMqtt = serverconfiguration.Mqtttestlocal
var mqttCredentialsNeeded = serverconfiguration.Mqttcredentialsneeded
var mqttRemoteUsername = serverconfiguration.Mqttremoteservercredentials.Username
var mqttRemotePassword = serverconfiguration.Mqttremoteservercredentials.Password
var mqttClientName = serverconfiguration.Mqttclientname
var mqttLocalhost = serverconfiguration.Mqttlocalhost
var mqttLocalPort = serverconfiguration.Mqttlocalport
var mqttRemoteHost = serverconfiguration.Mqtthost
var mqttRemotePort = serverconfiguration.Mqttport
var mqttPublishTopic = serverconfiguration.Mqttpublishtopic
var mqttSubscribeTopic = serverconfiguration.Mqttsubscribetopic

func retrieveConfig() {
	// Launch server config
	openconfig, err := os.Open("./config/serverConfig.json")
	if err != nil {
		logger.Error().Err(err).Msg("Error while opening config file from primary location, retrying location")
		openconfig, err := os.Open("/config/serverConfig.json")
		defer openconfig.Close()
		decoder := json.NewDecoder(openconfig)
		decodingerror := decoder.Decode(&serverconfiguration)
		if decodingerror != nil {
			logger.Error().Err(err).Msg("Error while decoding config file from retry")
		}
		if err != nil {
			logger.Error().Err(err).Msg("Fallback location also failed.")
		} else {
			logger.Info().Msg("Opened serverconfiguration from secondary location")
		}
	} else {
		logger.Info().Msg("Opened serverconfiguration from primary location")
	}
	defer openconfig.Close()
	decoder := json.NewDecoder(openconfig)
	decodingerror := decoder.Decode(&serverconfiguration)
	if decodingerror != nil {
		logger.Error().Err(err).Msg("Error while decoding config file from primary location.")
	}
	avrTransport = serverconfiguration.Avrtransport
	logger.Trace().Msg(avrTransport)
	testWithLocalMqtt = serverconfiguration.Mqtttestlocal
	logger.Trace().Msgf("Testlocal is %t", testWithLocalMqtt)
	mqttCredentialsNeeded = serverconfiguration.Mqttcredentialsneeded
	logger.Trace().Msgf("Credentials needed is %t", mqttCredentialsNeeded)
	mqttRemoteUsername = serverconfiguration.Mqttremoteservercredentials.Username
	logger.Trace().Msg(mqttRemoteUsername)
	mqttRemotePassword = serverconfiguration.Mqttremoteservercredentials.Password
	logger.Trace().Msg(mqttRemotePassword)
	mqttClientName = serverconfiguration.Mqttclientname
	logger.Trace().Msg(mqttClientName)
	mqttLocalhost = serverconfiguration.Mqttlocalhost
	logger.Trace().Msg(mqttLocalhost)
	mqttLocalPort = serverconfiguration.Mqttlocalport
	logger.Trace().Msg(mqttLocalPort)
	mqttRemoteHost = serverconfiguration.Mqtthost
	logger.Trace().Msg(mqttRemoteHost)
	mqttRemotePort = serverconfiguration.Mqttport
	logger.Trace().Msg(mqttRemotePort)
	mqttPublishTopic = serverconfiguration.Mqttpublishtopic
	logger.Trace().Msg(mqttPublishTopic)
	mqttSubscribeTopic = serverconfiguration.Mqttsubscribetopic
	logger.Trace().Msg(mqttSubscribeTopic)
}

func main() {
	// argument function
	arguments := os.Args

	// initiate zerologger for all logger functions
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if len(arguments) > 2 {
		switch arguments[2] {
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
			fmt.Println("Logging set to:", arguments[2])
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			fmt.Println("Logging set to:", arguments[2])
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
			fmt.Println("Logging set to:", arguments[2])
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			fmt.Println("Logging set to:", arguments[2])
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			fmt.Println("Logging set to:", arguments[2])
		case "trace":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
			fmt.Println("Logging set to:", arguments[2])
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			fmt.Println("Unrecognized loglevel, logging set to: info ")
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		fmt.Println("No loglevel provided, default logging set to: info ")
	}
	// Retrieve serverconfig
	retrieveConfig()

	// Create unique ID to be able to generate unique logfile names
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error while generating unique id")
	}
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	logger.Info().Msg("Unique id of AVR server generated: " + uuid)

	// Create naming conventions and folders for log files
	loggingBaseURL := serverconfiguration.Avrlogfile
	loggingExportFolder := serverconfiguration.Avrexportfolder
	loggingPath := loggingBaseURL + uuid + "/"
	logFileName := loggingPath + "avr-prod-go" + uuid + ".log"

	if avrTransport == "file" {
		// Check if folder to logfile location exists and if not create it
		if _, err := os.Stat(loggingPath); os.IsNotExist(err) {
			err := os.MkdirAll(loggingPath, 0755)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create logfile directory")
			} else {
				logger.Info().Msg("Creating logfile directory at: " + loggingPath)
				dir, err := os.Getwd()
				if err != nil {
					logger.Error().Err(err).Msg("Failed to get current folder")
				}
				logger.Debug().Msg("Working directory is: " + dir)
			}
		}
		logger.Info().Msg("LogFile path and file created: " + logFileName)
	}

	// Launch server based on input paramters during launch.
	if len(arguments) == 1 || avrTransport == "" {
		logger.Fatal().Msg("Missing parameters, please provide port number or transportmechanisms for launching the server")
		return
	} else if arguments[1] == "2498" {
		logger.Info().Msgf("Starting non-secure AVR server")
		if avrTransport == "file" {
			go createRotatingLogger(logFileName, loggingPath)
			go createExportFile(logFileName, loggingPath, loggingExportFolder)
		} else {
			logger.Info().Msg("Not logging to file, skipping generation of exportfiles or logrotation")
		}
		startAvr()
	} else if arguments[1] == "2499" {
		logger.Info().Msgf("Starting secure AVR server")
		if avrTransport == "file" {
			go createRotatingLogger(logFileName, loggingPath)
			go createExportFile(logFileName, loggingPath, loggingExportFolder)
		} else {
			logger.Info().Msg("Not logging to file, skipping generation of exportfiles or logrotation")
		}
		startAvrSecure()
	} else {
		// fmt.Println("incorrect arguments provided for server launch")
		logger.Fatal().Msg("incorrect arguments provided for server launch")
		return
	}
}
