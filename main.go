package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const userKey = "user"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("must specify config yaml as args")
		os.Exit(1)
	}
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		conf := &Config{}
		err := conf.LoadConfig(os.Args[1])
		if err != nil {
			return false, err
		}

		err = conf.Auth.LoadUserPass(os.Args[1])
		if err != nil {
			return false, err
		}
		if conf.Auth.IsPermitted(username, password) {
			c.Set(userKey, username)
			return true, nil
		} else {
			return false, fmt.Errorf("not permitted")
		}

	}))

	// Routes
	e.GET("/ping", renew)
	e.GET("/revoke/:ip", revoke)
	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}

func revoke(c echo.Context) (err error) {
	conf := &Config{}
	ip := c.Param("ip")
	err = conf.LoadConfig(os.Args[1])
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = conf.AwsSg.Init()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//err = conf.AwsSg.RevokeSgIngress(ip)
	//if err != nil {
	//	if !strings.Contains(err.Error(), "IP not exist") {
	//		return c.String(http.StatusInternalServerError, err.Error())
	//	}
	//}
	groupIds := make([]string, 1)
	groupIds[0] = conf.AwsSg.SgId
	sgs, err := conf.AwsSg.ListSg(groupIds)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	for _, v := range sgs[0].IpPermissions {
		fmt.Println(*(v.IpProtocol))
		if *(v.IpProtocol) == "-1" {
			iphive := false
			iplens := len(v.IpRanges)
			fmt.Println("lens is : ", iplens)
			for _, t := range v.IpRanges {
				cidrips := strings.Replace(*(t.CidrIp), "/32", "", -1)
				fmt.Println("cidrips is : ", cidrips)
				if cidrips == ip {
					iphive = true
				}
			}
			if !iphive {
				fmt.Println("ip is not exist")
				return c.String(http.StatusOK, "ip is not exist")
			} else {
				fmt.Println("ip is exist")
			}
			fmt.Println(" iphive is : ", iphive)
			ips := v.IpRanges[0].CidrIp
			fmt.Println("value is : ", *ips)
			fmt.Printf("type is : %T \n", v.IpRanges[0].CidrIp)
		} else {
			fmt.Println("not in IpProtocol is : ", *(v.IpProtocol))
		}
	}
	err = conf.AwsSg.RevokeSgIngress(ip)
	if err != nil {
		if !strings.Contains(err.Error(), "IP not exist") {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.String(http.StatusOK, sgs[0].String())
}

// Handler
func renew(c echo.Context) (err error) {
	remoteIP := strings.Split(c.Request().RemoteAddr, ":")[0]
	conf := &Config{}
	err = conf.LoadConfig(os.Args[1])
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = conf.AwsSg.Init()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	desc := fmt.Sprintf("%s %s", c.Get(userKey), time.Now().Format("20060102150405"))
	err = conf.AwsSg.AuthSgIngress(remoteIP, desc)
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate") {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	groupIds := make([]string, 1)
	groupIds[0] = conf.AwsSg.SgId
	sgs, err := conf.AwsSg.ListSg(groupIds)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, sgs[0].String())
}
