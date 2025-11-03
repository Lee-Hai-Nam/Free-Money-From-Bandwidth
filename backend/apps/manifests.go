package apps

// AppManifest represents a complete app configuration
type AppManifest struct {
	Name               string
	Dashboard          string
	Link               string
	Image              string
	Environment        map[string]string // Original env var mapping
	RequiredFields     map[string]bool
	Volumes            []string
	Ports              []string
	Command            string
	NetworkMode        string
	ResourceLimits     *ResourceLimits
	AutoGenerateFields map[string]*AutoGenerateConfig
	Disabled           bool // Flag to mark if the app is currently not working
}

// ResourceLimits represents resource constraints
type ResourceLimits struct {
	CPUs              string
	MemoryReservation string
	MemoryLimit       string
}

// AutoGenerateConfig represents auto-generation settings for fields
type AutoGenerateConfig struct {
	Length  int
	Prefix  string
	Charset string
}

// GetAppManifest returns the manifest for an app
func GetAppManifest(appID string) *AppManifest {
	manifests := GetAllManifests()
	return manifests[appID]
}

// GetAllManifests returns all app manifests
func GetAllManifests() map[string]*AppManifest {
	return map[string]*AppManifest{
		"earnapp": {
			Name:      "EarnApp",
			Dashboard: "https://earnapp.com/dashboard",
			Link:      "https://earnapp.com/i/3zulx7k",
			Image:     "fazalfarhan01/earnapp:lite",
			Environment: map[string]string{
				"EARNAPP_UUID": "$EARNAPP_UUID",
				"EARNAPP_TERM": "yes",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME": true,
			},
			Volumes: []string{".data/.earnapp:/etc/earnapp"},
			ResourceLimits: &ResourceLimits{
				CPUs:              "1.0",
				MemoryReservation: "128m",
				MemoryLimit:       "512m",
			},
			AutoGenerateFields: map[string]*AutoGenerateConfig{
				"EARNAPP_UUID": {
					Length:  32,
					Prefix:  "sdk-node-",
					Charset: "abcdefghijklmnopqrstuvwxyz0123456789",
				},
			},
		},
		"honeygain": {
			Name:      "Honeygain",
			Dashboard: "https://dashboard.honeygain.com/",
			Link:      "https://r.honeygain.me/MINDL15721",
			Image:     "honeygain/honeygain:latest",
			Environment: map[string]string{
				"HONEYGAIN_DUMMY": "",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":        true,
				"HONEYGAIN_EMAIL":    true,
				"HONEYGAIN_PASSWORD": true,
			},
			Command: "-tou-accept -email $HONEYGAIN_EMAIL -pass $HONEYGAIN_PASSWORD -device $DEVICE_NAME",
			ResourceLimits: &ResourceLimits{
				CPUs:              "1.0",
				MemoryReservation: "128m",
				MemoryLimit:       "512m",
			},
		},
		"iproyalpawns": {
			Name:      "IPRoyal Pawns",
			Dashboard: "https://dashboard.pawns.app/",
			Link:      "https://pawns.app?r=MiNe",
			Image:     "iproyal/pawns-cli:latest",
			Environment: map[string]string{
				"IPROYALPAWNS_DUMMY": "",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":           true,
				"IPROYALPAWNS_EMAIL":    true,
				"IPROYALPAWNS_PASSWORD": true,
			},
			Command: "-accept-tos -email=$IPROYALPAWNS_EMAIL -password=$IPROYALPAWNS_PASSWORD -device-name=$DEVICE_NAME -device-id=id_$DEVICE_NAME",
			ResourceLimits: &ResourceLimits{
				CPUs:              "0.5",
				MemoryReservation: "64m",
				MemoryLimit:       "256m",
			},
		},
		"packetstream": {
			Name:      "PacketStream",
			Dashboard: "https://packetstream.io/dashboard",
			Link:      "https://packetstream.io/?psr=3zSD",
			Image:     "packetstream/psclient:latest",
			Environment: map[string]string{
				"CID": "$PACKETSTREAM_CID",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"PACKETSTREAM_CID": true,
			},
			ResourceLimits: &ResourceLimits{
				CPUs:              "1.0",
				MemoryReservation: "128m",
				MemoryLimit:       "512m",
			},
		},
		"traffmonetizer": {
			Name:      "TraffMonetizer",
			Dashboard: "https://app.traffmonetizer.com/dashboard",
			Link:      "https://traffmonetizer.com/?aff=366499",
			Image:     "traffmonetizer/cli_v2:latest",
			Environment: map[string]string{
				"TRAFFMONETIZER_DUMMY": "",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":          true,
				"TRAFFMONETIZER_TOKEN": true,
			},
			Command: "start accept status --token $TRAFFMONETIZER_TOKEN --device-name $DEVICE_NAME",
			ResourceLimits: &ResourceLimits{
				CPUs:              "0.5",
				MemoryReservation: "64m",
				MemoryLimit:       "256m",
			},
		},
		// Additional apps
		"repocket": {
			Name:      "Repocket",
			Dashboard: "https://app.repocket.co/#home",
			Link:      "https://link.repocket.co/hr8i",
			Image:     "repocket/repocket:latest",
			Environment: map[string]string{
				"RP_EMAIL":   "$REPOCKET_EMAIL",
				"RP_API_KEY": "$REPOCKET_APIKEY",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":     true,
				"REPOCKET_EMAIL":  true,
				"REPOCKET_APIKEY": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"earnfm": {
			Name:      "EarnFM",
			Dashboard: "https://app.earn.fm/",
			Link:      "https://earn.fm/ref/MATTTAV6",
			Image:     "earnfm/earnfm-client:latest",
			Environment: map[string]string{
				"EARNFM_TOKEN": "$EARNFM_APIKEY",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":   true,
				"EARNFM_APIKEY": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"proxyrack": {
			Name:      "ProxyRack",
			Dashboard: "https://peer.proxyrack.com/dashboard",
			Link:      "https://peer.proxyrack.com/ref/myoas6qttvhuvkzh8ffx90ns1ouhwgilfgamo5ex",
			Image:     "proxyrack/pop:latest",
			Environment: map[string]string{
				"API_KEY":     "$PROXYRACK_APIKEY",
				"DEVICE_NAME": "$DEVICE_NAME",
				"UUID":        "$PROXYRACK_UUID",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"PROXYRACK_APIKEY": true,
				"PROXYRACK_UUID":   true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
		},
		"proxylite": {
			Name:      "ProxyLite",
			Dashboard: "https://proxylite.ru/",
			Link:      "https://proxylite.ru/?r=PJTKXWN3",
			Image:     "proxylite/proxyservice:latest",
			Environment: map[string]string{
				"USER_ID": "$PROXYLITE_USERID",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"PROXYLITE_USERID": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
		},
		"bitping": {
			Name:      "Bitping",
			Dashboard: "https://app.bitping.com/earnings",
			Link:      "https://app.bitping.com?r=qm7mIuX3",
			Image:     "bitping/bitpingd:latest",
			Environment: map[string]string{
				"BITPING_EMAIL":    "$BITPING_EMAIL",
				"BITPING_PASSWORD": "$BITPING_PASSWORD",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"BITPING_EMAIL":    true,
				"BITPING_PASSWORD": true,
			},
			Volumes:        []string{".data/.bitpingd:/root/.bitpingd"},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
		},
		"packetshare": {
			Name:      "PacketShare",
			Dashboard: "https://packetshare.io/ucenter.html",
			Link:      "https://www.packetshare.io/?code=A260871CFD822E35",
			Image:     "packetshare/packetshare:latest",
			Environment: map[string]string{
				"PACKETSHARE_EMAIL":    "$PACKETSHARE_EMAIL",
				"PACKETSHARE_PASSWORD": "$PACKETSHARE_PASSWORD",
				"PACKETSHARE_DUMMY":    "",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":          true,
				"PACKETSHARE_EMAIL":    true,
				"PACKETSHARE_PASSWORD": true,
			},
			Command:        "-accept-tos -email=$PACKETSHARE_EMAIL -password=$PACKETSHARE_PASSWORD",
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"grass": {
			Name:      "Grass",
			Dashboard: "https://app.getgrass.io/dashboard",
			Link:      "https://app.getgrass.io/register/?referralCode=qyvJmxgNUhcLo2f",
			Image:     "mrcolorrain/grass-node:latest",
			Environment: map[string]string{
				"USER_EMAIL":    "$GRASS_EMAIL",
				"USER_PASSWORD": "$GRASS_PASSWORD",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":    true,
				"GRASS_EMAIL":    true,
				"GRASS_PASSWORD": true,
			},
			Volumes:        []string{".data/.grass:/app/chrome_user_data"},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"gradient": {
			Name:      "Gradient",
			Dashboard: "https://app.gradient.network/dashboard",
			Link:      "https://app.gradient.network/signup?code=9WOBKP",
			Image:     "carbon2029/dockweb:latest",
			Environment: map[string]string{
				"GRADIENT_EMAIL": "$GRADIENT_EMAIL",
				"GRADIENT_PASS":  "$GRADIENT_PASSWORD",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":       true,
				"GRADIENT_EMAIL":    true,
				"GRADIENT_PASSWORD": true,
			},
			Volumes:        []string{".data/.gradient:/app/chrome_user_data"},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"dawn": {
			Name:      "Dawn",
			Dashboard: "https://dawninternet.com",
			Link:      "https://dawninternet.com?code=xo23vynw",
			Image:     "carbon2029/dockweb:latest",
			Environment: map[string]string{
				"DAWN_EMAIL": "$DAWN_EMAIL",
				"DAWN_PASS":  "$DAWN_PASSWORD",
			},
			Ports: []string{"${DAWN_PORT}:5000"},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":   true,
				"DAWN_EMAIL":    true,
				"DAWN_PASSWORD": true,
			},
			Volumes:        []string{".data/.dawn:/app/chrome_user_data"},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"teneo": {
			Name:      "Teneo",
			Dashboard: "https://dashboard.teneo.pro/",
			Link:      "https://dashboard.teneo.pro/?code=qPgLn",
			Image:     "carbon2029/dockweb:latest",
			Environment: map[string]string{
				"TENEO_EMAIL": "$TENEO_EMAIL",
				"TENEO_PASS":  "$TENEO_PASSWORD",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":    true,
				"TENEO_EMAIL":    true,
				"TENEO_PASSWORD": true,
			},
			Volumes:        []string{".data/.teneo:/app/chrome_user_data"},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"proxybase": {
			Name:      "Proxybase",
			Dashboard: "https://dash.proxybase.org/",
			Link:      "http://dash.proxybase.org/signup?ref=XfOz3zeURm",
			Image:     "proxybase/proxybase:latest",
			Environment: map[string]string{
				"USER_ID":     "$PROXYBASE_USERID",
				"DEVICE_NAME": "$DEVICE_NAME",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"PROXYBASE_USERID": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "2.0", MemoryReservation: "256m", MemoryLimit: "1g"},
		},
		"wipter": {
			Name:      "Wipter",
			Dashboard: "https://wipter.com/dashboard",
			Link:      "https://wipter.com/signup?ref=money4band",
			Image:     "ghcr.io/adfly8470/wipter/wipter@sha256:9b1a7742bfbbd68e86eea1719f606c7d10c884e2578a4fb35f109eed387619cd",
			Environment: map[string]string{
				"WIPTER_EMAIL":    "$WIPTER_EMAIL",
				"WIPTER_PASSWORD": "$WIPTER_PASSWORD",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":     true,
				"WIPTER_EMAIL":    true,
				"WIPTER_PASSWORD": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
			Ports:          []string{"${WIPTER_PORT_1}:5900", "${WIPTER_PORT_2}:6080"},
		},
		"ebesucher_chrome": {
			Name:      "Ebesucher Chrome",
			Dashboard: "",
			Link:      "",
			Image:     "lscr.io/linuxserver/chromium:latest",
			Environment: map[string]string{
				"CHROME_CLI":   "https://www.ebesucher.com/surfbar/$EBESUCHER_USERNAME",
				"CUSTOM_USER":  "fmfb",
				"PASSWORD":     "fmfb",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":        true,
				"EBESUCHER_USERNAME": true,
			},
			Volumes:        []string{".data/.ebesucher_chrome:/config"},
			Ports:          []string{"3000:3000"},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"ebesucher_firefox": {
			Name:      "Ebesucher Firefox",
			Dashboard: "",
			Link:      "",
			Image:     "jlesage/firefox",
			Environment: map[string]string{
				"FF_OPEN_URL":      "https://www.ebesucher.com/surfbar/$EBESUCHER_USERNAME",
				"VNC_LISTENING_PORT": "-1",
				"VNC_PASSWORD":     "fmfb",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":        true,
				"EBESUCHER_USERNAME": true,
			},
			Volumes:        []string{".data/.ebesucher_firefox:/config:rw"},
			Ports:          []string{"5800:5800"},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "256m", MemoryLimit: "1g"},
			Disabled:       true,
		},
		"adnade": {
			Name:      "AdnAde",
			Dashboard: "",
			Link:      "",
			Image:     "jlesage/firefox",
			Environment: map[string]string{
				"FF_OPEN_URL":        "https://adnade.net/view.php?user=$ADNADE_USERNAME&multi=4",
				"VNC_LISTENING_PORT": "-1",
				"WEB_LISTENING_PORT": "5900",
				"VNC_PASSWORD":       "fmfb",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":     true,
				"ADNADE_USERNAME": true,
			},
			Volumes:        []string{".data/.adnade:/config:rw"},
			Ports:          []string{"5900:5900"},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "256m", MemoryLimit: "1g"},
		},
		"packetsdk": {
			Name:      "PacketSDK",
			Dashboard: "",
			Link:      "",
			Image:     "packetsdk/packetsdk",
			Command:   "-appkey=$PACKET_SDK_APP_KEY",
			RequiredFields: map[string]bool{
				"DEVICE_NAME":           true,
				"PACKET_SDK_APP_KEY": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"gaganode": {
			Name:      "Gaganode",
			Dashboard: "",
			Link:      "",
			Image:     "xterna/gaga-node",
			Environment: map[string]string{
				"TOKEN": "$GAGANODE_TOKEN",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":     true,
				"GAGANODE_TOKEN": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"castarsdk": {
			Name:      "CastarSDK",
			Dashboard: "",
			Link:      "",
			Image:     "ghcr.io/adfly8470/castarsdk/castarsdk@sha256:30d7e9830c0144165b86dbb053eaea11e36d1b9f7ee0837fd4eda71cc6b48125",
			Environment: map[string]string{
				"KEY": "$CASTAR_SDK_KEY",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":      true,
				"CASTAR_SDK_KEY": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"peer2profit": {
			Name:      "Peer2Profit",
			Dashboard: "",
			Link:      "",
			Image:     "enwaiax/peer2profit",
			Environment: map[string]string{
				"email": "$PEER2PROFIT_EMAIL",
			},
			RequiredFields: map[string]bool{
				"DEVICE_NAME":       true,
				"PEER2PROFIT_EMAIL": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"urnetwork": {
			Name:      "UrNetwork",
			Dashboard: "",
			Link:      "",
			Image:     "bringyour/community-provider:latest",
			Command:   "provide",
			RequiredFields: map[string]bool{
				"DEVICE_NAME":   true,
				"UR_AUTH_TOKEN": true,
			},
			Volumes:        []string{".data/.urnetwork:/root/.urnetwork"},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		"titan": {
			Name:      "Titan",
			Dashboard: "",
			Link:      "",
			Image:     "nezha123/titan-edge",
			RequiredFields: map[string]bool{
				"DEVICE_NAME": true,
				"TITAN_HASH":  true,
			},
			Volumes:        []string{".data/.titanedge:/root/.titanedge"},
			ResourceLimits: &ResourceLimits{CPUs: "1.0", MemoryReservation: "128m", MemoryLimit: "512m"},
		},
		// Extra app: mystnode
		"mystnode": {
			Name:      "MystNode",
			Dashboard: "https://mystnodes.com/nodes",
			Link:      "https://mystnodes.co/?referral_code=Tc7RaS7Fm12K3Xun6mlU9q9hbnjojjl9aRBW8ZA9",
			Image:     "mysteriumnetwork/myst:latest",
			Environment: map[string]string{
				"MYSTNODE_DUMMY": "",
			},
			Command: "service --agreed-terms-and-conditions",
			Volumes: []string{".data/mysterium-node:/var/lib/mysterium-node"},
			Ports:   []string{"${MYSTNODE_PORT}:4449"},
			RequiredFields: map[string]bool{
				"DEVICE_NAME": true,
			},
			ResourceLimits: &ResourceLimits{CPUs: "4.0", MemoryReservation: "512m", MemoryLimit: "2g"},
		},
	}
}
