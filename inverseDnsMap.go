package main

import "github.com/miekg/dns"

type inverseDnsMap struct {
	ipMap    map[string][]string
	cNameMap map[string][]string
}

func newInverseDnsMap() *inverseDnsMap {
	return &inverseDnsMap{
		ipMap:    make(map[string][]string),
		cNameMap: make(map[string][]string),
	}
}

func (this *inverseDnsMap) Add(byteMessage []byte) {
	msg := new(dns.Msg)
	err := msg.Unpack(byteMessage)
	if err != nil {
		return
	}

	for _, answer := range msg.Answer {
		switch answer.Header().Rrtype {
		case dns.TypeA:
			if t, ok := answer.(*dns.A); ok {

				this.ipMap[t.A.String()] = AppendUnique(this.ipMap[t.A.String()], t.Hdr.Name)
			}
			break
		case dns.TypeCNAME:
			if t, ok := answer.(*dns.CNAME); ok {
				this.cNameMap[t.Target] = AppendUnique(this.cNameMap[t.Target], t.Hdr.Name)
			}
		default:
		}
	}
}

func AppendUnique(slice []string, item string) []string {
	for _, element := range slice {
		if element == item {
			return slice
		}
	}
	return append(slice, item)
}

func (this *inverseDnsMap) Get(ip string) []string {
	var stringList []string
	for _, domain := range this.ipMap[ip] {
		stringList = append(stringList, domain)
		stringList = append(stringList, this.getCnames(domain)...)

	}

	return stringList
}

func (this *inverseDnsMap) getCnames(domain string) []string {
	var stringList []string
	for _, domain := range this.cNameMap[domain] {
		stringList = append(stringList, domain)
		stringList = append(stringList, this.getCnames(domain)...)
	}
	return stringList
}
