package aws

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSecurityGroup_ListSg(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	sg := &SecurityGroup{
		Region: "cn-northwest-1",
		SgName: "Hadoop",
		SgId:   "sg-0e0c5cd076cf1fb51",
	}
	Convey("List Security Group", t, func() {
		err := sg.Init()
		So(err, ShouldBeNil)
		sgs, err := sg.ListSg([]string{sg.SgId})
		So(err, ShouldBeNil)
		So(sgs, ShouldNotBeNil)
	})
	Convey("Auth new ip", t, func() {
		testIP := "1.1.1.1"
		_ = sg.RevokeSgIngress(testIP)
		err := sg.AuthSgIngress(testIP, "unit test")
		So(err, ShouldBeNil)
		desc := fmt.Sprintf("unittest %s", time.Now().Format("20060102150405"))
		err = sg.AuthSgIngress(testIP, desc)
		So(err.Error(), ShouldContainSubstring, "Duplicate")
		sgs, err := sg.ListSg([]string{sg.SgId})
		So(err, ShouldBeNil)
		var found bool
	findIP:
		for _, group := range sgs {
			for _, perm := range group.IpPermissions {
				if *perm.IpProtocol == "-1" {
					for _, ipr := range perm.IpRanges {
						if strings.HasPrefix(*ipr.CidrIp, testIP) {
							found = true
							break findIP
						}
					}
				}
			}
		}
		So(found, ShouldBeTrue)

		found = false
		err = sg.RevokeSgIngress(testIP)
		sgs, err = sg.ListSg([]string{sg.SgId})
	findIP2:
		for _, group := range sgs {
			for _, perm := range group.IpPermissions {
				if *perm.IpProtocol == "-1" {
					for _, ipr := range perm.IpRanges {
						if strings.HasPrefix(*ipr.CidrIp, testIP) {
							found = true
							break findIP2
						}
					}
				}
			}
		}
		So(found, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}
