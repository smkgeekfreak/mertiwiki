// util
package util

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	//"os/signal"
	//"reflect"
	//"sync"
	//"syscall"
	"fmt"
	"path/filepath"
	"runtime"
)

type TestCounter struct {
	SuccessCount, FailCount int
}

/*
 Get package and name stirng from calling func
*/
func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

type Config struct {
	Port                   string
	ApiVersion             string
	RootPath               string
	Swagger_WebServicesUrl string
	SwaggerFilePath        string
	ConfigRoot             string
}

var (
	config *Config
	//configLock = new(sync.RWMutex)
)

func loadConfig(fail bool) {
	wd, _ := os.Getwd()
	log.Printf("WORKING DIR %s", wd)
	cfgFile := filepath.Join("./", "config.json")
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	log.Printf("config = %#v", temp)
	//configLock.Lock()
	config = temp

	//configLock.Unlock()

	err = os.Setenv("WIKICONFIGROOT", config.ConfigRoot)
	if err != nil {
		log.Fatalf("Could not set WIKICONFIGROOT - %#v", err)
	}
}

func GetConfig() *Config {
	//configLock.RLock()
	//defer configLock.RUnlock()
	return config
}
func init() {
	loadConfig(false)
	//s := make(chan os.Signal, 1)
	//signal.Notify(s, syscall.SIGUSR2)
	//go func() {
	//	for {
	//		<-s
	//		loadConfig(false)
	//		log.Println("Reloaded")
	//	}
	//}()
}

/*
Reflectively call the specified method on the provided interface
TODO: deprecated .. remove
Yes I should move
*/
//func InvokeModel(any interface{}, name string, args ...interface{}) (result []reflect.Value) { //interface{} {
//inputs := make([]reflect.Value, len(args))
//for i, _ := range args {
//	inputs[i] = reflect.ValueOf(args[i])
//}
//retval := reflect.ValueOf(any).MethodByName(name).Call(inputs)
////log.Printf("return from model invoke %#v", retval)
//result = retval
//return
//}

func ReadFlags() {
	//
	// Read command line flags
	wikiEnvPtr := flag.String("wikiEnv", "dev", "Setting the configuration environment")
	flag.Parse()
	log.Printf("Staring wikiServer with -wikiEnv %#v", *wikiEnvPtr)
	//
	// Set the environment variable for configuration properites
	err := os.Setenv("WIKIENV", *wikiEnvPtr)
	if err != nil {
		log.Fatalf("Could not set WIKENV - %#v", err)
	}
}

func toString(a interface{}) (string, bool) {
	aString, isString := a.(string)
	if isString {
		return aString, true
	}

	aBytes, isBytes := a.([]byte)
	if isBytes {
		return string(aBytes), true
	}

	aStringer, isStringer := a.(fmt.Stringer)
	if isStringer {
		return aStringer.String(), true
	}

	return "", false
}
