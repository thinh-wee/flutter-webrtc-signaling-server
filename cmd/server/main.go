package main

import (
	"os"

	"github.com/flutter-webrtc/flutter-webrtc-server/pkg/logger"
	"github.com/flutter-webrtc/flutter-webrtc-server/pkg/signaler"
	"github.com/flutter-webrtc/flutter-webrtc-server/pkg/turn"
	"github.com/flutter-webrtc/flutter-webrtc-server/pkg/websocket"
	"gopkg.in/ini.v1"
)

func main() {

	cfg, err := ini.Load("configs/config.ini")
	if err != nil {
		logger.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}

	publicIP := cfg.Section("turn").Key("public_ip").String()
	stunPort, err := cfg.Section("turn").Key("port").Int()
	if err != nil {
		stunPort = 3478
	}
	realm := cfg.Section("turn").Key("realm").String()

	turnConfig := turn.TurnServerConfig{
		PublicIP: publicIP,
		Port:     stunPort,
		Realm:    realm,
		Username: cfg.Section("turn").Key("username").String(),
		Password: cfg.Section("turn").Key("password").String(),
	}

	signaler := signaler.NewSignaler(&turnConfig)
	wsServer := websocket.NewWebSocketServer(signaler.HandleNewWebSocket, signaler.HandleTurnServerCredentials)

	// sslCert := cfg.Section("general").Key("cert").String()
	// sslKey := cfg.Section("general").Key("key").String()
	bindAddress := cfg.Section("general").Key("bind").String()

	port, err := cfg.Section("general").Key("port").Int()
	if err != nil {
		port = 8086
	}

	htmlRoot := cfg.Section("general").Key("html_root").String()

	config := websocket.DefaultConfig()
	config.Host = bindAddress
	config.Port = port
	config.HTMLRoot = htmlRoot
	config.TurnServerPath = "/turn"
	config.WebSocketPath = "/ws/connect"

	wsServer.Bind(config)
}
