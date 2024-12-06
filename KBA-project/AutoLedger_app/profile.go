package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"tnrto": {
		CertPath:     "../AutoLedger/organizations/peerOrganizations/tnrto.vaahan.com/users/User1@tnrto.vaahan.com/msp/signcerts/cert.pem",
		KeyDirectory: "../AutoLedger/organizations/peerOrganizations/tnrto.vaahan.com/users/User1@tnrto.vaahan.com/msp/keystore/",
		TLSCertPath:  "../AutoLedger/organizations/peerOrganizations/tnrto.vaahan.com/peers/peer0.tnrto.vaahan.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.tnrto.vaahan.com",
		MSPID:        "tnrtoMSP",
	},

	"klrto": {
		CertPath:     "../AutoLedger/organizations/peerOrganizations/klrto.vaahan.com/users/User1@klrto.vaahan.com/msp/signcerts/cert.pem",
		KeyDirectory: "../AutoLedger/organizations/peerOrganizations/klrto.vaahan.com/users/User1@klrto.vaahan.com/msp/keystore/",
		TLSCertPath:  "../AutoLedger/organizations/peerOrganizations/klrto.vaahan.com/peers/peer0.klrto.vaahan.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.klrto.vaahan.com",
		MSPID:        "klrtoMSP",
	},

	"knrto": {
		CertPath:     "../AutoLedger/organizations/peerOrganizations/knrto.vaahan.com/users/User1@knrto.vaahan.com/msp/signcerts/cert.pem",
		KeyDirectory: "../AutoLedger/organizations/peerOrganizations/knrto.vaahan.com/users/User1@knrto.vaahan.com/msp/keystore/",
		TLSCertPath:  "../AutoLedger/organizations/peerOrganizations/knrto.vaahan.com/peers/peer0.knrto.vaahan.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.knrto.vaahan.com",
		MSPID:        "knrtoMSP",
	},
}
