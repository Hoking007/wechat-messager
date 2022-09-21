package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
)

var closeIDS = sync.Map{}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Notifier interface with event delivery
	notifier, err := notify.New(
		conn,
		// action event handler
		notify.WithOnAction(func(action *notify.ActionInvokedSignal) {
			log.Printf("ActionInvoked: %v Key: %v", action.ID, action.ActionKey)
		}),
		// closed event handler
		notify.WithOnClosed(func(closer *notify.NotificationClosedSignal) {
			log.Printf("NotificationClosed: %v Reason: %v", closer.ID, closer.Reason)
			closeIDS.Store(closer.ID, true)
		}),
		// override with custom logger
		notify.WithLogger(log.New(os.Stdout, "notify: ", log.Flags())),
	)

	if err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.kde.StatusNotifierItem"),
		dbus.WithMatchMember("NewIcon"),
	); err != nil {
		panic(err)
	}
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	var prevId uint32 = 0
	for v := range c {
		pid := getSenderPID(conn, v.Sender)
		procName := processName(pid)
		if isXembedsniproxy(procName) && (prevId == 0 || notifyClose(prevId)) {
			// notify
			prevId = notifyWechat(notifier, v)
		}
	}
}

func listAllNames(conn *dbus.Conn) []string {
	var s []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get list of owned names:", err)
		os.Exit(1)
	}
	fmt.Println("Currently owned names on the session bus:")
	for _, v := range s {
		fmt.Println(v)
	}
	return s
}

func getSenderPID(conn *dbus.Conn, name string) uint32 {
	var s uint32
	err := conn.BusObject().Call("org.freedesktop.DBus.GetConnectionUnixProcessID", 0, name).Store(&s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get GetConnectionUnixProcessID:", err)
		os.Exit(1)
	}
	return s
}

func processName(pid uint32) string {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return ""
	}
	name, _ := p.Name()
	return name
}

func isXembedsniproxy(name string) bool {
	return "xembedsniproxy" == name
}

func notifyWechat(notifier notify.Notifier, signal *dbus.Signal) uint32 {
	body := "你收到了新的消息"
	if signal.Body != nil {
		body = body + "   " + fmt.Sprint(signal.Body)
	}
	n := notify.Notification{
		AppName:       "微信",
		ReplacesID:    uint32(0),
		AppIcon:       "mail-unread",
		Summary:       "通知消息",
		Body:          body,
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: time.Second * 10,
	}
	// Ship it!
	createdID, err := notifier.SendNotification(n)
	if err != nil {
		log.Printf("error sending notification: %v", err.Error())
	}
	log.Printf("created notification with id: %v", createdID)
	return createdID
}
func notifyClose(prevId uint32) bool {
	_, exist := closeIDS.LoadAndDelete(prevId)
	return exist
}
