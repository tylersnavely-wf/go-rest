/*
Copyright 2014 - 2015 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rest

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/newrelic/go-agent"
)

var appSetupLock = &sync.Mutex{}
var newRelicApp newrelic.Application

// setUpAPM loads the New Relic configuration and loads the app needed for reporting.
func setUpAPM() (newrelic.Application, error) {
	appSetupLock.Lock()
	defer appSetupLock.Unlock()
	if newRelicApp != nil {
		return newRelicApp, nil
	}

	relicAppKey := os.Getenv("NEW_RELIC_APP_NAME")
	relicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	cfg := newrelic.NewConfig(relicAppKey, relicLicenseKey)

	cfg.DistributedTracer.Enabled = true
	cfg.Logger = newrelic.NewLogger(os.Stdout)

	relicLabels := os.Getenv("NEW_RELIC_LABELS")
	labelList := strings.Split(relicLabels, ";")
	for _, pair := range labelList {
		l := strings.Split(pair, ":")
		if len(l) != 2 {
			log.Println("Bad formatting on the NEW_RELIC_LABELS env var")
		} else {
			cfg.Labels[l[0]] = l[1]
		}
	}

	app, err := newrelic.NewApplication(cfg)
	if err != nil {
		log.Println("Error starting New Relic application")
		return nil, err
	}

	newRelicApp = app

	return newRelicApp, nil
}
