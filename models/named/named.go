package named

import (
    "fmt"
    "bytes"
    "sync"
    "os/exec"
    "io"
    "text/template"
    "strings"

    netaddr "github.com/dspinhirne/netaddr-go"

    "github.com/NCKU-NASA/nasa-judge-lib/schema/user"

    "github.com/NCKU-NASA/nasa-judge-named/utils/config"
)

var lock *sync.RWMutex

type Record struct {
    Method string `json:"method"`
    Name string `json:"name"`
    Type string `json:"type"`
    Data string `json:"data"`
}

func init() {
    lock = new(sync.RWMutex)
    users, err := user.GetUsers()
    if err != nil {
        panic(err)
    }
    for _, now := range users {
        SetRecord(now)
    }
}

func nthnet(network string, index uint) string {
    net, err := netaddr.ParseIPv4Net(network)
    if err != nil {
        panic(err)
    }
    return net.Nth(uint32(index)).String()
}

func (c Record) Set() {
    lock.Lock()
    defer lock.Unlock()
    cmd := exec.Command("nsupdate", "-k", config.UpdateKey)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        panic(err)
    }
    io.WriteString(stdin, fmt.Sprintf(`server 127.0.0.1
update %s %s 604800 %s %s
send
quit`, c.Method, c.Name, c.Type, c.Data))
    stdin.Close()
    if err = cmd.Start(); err != nil {
        panic(err)
    }
}

func SetRecord(userdata user.User) {
    tmplfunc := template.FuncMap{
        "nthnet": nthnet,
    }
    for _, tmpl := range config.RecordTmpl {
        t := template.New("").Funcs(tmplfunc)
        t = template.Must(t.Parse(tmpl))
        var buf bytes.Buffer
        t.Execute(&buf, userdata)
        recordpart := strings.Split(buf.String(), " ")
        if len(recordpart) != 3 || recordpart[0] == "" || recordpart[1] == "" || recordpart[2] == "" {
            continue
        }
        nowrecord := Record{
            Method: "add",
            Name: recordpart[0],
            Type: recordpart[1],
            Data: recordpart[2],
        }
        nowrecord.Set()
    }
}

