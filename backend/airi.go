package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kiosk404/airi-go/backend/api/middleware"
	"github.com/kiosk404/airi-go/backend/api/router"
	"github.com/kiosk404/airi-go/backend/application"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ternary"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/version/verflag"
	"github.com/spf13/cobra"
)

const (
	// defaultCrashLogFile
	defaultCrashLogFile = "crash.log"

	// recommendedLogDir 定义日志输出的地址
	recommendedLogDir = "./output/"

	// appName defines the executable binary filename for route scheduler component
	appName = "airi-go"
)

func NewAiriGoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appName,
		Short: "Airi-Go",
		Long: `With the power of modern large language models like ChatGPT and famous Claude, 
		asking a virtual being to roleplay and chat with us is already easy enough for everyone. 
		Platforms like Character.ai (a.k.a. c.ai) and JanitorAI as well as local playgrounds like 
		SillyTavern are already good-enough solutions for a chat based or visual adventure game like experience.`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()
			return run()
		},
		PostRun: func(cmd *cobra.Command, args []string) {

		},
	}

	cobra.OnInitialize(setCrashOutput, loadEnv, initLog)

	return cmd
}

func run() error {
	if err := application.Init(context.Background()); err != nil {
		panic("InitializeInfra failed, err=" + err.Error())
	}

	return startHTTPServer()
}

func startHTTPServer() error {
	maxRequestBodySize := os.Getenv("MAX_REQUEST_BODY_SIZE")
	maxsize := conv.StrToInt64D(maxRequestBodySize, 1024*1024*200)
	addr := getEnv("HTTP_ADDR", ":8888")

	s := gin.Default()

	s.Use(middleware.MaxBodySizeMiddleware(maxsize))

	ginMode := getEnv("GIN_MODE", gin.DebugMode)
	gin.SetMode(ginMode)

	// cors option
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	s.Use(cors.New(config))

	//Middleware order matters
	s.Use(middleware.ContextCacheMW())     // must be first
	s.Use(middleware.RequestInspectorMW()) // must be second
	s.Use(middleware.SetHostMW())
	s.Use(middleware.SetLogIDMW())
	s.Use(middleware.AccessLogMW())
	s.Use(middleware.OpenapiAuthMW())
	s.Use(middleware.SessionAuthMW())

	router.GeneratedRegister(s)

	logs.Info("server start !!")
	return s.Run(addr)
}

func loadEnv() {
	appEnv := os.Getenv("APP_ENV")
	fileName := ternary.IFElse(appEnv == "", ".env", ".env"+appEnv)

	logs.Info("load env from %s", fileName)

	err := godotenv.Load(fileName)
	if err != nil {
		panic("load env file failed, err=" + err.Error())
	}
}

func initLog() {
	logBasePath := recommendedLogDir
	logPath := fmt.Sprintf("%s%s", logBasePath, "log/common.log")
	// 初始化日志打印
	if err := logs.InitLog(logPath); err != nil {
		panic(err)
	}
}

func setCrashOutput() {
	crashFile, _ := os.Create(defaultCrashLogFile)

	if err := debug.SetCrashOutput(crashFile, debug.CrashOptions{}); err != nil {
		return
	}
}

func getEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
