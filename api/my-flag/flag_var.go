package api_flag

import "flag"

var ConfPath string

func init() {
	flag.StringVar(&ConfPath, "conf", "config.yml", "config path")
	flag.Parse()
}
