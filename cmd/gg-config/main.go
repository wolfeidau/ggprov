package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"path"

	"github.com/alecthomas/kingpin"
	"github.com/davecgh/go-spew/spew"
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
	ggLoginUser      = "linaro"
	ggTarFile        = "greengrass-linux-aarch64-1.0.0.tar.gz"
	ggConfigFilePath = "/greengrass/configuration/config.json"
	ggCertFilePath   = "/greengrass/configuration/certs/"

	ggCACertURL = "https://www.symantec.com/content/en/us/enterprise/verisign/roots/VeriSign-Class%203-Public-Primary-Certification-Authority-G5.pem"
)

var (
	ggName = kingpin.Arg("name", "Name of greengrass core.").Required().String()

	configTmpl = template.Must(template.New("config.json").Parse(ggConfigTmpl))
)

func main() {
	kingpin.Parse()

	configPath := path.Join("/home", ggLoginUser, fmt.Sprintf("%s.yaml", *ggName))

	config, err := ggprov.Load(configPath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	spew.Dump(config)

	log.Println("Creating /etc/sysctl.d/local-ggc.conf")

	err = ioutil.WriteFile("/etc/sysctl.d/local-ggc.conf", []byte(ggSysLocalConf), 0660)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// does the greengrass exist
	exists, err := ggprov.UserExist("ggc_user")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if !exists {
		err = ggprov.NewSystemUser("ggc_user")
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

	// does the greengrass exist
	exists, err = ggprov.GroupExist("ggc_group")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if !exists {
		err = ggprov.NewSystemGroup("ggc_group")
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

	tarPath := path.Join("/home", ggLoginUser, ggTarFile)

	err = ggprov.RunCommand("tar", []string{"xf", tarPath, "-C", "/"})
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ggprov.DownloadFromURL(ggCACertURL, "/greengrass/configuration/certs/root-ca.pem")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	ggConfigFile, err := os.Create(ggConfigFilePath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = configTmpl.Execute(ggConfigFile, config)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ioutil.WriteFile(path.Join(ggCertFilePath, "cloud.pem.crt"), []byte(config.ThingCreds.CertificatePem), 0660)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ioutil.WriteFile(path.Join(ggCertFilePath, "private.pem.key"), []byte(config.ThingCreds.PrivateKey), 0660)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ggprov.AptUpdate()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = ggprov.AptInstall([]string{"sqlite3"})
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
