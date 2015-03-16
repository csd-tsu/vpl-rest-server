package main
import (
    "os"
    "flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"os/exec"
	"encoding/json"
)


type HandlerConfigLine struct {
    UrlPattern string
    Executable string
    Argument string
}

func CreateHandler(config HandlerConfigLine) func(http.ResponseWriter,*http.Request) {
    log.Println("Registering ", config.Executable, config.Argument);

    return func(w http.ResponseWriter, r *http.Request) {
        cmd:= exec.Command(config.Executable, config.Argument)
    
        cmd.Stdin = r.Body
        stdout,_ := cmd.StdoutPipe()
        cmd.Start()
    
        result, _ := ioutil.ReadAll(stdout)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        fmt.Fprintf(w, "%s", result)
    }
}
    

func ReadConfig(filename string) (configuration []HandlerConfigLine, err error) {
    file, err := os.Open(filename)
    if (err != nil) {
        return
    }
    decoder := json.NewDecoder(file);
    err = decoder.Decode(&configuration)
    return

}

func main() {
    var configFilename string
    flag.StringVar(&configFilename, "config", "", "path to a configuration file")
    flag.Parse()

    if(configFilename == "") {
        log.Fatal("Please use --config flag to define a path to a configuration file")
    }

    handlerLines, err := ReadConfig(configFilename)
    
    if (err !=nil) {
        log.Fatal(err)
    }

    for _, line := range handlerLines {
        http.HandleFunc(line.UrlPattern, CreateHandler(line))
    }

    log.Println("Running...")
    err = http.ListenAndServe(":8080", nil);
    if (err != nil){
        log.Fatal(err);
    }
}
