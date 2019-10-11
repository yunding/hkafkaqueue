package HKafkaQueue

import (
	"fmt"
	"sync"
	"testing"
)

var dataDir = "/tmp/hqueue"
var queueName = "test"

func TestNewHQueue(t *testing.T) {
	hqueue, err := NewHQueue(queueName, dataDir)
	if err != nil {
		t.Fatalf("create hqueue error %v", err)
	}
	fmt.Println(hqueue)
}

func TestQueueWrite(t *testing.T) {
	hqueue, err := NewHQueue(queueName, dataDir)
	var i = 0
	if err != nil {
		t.Fatalf("create hqueue error %v", err)
	}
	for {
		if hqueue.producerIndex.blockNum > 1 {
			break
		}
		_, err := hqueue.Offer([]byte("Bridgewater Associates AQR Capital Management Millennium Management Citadel Soros Fund Management Winton Capital Management D.E. Shaw& Co. enaissance Technologies LLC Two Sigma Paulson & Co."))
		if err != nil {
			t.Fatalf("write error :%v", err)
			break
		}
		i = i + 1
	}
	fmt.Println(hqueue.producerIndex.counter)
	fmt.Println(i)
	hqueue.Sync()

}

func TestQueueRead(t *testing.T) {
	hqueue, err := NewHQueue(queueName, dataDir, "p1")
	var i = 0
	if err != nil {
		t.Fatalf("create hqueue error: %v", err)
	}
	for {
		_, err := hqueue.Poll()
		if err != nil {
			break
		} else {
			//fmt.Println(String(bytes))
			i = i + 1
		}
	}
	fmt.Println(hqueue.consumerIndex.counter)
	fmt.Println(i)
	hqueue.Sync()
}

func BenchmarkWrite(t *testing.B) {
	hqueue, _ := NewHQueue(queueName, dataDir)
	var w sync.WaitGroup
	w.Add(1)
	write(hqueue, &w)
}

func BenchmarkRead(t *testing.B) {
	hqueue, _ := NewHQueue(queueName, dataDir, "p4")
	var w sync.WaitGroup
	w.Add(1)
	read(hqueue, &w)
}

func BenchmarkMultipleReadWrite(t *testing.B) {
	defer func() {
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	hqueue, _ := NewHQueue(queueName, dataDir)
	hqueue2, _ := NewHQueue(queueName, dataDir, "p5")
	var w sync.WaitGroup
	w.Add(2)
	go read(hqueue2, &w)
	go write(hqueue, &w)
	w.Wait()
}

func write(hqueue *HQueue, w *sync.WaitGroup) {
	var i = 0
	for {
		_, err := hqueue.Offer([]byte("Bridgewater Associates AQR Capital Management Millennium Management Citadel Soros Fund Management Winton Capital Management D.E. Shaw& Co. enaissance Technologies LLC Two Sigma Paulson & Co." +
			"Bridgewater Associates AQR Capital Management Millennium Management Citadel Soros Fund Management Winton Capital Management D.E. Shaw& Co. enaissance Technologies LLC Two Sigma Paulson & Co."))
		if err != nil {
			continue
		}
		i = i + 1
		if i%100000 == 0 {
			//fmt.Printf("write msg count:%d",i)
			//fmt.Println(hqueue)
			//time.Sleep(time.Second*2)
			w.Done()
			break
		}
	}
	//os.RemoveAll(dataDir+string(os.PathSeparator)+queueName)
}

func read(hqueue *HQueue, w *sync.WaitGroup) {
	var i = 0
	for {
		//fmt.Println("read......")
		b, err := hqueue.Poll()
		//release memory or will cause out of range
		if err != nil {
			if _, ok := err.(*ReadZeroError); ok {
				//fmt.Println("read Zero")
				//time.Sleep(time.Second)
				//if hqueue.consumerIndex.counter != 0 {
				//	fmt.Println(hqueue.consumerIndex.counter)
				//	fmt.Println(i)
				//}
				//time.Sleep(time.Second)
			} else {
				fmt.Errorf("%s", err)
			}
		} else {
			i = i + 1
			//fmt.Printf("read msg i:%d :%s\n",i,string(msg))
		}
		if i == 100000 {
			fmt.Println("break")
			w.Done()
			break
		}
		_ = b

	}

}
