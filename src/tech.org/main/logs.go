package main

import "log"

const LOG_ENABLED  = true

func logs(msg string){

	if LOG_ENABLED {
		println(msg)
	}
}


func logs_f(msg string, args interface{}){

	if LOG_ENABLED {
		log.Printf(msg, args)
	}
}
