package snmp

import (
	"fmt"
)

const (
	trapExchange         = "SnmpTrapExchanger"
	notificationExchange = "NotificationExchange"
)

var (
	// notificationConverterMapping is a mapping from notification to it's converter.
	// The Address in the notification will be taken as the key here, since only this information is contained in the notification for sure.
	notificationConverterMapping = make(map[string]Converter)
	// snURLMapping is a mapping from server's serial number to server's URL.
	snURLMapping = make(map[string]string)

	huaweiNotificationConverter = &HuaweiNotificationConverter{}
	dellNotificationConverter   = &DellNotificationConverter{}
	hpeNotificationConverter    = &HpeNotificationConverter{}
	ibmNotificationConverter    = &IBMNotificationConverter{}
	lenovoNotificationConverter = &LenovoNotificationConverter{}
)

// Init includes all kinds of work that should be done before notification processing, it includes and should be performed in the order below:
// 1. Prepare the message queue so that when the notification comes we can put it into immediately.
// 2. Prepare the mapping from the notification to the converter.
func Init() {
	generateConverterMapping()
}

func generateConverterMapping() {
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = huaweiNotificationConverter
		snURLMapping[sn] = "/api/v1/servers/" + sn
		sn = fmt.Sprintf("sn-hpe-%d", i)
		notificationConverterMapping[sn] = dellNotificationConverter
		snURLMapping[sn] = "/api/v1/servers/" + sn
		sn = fmt.Sprintf("sn-dell-%d", i)
		notificationConverterMapping[sn] = dellNotificationConverter
		snURLMapping[sn] = "/api/v1/servers/" + sn
		sn = fmt.Sprintf("sn-ibm-%d", i)
		notificationConverterMapping[sn] = dellNotificationConverter
		snURLMapping[sn] = "/api/v1/servers/" + sn
		sn = fmt.Sprintf("sn-lenovo-%d", i)
		notificationConverterMapping[sn] = dellNotificationConverter
		snURLMapping[sn] = "/api/v1/servers/" + sn
	}
}