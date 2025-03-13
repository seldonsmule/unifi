package unifi


import (
	//"os"
	"fmt"
	//"flag"
	//"strings"
        "encoding/json"
	//"math"
       //"github.com/seldonsmule/powerwall"
       //"github.com/seldonsmule/simpleconffile"
        //"time"
        "github.com/seldonsmule/restapi"
        "github.com/seldonsmule/logmsg"

)

type Unifi struct {

  sToken string // User created personal token to ST account

  sApiKey string // From the Unifi console

  sSiteId string // Current SiteID for console

  sBaseEndpoint string // Base ST endpoint

  Devices UnDevices // List of all devices from the hub
  DeviceStats UnDeviceStats // Stats for the last called device
  Clients UnClients // List of all scenes from the hub
  Sites   UnSites // List of all scenes from the hub

}

func New(URL string) *Unifi {

  un := new(Unifi)

  un.sBaseEndpoint = URL + "/proxy/network/integrations/v1"

  return(un)

}

func (pUn *Unifi) Dump(){

  //fmt.Printf("Unifi.sToken[%s]\n", pUn.sToken)
  fmt.Printf("Unifi.sBaseEndpoint[%s]\n", pUn.sBaseEndpoint)

}

func (pUn *Unifi) SetApiKey(sApiKey string) bool{

  pUn.sApiKey = sApiKey

  return true
}

func (pUn *Unifi) SetSiteId(sSiteId string) bool{

  pUn.sSiteId = sSiteId

  return true
}

func (pUn *Unifi) GetStructs() bool{

  fmt.Println("Getting base structs to save off")

  if(!pUn.GetStructSites()){
    logmsg.Print(logmsg.Error,"GetStructSites() failed")
    return false
  }

  if(!pUn.GetStructDevices(pUn.sSiteId)){
    logmsg.Print(logmsg.Error,"GetStructDevices() failed")
    return false
  }

  // lets get the struct for getting stats too while here

  if(pUn.Devices.Count > 0){

    fmt.Println("We have Adopted devices - saving off stats struct too")
    if(!pUn.GetStructDeviceStats(pUn.sSiteId, pUn.Devices.Data[0].ID)){
      logmsg.Print(logmsg.Error,"GetStructDevices() failed")
      return false
    }

  }

  if(!pUn.GetStructClients(pUn.sSiteId)){
    logmsg.Print(logmsg.Error,"GetStructClients() failed")
    return false
  }

  return true

}

func (pUn *Unifi) ListSitesIDs() bool{

  endpointname := pUn.sBaseEndpoint + "/sites"

  r := restapi.NewGet("listsites", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  //r.JsonOnly()

  //r.Dump()

  //r.DebugOn()


  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }


  tmp1 := r.GetValue("data")
  //fmt.Printf("data[%s]\n", tmp1)

  myarray := restapi.CastArray(tmp1)

  //fmt.Println("array len:", len(myarray))

  for i:=0; i < len(myarray); i++ {

    tmpmap := restapi.CastMap(myarray[i])

    fmt.Printf("[%d]\n", i)

    for name, value := range tmpmap{
 
      fmt.Printf("    %s = %s\n", name, value)
    }


  }


  return true

}


func (pUn *Unifi) GetStructSites() bool{
  return(pUn.GetSites(true))
}

func (pUn *Unifi) GetSites(bSave bool) bool{

  endpointname := pUn.sBaseEndpoint + "/sites"

  r := restapi.NewGet("getsites", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  if(bSave){
    r.SaveResponseBody("un_sites", "UnSites", false)
  }

  // cheating and saving off

  json.Unmarshal(r.BodyBytes, &pUn.Sites)

  return true

}

func (pUn *Unifi) GetStructDevices(SiteID string ) bool{

  return(pUn.GetDevices(SiteID, true))
}


func (pUn *Unifi) GetDevices(SiteID string, bSave bool ) bool{

  endpointname := pUn.sBaseEndpoint + "/sites/" + SiteID + "/devices"

  r := restapi.NewGet("getdevices", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  if(bSave){
    r.SaveResponseBody("un_devices", "UnDevices", false)
  }

  // cheating and saving off

  json.Unmarshal(r.BodyBytes, &pUn.Devices)

  return true

}

func (pUn *Unifi) GetStructDeviceStats(SiteID string , DeviceID string) bool{

  return(pUn.GetDeviceStats(SiteID, DeviceID, true))
  
}

func (pUn *Unifi) GetDeviceStats(SiteID string , DeviceID string, bSave bool) bool{

  endpointname := pUn.sBaseEndpoint + "/sites/" + SiteID + "/devices/" + DeviceID + "/statistics/latest"

  r := restapi.NewGet("getdevicesstats", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  if(bSave){
    r.SaveResponseBody("un_devices_stats", "UnDeviceStats", false)
  }

  // cheating and saving off

  json.Unmarshal(r.BodyBytes, &pUn.DeviceStats)

  return true

}
func (pUn *Unifi) GetStructClients(SiteID string ) bool{
  return(pUn.GetClients(SiteID, true))
}

func (pUn *Unifi) GetClients(SiteID string, bSave bool) bool{

  endpointname := pUn.sBaseEndpoint + "/sites/" + SiteID + "/clients?limit=1000"

  r := restapi.NewGet("getclients", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  if(bSave){
    r.SaveResponseBody("un_clients", "UnClients", false)
  }

  // cheating and saving off

  json.Unmarshal(r.BodyBytes, &pUn.Clients)

  return true

}

func (pUn *Unifi) ListDevices(SiteID string ) bool{

  endpointname := pUn.sBaseEndpoint + "/sites/" + SiteID + "/devices"

  r := restapi.NewGet("getdevices", endpointname)

  r.SetApiKey(pUn.sApiKey)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  r.SaveResponseBody("un_devices", "UnDevices", false)

  return true

}

