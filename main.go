package main

import (
	"fmt"
)


func main(){
  watcher,err := fsnotify.NewWatcher()
  if err!= nil {
    fmt.Println(err)
    return
  }
  defer watcher.Close()
  
  err = watcher.Add("/home/student07")
  if err != nil {
    fmt.Println(err)
    return
  } 

  for {
        select{
    case event,ok := <-watcher.Events:
    if !ok {
      return
    }
    if event.Op&fsnotify.Write == fsnotify.Write {
     fmt.Println("Fichier modifié", event.Name) 
    }
    case err,ok := <- watcher.Errors:
    if !ok {
      return
    }
    fmt.Println("Erreur de surveillance:",err)
  }
  }


}
