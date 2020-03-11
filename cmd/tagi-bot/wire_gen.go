// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/fsufitch/tagialisi-bot/bot"
	"github.com/fsufitch/tagialisi-bot/bot/groups-module"
	log2 "github.com/fsufitch/tagialisi-bot/bot/log-module"
	"github.com/fsufitch/tagialisi-bot/bot/memelink-module"
	"github.com/fsufitch/tagialisi-bot/bot/ping-module"
	"github.com/fsufitch/tagialisi-bot/bot/sockpuppet-module"
	"github.com/fsufitch/tagialisi-bot/config"
	"github.com/fsufitch/tagialisi-bot/db/acl-dao"
	"github.com/fsufitch/tagialisi-bot/db/connection"
	"github.com/fsufitch/tagialisi-bot/db/memes-dao"
	"github.com/fsufitch/tagialisi-bot/log"
	"github.com/fsufitch/tagialisi-bot/web"
	"github.com/fsufitch/tagialisi-bot/web/auth"
)

// Injectors from wire.go:

func InitializeMain() (Main, func(), error) {
	interruptContext := ProvideInterruptContext()
	logger := log.ProvideLogger()
	module := &ping.Module{
		Log: logger,
	}
	debugMode, err := config.ProvideDebugModeFromEnvironment()
	if err != nil {
		return Main{}, nil, err
	}
	discordLogChannel := config.ProvideDiscordLogChannelFromEnvironment()
	logModule := &log2.Module{
		Log:        logger,
		DebugMode:  debugMode,
		LogChannel: discordLogChannel,
	}
	sockpuppetModule := &sockpuppet.Module{
		Log: logger,
	}
	databaseString, err := config.ProvideDatabaseStringFromEnvironment()
	if err != nil {
		return Main{}, nil, err
	}
	databaseConnection, cleanup, err := connection.ProvidePostgresDatabaseConnection(logger, databaseString)
	if err != nil {
		return Main{}, nil, err
	}
	dao := &memes.DAO{
		Conn: databaseConnection,
	}
	aclDAO := &acl.DAO{
		Conn: databaseConnection,
	}
	memelinkModule := &memelink.Module{
		Log:     logger,
		MemeDAO: dao,
		ACLDAO:  aclDAO,
	}
	managedGroupPrefix := config.ProvideManagedGroupPrefixFromEnvironment()
	groupsModule := &groups.Module{
		Log:    logger,
		Prefix: managedGroupPrefix,
	}
	modules := bot.Modules{
		Ping:       module,
		Log:        logModule,
		SockPuppet: sockpuppetModule,
		MemeLink:   memelinkModule,
		Groups:     groupsModule,
	}
	moduleList := bot.ProvideModuleList(modules)
	botModuleBlacklist := config.ProvideBotModuleBlacklistFromEnvironment()
	discordBotToken, err := config.ProvideDiscordBotTokenFromEnvironment()
	if err != nil {
		cleanup()
		return Main{}, nil, err
	}
	tagioalisiBot := &bot.TagioalisiBot{
		Log:             logger,
		Modules:         moduleList,
		ModuleBlacklist: botModuleBlacklist,
		Token:           discordBotToken,
	}
	stdOutReceiver := log.ProvideStdOutReceiver(debugMode)
	stdErrReceiver := log.ProvideStdErrReceiver(debugMode)
	cliLoggingBootstrapper := log.CLILoggingBootstrapper{
		Logger:         logger,
		StdOutReceiver: stdOutReceiver,
		StdErrReceiver: stdErrReceiver,
	}
	webEnabled, err := config.ProvideWebEnabledFromEnvironment()
	if err != nil {
		cleanup()
		return Main{}, nil, err
	}
	webPort, err := config.ProvideWebPortFromEnvironment()
	if err != nil {
		cleanup()
		return Main{}, nil, err
	}
	webSecret := config.ProvideWebSecretFromEnvironment()
	secretBearerAuthorizationWrapper := &web.SecretBearerAuthorizationWrapper{
		Secret: webSecret,
		Log:    logger,
	}
	helloHandler := &web.HelloHandler{
		Log: logger,
	}
	sockpuppetHandler := &web.SockpuppetHandler{
		BotModule: sockpuppetModule,
		Log:       logger,
	}
	oAuth2Config := config.ProvideOAuth2ConfigFromEnvironment()
	loginStates := _wireLoginStatesValue
	loginHandler := &web.LoginHandler{
		OAuth2Config: oAuth2Config,
		LoginStates:  loginStates,
	}
	memorySessionStorage := auth.ProvideMemorySessionStorage()
	jwthmacSecret := config.ProvideJWTHMACSecretFromEnvironment()
	jwtSupport := auth.JWTSupport{
		JWTSecret: jwthmacSecret,
	}
	cookieSupport := auth.CookieSupport{
		JWT: jwtSupport,
	}
	authCodeHandler := &web.AuthCodeHandler{
		OAuth2Config:   oAuth2Config,
		LoginStates:    loginStates,
		SessionStorage: memorySessionStorage,
		AuthCookie:     cookieSupport,
	}
	logoutHandler := &web.LogoutHandler{
		SessionStorage: memorySessionStorage,
		AuthCookie:     cookieSupport,
	}
	whoAmIHandler := &web.WhoAmIHandler{
		Log:            logger,
		SessionStorage: memorySessionStorage,
		AuthCookie:     cookieSupport,
	}
	router := web.ProvideRouter(secretBearerAuthorizationWrapper, helloHandler, sockpuppetHandler, loginHandler, authCodeHandler, logoutHandler, whoAmIHandler)
	tagioalisiAPIServer := web.TagioalisiAPIServer{
		WebPort: webPort,
		Log:     logger,
		Router:  router,
	}
	webRunFunc := ProvideWebRunFunc(webEnabled, tagioalisiAPIServer)
	mainMain, cleanup2, err := ProvideMain(interruptContext, tagioalisiBot, logger, debugMode, cliLoggingBootstrapper, webRunFunc)
	if err != nil {
		cleanup()
		return Main{}, nil, err
	}
	return mainMain, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireLoginStatesValue = auth.LoginStates{}
)
