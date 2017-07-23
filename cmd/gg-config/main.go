package main

import (
	"io/ioutil"
	"log"

	"path"

	"github.com/wolfeidau/ggprov"
)

const ggSysLocalConf = `# created by gg-conf
fs.protected_hardlinks = 1
fs.protected_symlinks = 1
`

const ggConfigTmpl = `{
    "coreThing": {
        "caPath": "root-ca.pem",
        "certPath": "cloud.pem.crt",
        "keyPath": "private.pem.key",
        "thingArn": "{{.Thing.Arn}}",
        "iotHost": "{{.IotEndpoint.Hostname}}",
        "ggHost": "greengrass.iot.ap-southeast-2.amazonaws.com"
    },
    "runtime": {
        "cgroup": {
            "useSystemd": "yes"
        }
    }
}
`

const (
	ggLoginUser = "linaro"
	ggTarFile   = "greengrass-linux-aarch64-1.0.0.tar.gz"

	ggCACertURL = "https://www.symantec.com/content/en/us/enterprise/verisign/roots/VeriSign-Class%203-Public-Primary-Certification-Authority-G5.pem"
)

func main() {

	log.Println("Creating /etc/sysctl.d/local-ggc.conf")

	err := ioutil.WriteFile("/etc/sysctl.d/local-ggc.conf", []byte(ggSysLocalConf), 0660)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// does the greengrass exist
	exists, err := ggprov.UserExist("greengrass")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if !exists {
		// useradd --system --no-create-home --user-group greengrass
		err = ggprov.NewSystemUserAndGroup("greengrass")
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

	tarPath := path.Join("/home", ggLoginUser, ggTarFile)

	err = ggprov.RunCommand("tar", []string{"xvf", tarPath, "-C", "/"})
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ggprov.DownloadFromURL(ggCACertURL, "/greengrass/configuration/certs/root-ca.pem")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

}
