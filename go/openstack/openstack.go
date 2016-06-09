package main

import (
  "os"
  "fmt"
  "regexp"
  "io/ioutil"
  "encoding/json"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

type ControllerObj struct {
  IP string `json:"ip"`
}

type ComputeObj struct {
  Name string `json:"name"`
  OriginalIP string `json:"original_ip"`
  IP string `json:"ip"`
  TunnelIP string `json:"tunnel_ip"`
  EthNeutron string `json:"eth-neutron"`
  EthCeph string `json:"eth-ceph"`
}

type Machines struct {
  Controller ControllerObj
  Compute []ComputeObj
}

func ReplaceAllSubmatchIndex(re *regexp.Regexp, str string, replaceWith string) string {
  result := ""
  lastIndex := 0

  for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
    groups := []string{}
    for i:=0; i<len(v); i+=2 {
      groups = append(groups, str[v[i]:v[i+1]])
    }
    result += str[lastIndex:v[0]] + groups[1] + replaceWith
    lastIndex = v[1]
  }
  return result + str[lastIndex:]
}

func main() {
  DIRECTORY := "/Users/ccoleman/it-ops/openstack/liberty/"

  machines_raw, err := ioutil.ReadFile(DIRECTORY + "machines.json")
  check(err)

  machines := &Machines{}
  if err := json.Unmarshal(machines_raw, &machines); err != nil {
    panic(err)
  }

  compute_template_raw, err := ioutil.ReadFile(DIRECTORY + "template/openstack-compute/config")
  check(err)
  compute_template := string(compute_template_raw)

  var reController = regexp.MustCompile(`(CONTROLLER_IP=)(.*)`)
  var reThisHostName = regexp.MustCompile(`(THISHOST_NAME=)(.*)`)
  var reThisHostIP = regexp.MustCompile(`(THISHOST_IP=)(.*)`)
  var reThisHostTunnel = regexp.MustCompile(`(THISHOST_TUNNEL_IP=)(.*)`)
  var reAvailabilityZone = regexp.MustCompile(`(DEFAULT_AZ=)(.*)`)

  fmt.Println(machines.Controller.IP)
  for _, computeNode := range machines.Compute {
    fmt.Println(computeNode.Name)
    os.Mkdir(DIRECTORY + "/hosts/" + computeNode.Name, 0777)

    compute_template = ReplaceAllSubmatchIndex(reController, compute_template, machines.Controller.IP)
    compute_template = ReplaceAllSubmatchIndex(reThisHostName, compute_template, computeNode.Name)
    compute_template = ReplaceAllSubmatchIndex(reThisHostIP, compute_template, computeNode.IP)
    compute_template = ReplaceAllSubmatchIndex(reThisHostTunnel, compute_template, computeNode.TunnelIP)
    compute_template = ReplaceAllSubmatchIndex(reAvailabilityZone, compute_template, "\"Uncharted Software Toronto Production\"")

    err := ioutil.WriteFile(DIRECTORY + "/hosts/" + computeNode.Name + "/config", []byte(compute_template), 0644)
    check(err)

    install := "scp -r template/openstack-compute/* root@"+computeNode.OriginalIP+":\nscp -r template/common root@"+computeNode.OriginalIP+":\nscp -r template/common.sh root@"+computeNode.OriginalIP+":\nscp -r " + DIRECTORY + "/hosts/" + computeNode.Name + "/config  root@"+computeNode.OriginalIP+":\nssh root@"+computeNode.OriginalIP+" sh install-common-and-network.sh "+computeNode.EthNeutron+" "+computeNode.EthCeph
    err = ioutil.WriteFile(DIRECTORY + "/hosts/" + computeNode.Name + "/install.sh", []byte(install), 0744)
    check(err)
  }
}
