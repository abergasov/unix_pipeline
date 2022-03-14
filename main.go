package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// multiHash struct for collect data in same order
type multiHash struct {
	iteration int
	hash      string
}

const bufferSize = 100

// ExecutePipeline target function for pipeline realization
func ExecutePipeline(jobs ...job) {
	start := time.Now()
	in := make(chan interface{}, bufferSize)
	var wg sync.WaitGroup
	for _, j := range jobs {
		tmpOut := make(chan interface{}, bufferSize)
		wg.Add(1)
		go func(tmpIn, tmpOut chan interface{}, jb job, wg *sync.WaitGroup) {
			defer close(tmpOut)
			defer wg.Done()
			jb(tmpIn, tmpOut)
		}(in, tmpOut, j, &wg)
		in = tmpOut
	}

	wg.Wait()
	end := time.Since(start)
	println(fmt.Sprintf("pipeline finished at %d second", end/time.Second))
	println(3 * time.Second)

}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for i := range in {
		str := strconv.Itoa(i.(int))
		md5 := DataSignerMd5(str)
		wg.Add(1)
		go func() {
			defer wg.Done()
			wrapSingleHash(str, md5, out)
		}()
	}
	wg.Wait()
}

func wrapSingleHash(str, md5 string, res chan interface{}) {
	pt1 := ""
	pt2 := ""
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		pt1 = DataSignerCrc32(str)
		wg.Done()
	}()
	go func() {
		pt2 = DataSignerCrc32(md5)
		wg.Done()
	}()
	wg.Wait()
	res <- pt1 + "~" + pt2
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for i := range in {
		str := i.(string)
		wg.Add(1)
		go func() {
			defer wg.Done()
			calculateSingleMultiHash(str, out)
		}()
	}
	wg.Wait()
}

func calculateSingleMultiHash(calc string, res chan interface{}) {
	ch := make(chan multiHash, 6)
	for j := 0; j < 6; j++ {
		go func(s string, c chan multiHash, counter int) {
			c <- multiHash{
				iteration: counter,
				hash:      DataSignerCrc32(strconv.Itoa(counter) + s),
			}
		}(calc, ch, j)
	}
	mHash := make([]string, 6, 6)
	iteration := 0
	for i := range ch {
		mHash[i.iteration] = i.hash
		iteration += 1
		if iteration == 6 {
			break
		}
	}
	res <- strings.Join(mHash, "")
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for i := range in {
		str := i.(string)
		res = append(res, str)
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	ExecutePipeline([]job{
		func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
	}...)
}
