package main

import (
	"os"
	"strconv"
	"time"

	"github.com/DENICeG/go-rriclient/pkg/rri"
	whiteflag "github.com/danielb42/whiteflag" // MIT
	"github.com/sirupsen/logrus"
)

var (
	regacc          string
	password        string
	rriServer       string
	domainToCheck   string
	rriport                       = 51131
	rriResponseTime time.Duration = 0
	log                           = logrus.New()
)

func setAliasesViaWhiteflag() {
	whiteflag.Alias("dom", "domain", "use the given domain for check order")
	whiteflag.Alias("reg", "regacc", "sets the regacc used to perform the check")
	whiteflag.Alias("pw", "password", "sets the password used to perform the check")
	whiteflag.Alias("host", "rriserver", "sets the rri-server used to perform the check")
	whiteflag.Alias("port", "rriport", "sets the rri-server used to perform the check")
}

func main() {
	// Parse commandline parameters
	setAliasesViaWhiteflag()
	whiteflag.ParseCommandLine()

	regacc = whiteflag.GetString("regacc")
	password = whiteflag.GetString("password")
	domainToCheck = whiteflag.GetString("domain")

	// Use custom rri port if given
	if whiteflag.CheckInt("rriport") {
		rriport = whiteflag.GetInt("rriport")
	}

	rriServer = whiteflag.GetString("rriserver") + ":" + strconv.Itoa(rriport)

	// create client and perform command
	client, err := rri.NewClient(rriServer)
	if err != nil {
		log.Errorln(err)
		os.Exit(2)
	}

	err = client.Login(regacc, password)
	if err != nil {
		log.Errorln(err)
		os.Exit(2)
	}

	timeBegin := time.Now()
	checkQuery := rri.NewCheckDomainQuery(domainToCheck)
	rriResponse, err := client.SendQuery(checkQuery)

	if err != nil {
		rriResponseTime = 0
		log.Errorln("technical error: ", err)
		os.Exit(2)
	} else {
		if rriResponse.IsSuccessful() {
			rriResponseTime = time.Since(timeBegin)

			log.WithFields(logrus.Fields{
				"ResultField":       rriResponse.FirstField("Status"),
				"Result":            rriResponse.Result(),
				"DurationFromAgent": rriResponseTime,
			}).Info("SENSU-RRI-Check successful")
			_ = client.Logout()

			os.Exit(0)
		} else {
			rriResponseTime = 0
			log.Errorf("Result: %s, ErrorMsg: %s", rriResponse.Result(), rriResponse.ErrorMsg())
			_ = client.Logout()
		}
	}

	os.Exit(2)
}
