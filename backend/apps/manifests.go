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
			Name:      "EARNAPP",
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
			Name:      "HONEYGAIN",
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
			Name:      "IPROYALPAWNS",
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
			Name:      "PACKETSTREAM",
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
			Name:      "TRAFFMONETIZER",
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
			Name:      "REPOCKET",
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
			Name:      "EARNFM",
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
			Name:      "PROXYRACK",
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
			Name:      "PROXYLITE",
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
			Name:      "BITPING",
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
			Name:      "PACKETSHARE",
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
			Name:      "GRASS",
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
		},
		"gradient": {
			Name:      "GRADIENT",
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
		},
		"dawn": {
			Name:      "DAWN",
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
		},
		"teneo": {
			Name:      "TENEO",
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
		},
		"proxybase": {
			Name:      "PROXYBASE",
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
			Name:      "WIPTER",
			Dashboard: "https://wipter.com/dashboard",
			Link:      "https://wipter.com/signup?ref=money4band",
			Image:     "ghcr.io/techroy23/docker-wipter:latest",
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
		// Extra app: mystnode
		"mystnode": {
			Name:      "MYSTNODE",
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
