package config

import (
    "net"
    "os"
    "strconv"
    "strings"
    "encoding/json"
)

var Debug bool
var Trust []string
var Port string
var Secret string
var Sessionname string
var RecordTmpl []string
var UpdateKey string

func init() {
    loadenv()
    var err error
    debugstr, exists := os.LookupEnv("DEBUG")
    if !exists {
        Debug = false
    } else {
        Debug, err = strconv.ParseBool(debugstr)
        if err != nil {
            Debug = false
        }
    }
    truststr := os.Getenv("TRUST")
    if truststr == "" {
        Trust = []string{"127.0.0.1", "::1"}
    } else {
        var tmp []string
        err = json.Unmarshal([]byte(truststr), &tmp)
        if err != nil {
            panic(err)
        }
        for _, now := range tmp {
            ips, _ := net.LookupIP(now)
            for _, ip := range ips {
                Trust = append(Trust, ip.String())
            }
        }
    }
    RecordTmpl = strings.Split(os.Getenv("RECORDTMPL"), "\n")
    for idx, tmpl := range RecordTmpl {
        RecordTmpl[idx] = strings.Trim(tmpl, " \t\n")
    }
    Port = os.Getenv("PORT")
    Secret = os.Getenv("SECRET")
    Sessionname = os.Getenv("SESSIONNAME")
    UpdateKey = os.Getenv("UPDATEKEY")
}
