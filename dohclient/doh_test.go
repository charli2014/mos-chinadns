//     Copyright (C) 2020, IrineSistiana
//
//     This file is part of mos-chinadns.
//
//     mosdns is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     mosdns is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

package dohclient

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/miekg/dns"
)

//注意：这个test不正规，而且永远会pass。手动debug用。
func Test_dohClient(t *testing.T) {
	requestLogger := logrus.NewEntry(logrus.StandardLogger())
	logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	c := NewClient("https://223.5.5.5/dns-query", "223.5.5.5:443", nil, dns.MaxMsgSize, time.Second*3)
	q := new(dns.Msg)
	q.SetQuestion("www.baidu.com.", dns.TypeA)

	start := time.Now()
	fmt.Printf("First query\n")
	r, err := c.Exchange(q, requestLogger)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", r)
	fmt.Printf("query 0: %s\n", time.Since(start))

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		n := i
		go func() {
			defer wg.Done()
			start := time.Now()
			_, err := c.Exchange(q, requestLogger)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("query %d: %s\n", n, time.Since(start))
		}()
	}
	wg.Wait()
	fmt.Printf("total: %s\n", time.Since(start))
}
