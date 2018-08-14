package cpu

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"sort"

	"github.com/shirou/gopsutil/process"
)

type CPU struct {
	CPUByProcess   []ProcessCPU `json:"cpu_by_process"`
	Debug bool
}

type ProcessCPU struct {
	Name           string  `json:"name"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
	Chars          [10]int
	Rank           int
}

func (p *ProcessCPU) Encode() {
  for i, _ := range p.Name {
    p.Chars[i] = int([]rune(p.Name)[i])
  }
}

func (c *CPU) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()

	c.CPUByProcess = nil
	reversed_freq := map[float64][]ProcessCPU{}

	ps, err := process.Processes()
	if err != nil {
		log.Println(err)
	}
	for _, proc := range ps {
		name := fmt.Sprintf("/proc/%v", proc.Pid)
		if _, err := os.Stat(name); err == nil {
			cpuPercent, err := proc.CPUPercent()
			if err != nil {
				log.Println(err)
			}
			if cpuPercent > 0 {
				name, err := proc.Name()
				if err != nil {
					log.Println(err)
				}
				if name != "prom-top" {
					p := ProcessCPU{Name: name, CPUUsedPercent: cpuPercent}
					p.Encode()
					reversed_freq[p.CPUUsedPercent] = append(reversed_freq[p.CPUUsedPercent], p)
				}
			}
		}
	}

	var numbers []float64
	for val := range reversed_freq {
		numbers = append(numbers, val)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(numbers)))
	if len(numbers) > 5 {
		numbers = numbers[:5]
	}
	for i, number := range numbers {
		for _, p := range reversed_freq[number] {
			p.Rank = i
			c.CPUByProcess = append(c.CPUByProcess, p)
		}
	}

	if c.Debug == true {
		ser, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(ser))
	}

}
