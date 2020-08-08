package runner

import (
	"os"
	"reflect"
	"strings"

	"gopkg.in/validator.v2"
	"ktbs.me/teler/common"
	"ktbs.me/teler/pkg/errors"
	"ktbs.me/teler/pkg/matchers"
	"ktbs.me/teler/pkg/parsers"
)

func validate(options *common.Options) {
	if !options.Stdin {
		if options.Input == "" {
			errors.Exit(errors.ErrNoInputLog)
		}
	}

	if options.ConfigFile == "" {
		telerEnv := os.Getenv("TELER_CONFIG")
		if telerEnv == "" {
			errors.Exit(errors.ErrNoInputConfig)
		} else {
			options.ConfigFile = telerEnv
		}
	}

	config, errConfig := parsers.GetConfig(options.ConfigFile)
	if errConfig != nil {
		errors.Exit(errConfig.Error())
	}

	// Validates log format
	matchers.IsLogformat(config.Logformat)
	options.Configs = config

	// Validates notification parts on configuration files
	notification(options)

	if errVal := validator.Validate(options); errVal != nil {
		errors.Exit(errVal.Error())
	}
}

func notification(options *common.Options) {
	config := options.Configs

	if config.Alert.Active {
		provider := strings.Title(config.Alert.Provider)
		field := reflect.ValueOf(&config.Notifications).Elem().FieldByName(provider)

		switch provider {
		case "Slack":
			field.FieldByName("URL").SetString(SlackAPI)
			matchers.IsHexcolor(field.FieldByName("Color").String())
			matchers.IsChannel(field.FieldByName("Channel").String())
		case "Telegram":
			field.FieldByName("URL").SetString(strings.Replace(TelegramAPI, ":token", field.FieldByName("Token").String(), -1))
			matchers.IsChatID(field.FieldByName("ChatID").String())
			matchers.IsParseMode(field.FieldByName("ParseMode").String())
		default:
			errors.Exit(strings.Replace(errors.ErrAlertProvider, ":platform", config.Alert.Provider, -1))
		}

		matchers.IsToken(field.FieldByName("Token").String())
	}
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
