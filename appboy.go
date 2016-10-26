package main

type AppboyRecord struct {
	Devices []AppboyDevices   `json:"devices"`
	Tokens  []AppboyPushToken `json:"push_tokens"`
	Apps    []AppboyApps      `json:"apps"`
}

type AppboyPushToken struct {
	App      string `json:"app"`
	Platform string `json:"platform"`
	Token    string `json:"token"`
}

type AppboyDevices struct {
	Model   string `json:"model"`
	OS      string `json:"os"`
	Carrier string `json:"carrier"`
	IDFV    string `json:"device_id"`
}

type AppboyApps struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	Version   string `json:"versions"`
	Sessions  int    `json:"sessions"`
	FirstUsed string `json:"first_used"`
	LastUsed  string `json:"last_used"`
}
