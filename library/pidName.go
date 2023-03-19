package library

	

func getNamePID(){
	pid := process.GetSelPid()
	p,err:= process.Newprocess(pid)
	if err != nil {
		fmt.Println(err)
	}
	name,err:=p.Name()
	if err != nil {
		fmt.Println(err)
	}
fmt.Println("le Processus en cours d'éxécution est: %s\n",name)
}


