package main
import ( 
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"os/exec"
)



func handler(w http.ResponseWriter, r *http.Request) {
    cmd:= exec.Command("cat")
    
    cmd.Stdin = r.Body
    
    stdout,_ := cmd.StdoutPipe()
    
    
    cmd.Start()
    
    
    result,_ := ioutil.ReadAll(stdout)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w, "%s", result)
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":8080", nil);
    
    if (err != nil){
	log.Fatal(err);
    }
}
