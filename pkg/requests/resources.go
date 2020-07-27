package requests

import (
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kitabisa/teler/common"
	"github.com/kitabisa/teler/configs"
	"github.com/kitabisa/teler/pkg/parsers"
	log "github.com/projectdiscovery/gologger"
)

var resource *configs.Resources
var hasExclude bool

// Get all resources and saved all to memory
func Get(options *common.Options) {
	getThreat(options)
	getFilter(options)
}

func getThreat(options *common.Options) {
	client := Client()
	excludes := options.Configs.Rules.Threat.Excludes
	resource, _ = parsers.GetResources()

	for i := 0; i < len(resource.Threat); i++ {
		hasExclude = false
		threat := reflect.ValueOf(&resource.Threat[i]).Elem()

		for j := 0; j < len(excludes); j++ {
			if excludes[j] == threat.FieldByName("Category").String() {
				hasExclude = true
			}
		}

		if hasExclude {
			continue
		}

		log.Infof("Getting \"%s\" resource...\n", threat.FieldByName("Category").String())

		req, _ := http.NewRequest("GET", threat.FieldByName("URL").String(), nil)
		resp, _ := client.Do(req)

		body, _ := ioutil.ReadAll(resp.Body)
		threat.FieldByName("Content").SetString(string(body))
	}
}

func getFilter(options *common.Options) {
	client := Client()
	excludes := options.Configs.Rules.Filter.Excludes
	resource, _ = parsers.GetResources()

	for i := 0; i < len(resource.Filter); i++ {
		hasExclude = false
		threat := reflect.ValueOf(&resource.Filter[i]).Elem()

		for j := 0; j < len(excludes); j++ {
			if excludes[j] == threat.FieldByName("Category").String() {
				hasExclude = true
			}
		}

		if hasExclude {
			continue
		}

		log.Infof("Getting \"%s\" resource...\n", threat.FieldByName("Category").String())

		req, _ := http.NewRequest("GET", threat.FieldByName("URL").String(), nil)
		resp, _ := client.Do(req)

		body, _ := ioutil.ReadAll(resp.Body)
		threat.FieldByName("Content").SetString(string(body))
	}
}
