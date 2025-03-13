package main

import (

  "fmt"
  "flag"
  "os"
  "strings"
  "sort"
  "github.com/seldonsmule/unifi"
  "github.com/seldonsmule/simpleconffile"
  "github.com/seldonsmule/logmsg"

)

type Configuration struct {

  API_KEY string // Key from Unfi
  ConfFilename string // name of conf file
  Host string // name of the Unifi device
  SiteId string // default site id
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
    fmt.Printf("SiteID [%v]\n", gMyConf.SiteId)
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
  fmt.Println("      listdevices - List all the Uniquit devices (using siteid)")
  fmt.Println("        -csv Causes output to have a CSV format")
  fmt.Println("        -etchosts Causes output to be ready to cut/paste into /etc/hosts")
  fmt.Println("      listclients - List all the clients (using siteid)")
  fmt.Println("        -csv Causes output to have a CSV format")
  fmt.Println("        -etchosts Causes output to be ready to cut/paste into /etc/hosts")



}

func main(){

  cmdPtr := flag.String("cmd", "help", "Command to run")
  apikeyPtr := flag.String("apikey", "notset", "Unifi API KEY")
  hostPtr := flag.String("host", "notset", "URL of Host system")
  sitePtr := flag.String("siteid", "notset", "Default Site ID (needed for APIs)")
  confPtr := flag.String("conffile", ".unifi.conf", "config file name")
  bdebugPtr := flag.Bool("debug", false, "If true, do debug magic")
  bcsvPtr := flag.Bool("csv", false, "If true, separate columns with ','")
  betchostsPtr := flag.Bool("etchosts", false, "If true, Print a format to be copied into a /etc/hosts")


  flag.Parse()

  //fmt.Printf("cmd=%s\n", *cmdPtr)

  logmsg.SetLogFile("unifi_print.log");

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "apikeyPtr = ", *apikeyPtr)
  logmsg.Print(logmsg.Info, "hostPtr = ", *hostPtr)
  logmsg.Print(logmsg.Info, "sitePtr = ", *sitePtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "bdebugPtr = ", *bdebugPtr)
  logmsg.Print(logmsg.Info, "bcsvPtr = ", *bcsvPtr)
  logmsg.Print(logmsg.Info, "betchostsPtr = ", *betchostsPtr)

  var delimiter int

  if(*bcsvPtr){
    delimiter = ','
  }else{
    delimiter = '\t'
  }

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
        gMyConf.SiteId = *sitePtr;
      }
     
      gMyConf.ConfFilename = *confPtr


      simple.SaveConf(gMyConf)
 
    case "listsiteids":
     fmt.Printf("Site IDs for Host[%s]\n", gMyConf.Host)
     un.SetApiKey(gMyConf.API_KEY)
     un.GetSites(false)

     fmt.Printf("%s Sites\n\n", gMyConf.Host)

     c := delimiter // just being lazy with a short variable name

     fmt.Printf("Count%c Name%c ID\n", c, c)

     for i := 0; i < int(un.Sites.Count); i++ {
     
       fmt.Printf("%d%c %s%c %s\n", 
                  i,
                  c,
                  un.Sites.Data[i].Name,
                  c,
                  un.Sites.Data[i].ID)

     }

    case "listdevices":
     fmt.Printf("Devices for Host[%s]\n", gMyConf.Host)
     un.SetApiKey(gMyConf.API_KEY)
     un.GetDevices(gMyConf.SiteId, false)

     c := delimiter // just being lazy with a short variable name

     fmt.Printf("%s Devices\n\n", gMyConf.Host)

     fmt.Printf("Count%c Name%c ID%c IP%c MAC\n", c, c, c, c)

     for i := 0; i < int(un.Devices.Count); i++ {
     
       fmt.Printf("%d%c %s%c %s%c %s%c %s\n", 
                  i,
                  c,
                  un.Devices.Data[i].Name,
                  c,
                  un.Devices.Data[i].ID,
                  c,
                  un.Devices.Data[i].IpAddress,
                  c,
                  un.Devices.Data[i].MacAddress)
                  

     }

    case "listclients":
     if(!*betchostsPtr){
       fmt.Printf("Clients for Host[%s]\n", gMyConf.Host)
     }
     un.SetApiKey(gMyConf.API_KEY)
     un.GetClients(gMyConf.SiteId, false)

     c := delimiter // just being lazy with a short variable name

     if(!*betchostsPtr){
       fmt.Printf("%s Clients\n\n", gMyConf.Host)

       fmt.Printf("Count%c Name%c ID%c IP%c MAC\n", c, c, c, c)

       for i := 0; i < int(un.Clients.Count); i++ {
     
         fmt.Printf("%d%c %s%c %s%c %s%c %s\n", 
                  i,
                  c,
                  un.Clients.Data[i].Name,
                  c,
                  un.Clients.Data[i].ID,
                  c,
                  un.Clients.Data[i].IpAddress,
                  c,
                  un.Clients.Data[i].MacAddress)
                  

       }
     }else{


        sort.Slice(un.Clients.Data, func(i, j int) bool {
          return un.Clients.Data[i].IpAddress < un.Clients.Data[j].IpAddress
        })


       for i := 0; i < int(un.Clients.Count); i++ {

         if(strings.Compare(un.Clients.Data[i].Name, "") == 0){
           continue
         }

         if(strings.Compare(un.Clients.Data[i].IpAddress, "") == 0){
           continue
         }

         n := un.Clients.Data[i].Name

         name := strings.ReplaceAll(n, " ", "")



         fmt.Printf("%s\t%s\n", un.Clients.Data[i].IpAddress, name)

/*
         a := net.ParseIP(un.Clients.Data[i].IpAddress)
         fmt.Printf("%b %s\n", net.IP.To4(a))
*/
       }
     }

    case "getallinfo":

    default:
      help()
      os.Exit(2)

  }


  
}
