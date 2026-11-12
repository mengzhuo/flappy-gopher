package main

type soundEvent uint8

const (
	eventFlap soundEvent = iota
	eventScore
	eventHit
)
