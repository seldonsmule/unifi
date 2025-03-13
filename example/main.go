package main

import (

  "fmt"
  "flag"
  "os"
  "strings"
  //"/Users/rxe789/egc_NetworkDrive/dev/gocode/src/github.com/seldonsmule/unifi"
  "github.com/seldonsmule/unifi"
  "github.com/seldonsmule/simpleconffile"
  "github.com/seldonsmule/logmsg"

)

type Configuration struct {

  API_KEY string // Key from Unfi
  ConfFilename string // name of conf file
  Host string // name of the Unifi device
  Siteid string // default site id
  Encrypted bool

}

const COMPILE_IN_KEY = "example key 9999"

var gMyConf Configuration

func readconf(confFile string, printstd bool) bool{

  simple := simpleconffile.New(COMPILE_IN_KEY, confFile)
  
  if(!simple.ReadConf(&gMyConf)){
    msg := fmt.Sprintln("Error reading conf file: ", confFile)
    logmsg.Print(logmsg.Warning, msg)
    return false
  }
  
  if(gMyConf.Encrypted){    
    gMyConf.API_KEY = simple.DecryptString(gMyConf.API_KEY)
  }


  if(printstd){

    fmt.Printf("Encrypted [%v]\n", gMyConf.Encrypted)
    fmt.Printf("API_KEY [%v]\n", gMyConf.API_KEY)
    fmt.Printf("Host [%v]\n", gMyConf.Host)
    fmt.Printf("SiteID [%v]\n", gMyConf.Siteid)
    fmt.Printf("ConfFilename [%v]\n", gMyConf.ConfFilename)

  }

  return true

}

func help(){

  fmt.Println("cmd not found")

  flag.PrintDefaults()

  fmt.Println("cmds:")
  fmt.Println("      setconf - Setup Conf file")
  fmt.Println("        -apikey Unifi API Key")
  fmt.Println("        -host URL of Unifi device")
  fmt.Println("        -siteid Stored ID for APIs - see listsiteids")
  fmt.Println("        -conffile name of conffile (.unifi.conf default)")
  fmt.Println()
  fmt.Println("      readconf - Display Conf file info")
  fmt.Println("      listsiteids - List all the site IDs for the host URL")



}

func main(){

  cmdPtr := flag.String("cmd", "help", "Command to run")
  apikeyPtr := flag.String("apikey", "notset", "Unifi API KEY")
  hostPtr := flag.String("host", "notset", "URL of Host system")
  sitePtr := flag.String("siteid", "notset", "Default Site ID (needed for APIs)")
  confPtr := flag.String("conffile", ".unifi.conf", "config file name")
  bdebugPtr := flag.Bool("debug", false, "If true, do debug magic")

  flag.Parse()

  fmt.Printf("cmd=%s\n", *cmdPtr)

  logmsg.SetLogFile("example.log");

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "apikeyPtr = ", *apikeyPtr)
  logmsg.Print(logmsg.Info, "hostPtr = ", *hostPtr)
  logmsg.Print(logmsg.Info, "sitePtr = ", *sitePtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "bdebugPtr = ", *bdebugPtr)


  fmt.Println("starting");

  readconf(*confPtr, false);

  un := unifi.New(gMyConf.Host)

  //un.Dump()

  switch *cmdPtr {

    case "readconf":
      fmt.Println("Reading Conf File")
      readconf(*confPtr, true)

    case "setconf":

      readconf(*confPtr, false); // ignore errors

      fmt.Println("Setting conf file")

      simple := simpleconffile.New(COMPILE_IN_KEY, *confPtr);

      gMyConf.Encrypted = true

      if(strings.Compare(*apikeyPtr, "notset") != 0){
       gMyConf.API_KEY = simple.EncryptString(*apikeyPtr)
      }else{
       gMyConf.API_KEY = simple.EncryptString(gMyConf.API_KEY)
      } 

      if(strings.Compare(*hostPtr, "notset") != 0){
        gMyConf.Host = *hostPtr
      }

      if(strings.Compare(*sitePtr, "notset") != 0){
        gMyConf.Siteid = *sitePtr;
      }
     
      gMyConf.ConfFilename = *confPtr


      simple.SaveConf(gMyConf)
 
    case "listsiteids":
     fmt.Printf("Site IDs for Host[%s]\n", gMyConf.Host)
     un.SetApiKey(gMyConf.API_KEY)
     un.ListSitesIDs()

    case "get_structs":

    // un.SetApiKey(gMyConf.API_KEY)
     un.SetApiKey(gMyConf.API_KEY)
     un.SetSiteId(gMyConf.Siteid)
     un.GetStructs()
     //un.GetStructDevices(gMyConf.Siteid)
    

    default:
      help()
      os.Exit(2)

  }


  
}
