concurrent.Map 
======================

#### Test Example

```go

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

```

#### Output

###### concurrentNum = 100
```shell script
map_test.go:113: start test concurrent.Map
map_test.go:122: consumption:  1.35818294s
map_test.go:123: op/s: 736257.828255837
map_test.go:125: start test sync.Map
map_test.go:133: consumption:  513.661615ms
map_test.go:134: op/s: 1.9467192670212614e+06
```

###### concurrentNum = 1000
```shell script
map_test.go:113: start test concurrent.Map
map_test.go:122: consumption:  1.765569881s
map_test.go:123: op/s: 566376.4732215261
map_test.go:125: start test sync.Map
map_test.go:133: consumption:  4.923516887s
map_test.go:134: op/s: 203105.45651384036
```
