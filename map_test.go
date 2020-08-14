package concurrent

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func randString(n int) string {
	const tb = "qazxswedcvfrtgbnhyumkiolp1234567890QAZXSWEDCVFRTGBNHYUJMKIOLP"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		ret[i] = tb[rand.Intn(len(tb))]
	}
	return string(ret)
}

func Test_Map(t *testing.T) {
	m := NewMapWithMod(1)
	n := 1000 * 10000
	a := make([]string, n)
	for i := range a {
		a[i] = randString(10)
	}
	for i := range a {
		m.Set(ID(i), a[i])
	}
	for i := range a {
		if v, b := m.Get(ID(i)); !b {
			t.Fatalf("not found: %v", i)
		} else if v != a[i] {
			t.Fatal("incorrect values")
		}
	}

	//m.ForEach(func(v node) {
	//	t.Log(v.key, v.value)
	//})
	indexList := rand.Perm(n)
	t.Log("start test")
	st := time.Now()
	for _, i := range indexList {
		v, _ := m.Get(ID(i))
		//t.Log(i, v)
		if v != a[i] {
			t.Fatal("incorrect value: ", v)
		}
	}
	t.Log("consumption: ", time.Now().Sub(st))
	t.Log("op/s:", float64(n)/time.Now().Sub(st).Seconds())
}

func TestBalance(t *testing.T) {
	m := NewMapWithMod(1)
	n := 10
	a := make([]string, n)
	for i := range a {
		a[i] = randString(10)
	}
	for i := range a {
		m.Set(ID(i), a[i])
	}
	for i := range a {
		if v, b := m.Get(ID(i)); !b {
			t.Fatalf("not found: %v", i)
		} else if v != a[i] {
			t.Fatal("incorrect values")
		}
	}
	for i := range m.root {
		m.forEach(m.root[i], func(v node) {
			t.Log(v.key, v.value, max(v.leftDepth, v.rightDepth))
		})
	}
}

func TestPerformance(t *testing.T) {
	sm := sync.Map{}
	m := NewMap()
	n := 100 * 10000
	a := make([]string, n)
	for i := range a {
		a[i] = randString(10)
	}
	indexList := rand.Perm(len(a))
	for _, i := range indexList {
		m.Set(ID(i), a[i])
		sm.Store(i, a[i])
	}

	concurrentNum := 100
	for i := 0; i < concurrentNum; i++ {
		go func() {
			c := time.Tick(time.Millisecond)
			for range c {
				idx := rand.Intn(len(a))
				sm.Store(idx+len(a), a[idx])
			}
		}()
	}
	for i := 0; i < concurrentNum; i++ {
		go func() {
			c := time.Tick(time.Millisecond)
			for range c {
				idx := rand.Intn(len(a))
				m.Set(ID(idx+len(a)), a[idx])
			}
		}()
	}

	indexList = rand.Perm(n)
	t.Log("start test concurrent.Map")
	st := time.Now()
	for _, i := range indexList {
		v, _ := m.Get(ID(i))
		//t.Log(i, v)
		if v != a[i] {
			t.Fatal("incorrect value: ", v)
		}
	}
	t.Log("consumption: ", time.Now().Sub(st))
	t.Log("op/s:", float64(n)/time.Now().Sub(st).Seconds())

	t.Log("start test sync.Map")
	st = time.Now()
	for _, i := range indexList {
		v, _ := sm.Load(i)
		if v.(string) != a[i] {
			t.Fatal("incorrect value: ", v)
		}
	}
	t.Log("consumption: ", time.Now().Sub(st))
	t.Log("op/s:", float64(n)/time.Now().Sub(st).Seconds())
}
