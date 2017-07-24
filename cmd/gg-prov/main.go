package main

import (
	"log"

	"github.com/alecthomas/kingpin"
	"github.com/wolfeidau/ggprov"
)

var (
	ggName = kingpin.Arg("name", "Name of greengrass core.").Required().String()
)

func main() {

	kingpin.Parse()

	svcs, err := ggprov.CreateSvcs()
	if err != nil {
		log.Fatalf("Failed create aws svcs: %+v", err)
	}

	// aws greengrass get-service-role-for-account
	// if this returns 404 (not found) then create the policy
	role, err := svcs.CreateOrGetServiceRoleForAccount()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	log.Println("Role", role)

	// aws iot create-thing --thing-name "$GGC_DEPLOYMENT"
	thing, err := svcs.CreateThing(*ggName)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// aws iot create-keys-and-certificate --set-as-active > tmp-ggc-cert.json
	thingCreds, err := svcs.CreateKeysAndCertificates()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// aws iot create-policy --policy-name "$GGC_DEPLOYMENT-IOT-Policy" --policydocument
	policy, err := svcs.CreateThingPolicy(*ggName)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// aws iot attach-principal-policy --policy-name $POLICYNAME_IOT --principal $CERTARN_GGC
	err = svcs.AttachPrincipalPolicy(thingCreds, policy)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// aws iot attach-thing-principal --thing-name $THINGNAME_GGC --principal $CERTARN_GGC
	err = svcs.AttachThingPrincipal(thing, thingCreds)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	endpoint, err := svcs.GetIoTEndpoint()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	config := ggprov.NewThingConfig(role, policy, thing, thingCreds, endpoint)

	err = config.Save(thing.Name)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
