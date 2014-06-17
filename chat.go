package main

import (
    "encoding/json"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/googollee/go-socket.io"
    "log"
    "net"
    "net/http"
    "os"
)

func getName(id string, users map[string]string) string {
    oldname := id
    var present bool
    _, present = users[id]
    if present {
        oldname = users[id]
    }
    return oldname
}

func getOnlineMap(users map[string]string) string {
    slc, _ := json.Marshal(users)
    return string(slc)
}

func main() {
    m := martini.Classic()
    m.Use(render.Renderer())

    m.Use(martini.Logger())

    m.Get("/", func(r render.Render) {
        r.HTML(200, "index", nil)
    })

    users := make(map[string]string)

    sio := socketio.NewSocketIOServer(&socketio.Config{})

    sio.On("connect", func(ns *socketio.NameSpace) {
        name := getName(ns.Id(), users)
        users[ns.Id()] = name
        //log.Println("Connected: ", name)
        sio.Broadcast("chatresponse", name+" joined")
        sio.Broadcast("onlinechange", getOnlineMap(users))
    })

    sio.On("disconnect", func(ns *socketio.NameSpace) {
        name := getName(ns.Id(), users)
        //log.Println("Disonnected: ", name)
        delete(users, ns.Id())
        sio.Broadcast("chatresponse", name+" left")
        sio.Broadcast("onlinechange", getOnlineMap(users))
    })

    sio.On("register", func(ns *socketio.NameSpace, name string) {
        oldname := getName(ns.Id(), users)
        users[ns.Id()] = name
        sio.Broadcast("chatresponse", oldname+" is now known as "+name)
        sio.Broadcast("onlinechange", getOnlineMap(users))
    })

    sio.On("chat", func(ns *socketio.NameSpace, message string) {
        //log.Println("Chat: ", message)
        name := getName(ns.Id(), users)
        sio.Broadcast("chatresponse", name+": "+message)
    })

    sio.Handle("/", m)

    var listener net.Listener
    var err error
    if os.Getenv("SOCKET") != "" {
        listener, err = net.Listen("unix", os.Getenv("SOCKET"))
    } else if os.Getenv("PORT") != "" {
        listener, err = net.Listen("tcp", ":"+os.Getenv("PORT"))
    } else {
        listener, err = net.Listen("tcp", ":3000")
    }
    if err != nil {
        panic(err)
    }

    log.Fatal(http.Serve(listener, sio))
}
