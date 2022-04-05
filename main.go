package main

import (
	"encoding/json"
	"fmt"
	"log"

	apiclient "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/connector"
)

const caCert = `-----BEGIN CERTIFICATE-----
MIIEEzCCAnugAwIBAgIVAPrw9qKSpERS1eXS1LsoeKvMVg6mMA0GCSqGSIb3DQEB
CwUAMCExDTALBgNVBAoTBEp1anUxEDAOBgNVBAMTB2p1anUtY2EwHhcNMjIwMjE1
MTQzNzUxWhcNMzIwMjE1MTQ0MjUxWjAhMQ0wCwYDVQQKEwRKdWp1MRAwDgYDVQQD
EwdqdWp1LWNhMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAwxa8patl
/Qs7dIaRzJQNjXRykPlwQcCYgEy02j1R0Ci3ylYn9h93CSRUdTmNZkGF91JcVs+X
Y17nVRxxt5I0zxl1KMANJbWcQ4pHoMIeLQU+p1kSSQ9FUBH1kLcHT04FtP6YNMMG
SdXAf0QxSG6bNe/+/jBMw+8WCGO/LhZoQN5JvE9bMjcbth1r4QEmbOon0c6a2MER
Nt7O1hfkyc1/W7ZxgvV9o6SHYMbXhRynj8j7dEP0niEAJl9MCO1tyYIIXr3QNirN
fWH+633C9CO3NBlulAka6eKuGCZwe1wYQw1pEJbrLL5gno46+790IKud3ttq+DmP
Oc/unXIHv1bL1oYpHLV1mgcD+tZiqQdHFNTQe4AncxubTibX22GhrPQbpVp6e04K
eitgr8DgDc65WreYi3ZnVazpZt/9VmnJuYPIahwRmu9wKaambnMLsWzY3kC4wwss
aGDp2iYQ/HGP9AEVRwkFc9inYYJV+CjVXZc0m/gYzEUz7CD4Os6SZVz7AgMBAAGj
QjBAMA4GA1UdDwEB/wQEAwICpDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQu
Do1DXrXo2G5faRpLnxCGPsJxgTANBgkqhkiG9w0BAQsFAAOCAYEAJjorwT77X2Eo
DkLvqFTQr1TEzZBlGWExINThkSKRCWr+4iLTKX0csIpZ6TZrqg2aMr5B0fltTFNs
cpxNdf1myaXYqtcPsIsTrYut/KWrml5I4NgNBqMpCoMoWG39xyJi92mGj3Ppp/hk
eDm9ovcM5EtgrxnHlRhGlalNzlom/hDt/lWpLoESKv4eqRq0uiDQkQ+9dmEjAKrk
fS1cnH4Stq3WdkA9Dvnutdx3mg4qc5O31hSMisCOaNlypLL9cnP0J42hVfetv9nR
yEQq4fryT4EeWVT6jrFMHphtaLALzvDnSOjT55nMWYKXdr8dnl8ixDF5aIZFH4ft
RYdbq0iE2nToGJK3ZD3PHqZnilss9A00cTJzENiy4pFotQ3A530Zv9O5DNv+ZEhd
76IAn8E3Qqwq0k7JcYwK4XubxbJpF4bAKC/OL8aFaksozTRRuzTFL45CZEO7rulN
Dt5NVEFyoE24ETwnu4FzcFH59GbF0ShU+VcBdk9Y2oO26dH7jASX
-----END CERTIFICATE-----
`

func main() {
	// connr, err := connector.NewClientStore(connector.ClientStoreConfig{
	// 	ControllerName: "overlord",
	// 	ModelUUID:      "1a5b37f4-346a-450b-8d74-3e78e7bb47e8",
	// })
	connr, err := connector.NewSimple(connector.SimpleConfig{
		ControllerAddresses: []string{"10.225.205.241:17070"},
		ModelUUID:           "1a5b37f4-346a-450b-8d74-3e78e7bb47e8",
		CACert:              caCert,
		Username:            "admin",
		Password:            "password1",
	})
	if err != nil {
		log.Fatalf("Error getting connector: %s", err)
	}

	// Get a connection
	conn, err := connr.Connect()
	if err != nil {
		log.Fatalf("Error opening connection: %s", err)
	}
	defer conn.Close()

	fmt.Println("Deploying postgresql charm...")
	// Try deploying a charm
	err = DeployCharm(conn, DeployCharmArgs{
		CharmName: "postgresql",
		NumUnits:  1,
		Revision:  -1,
	})
	if err != nil {
		fmt.Printf("Error deploying postgresql charm: %s\n", err)
	} else {
		fmt.Println("Success deploying postgresql charm")
	}

	// Get a Client facade client
	client := apiclient.NewClient(conn)

	// Call the Status endpoint of the client facade
	status, err := client.Status(nil)
	if err != nil {
		log.Fatalf("Error requesting status: %s", err)
	}

	// Print to stdout.
	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling response: %s", err)
	}
	fmt.Printf("%s\n", b)
}
