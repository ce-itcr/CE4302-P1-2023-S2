package cache_controller

type Cache struct{
	Data	[4]int
	Address	[4]int
	Status	[4]string
}

func NewCache() *Cache{
	return &Cache{
		Data: [4]int{0, 0, 0, 0},
		Address: [4]int{0, 0, 0, 0},
		Status: [4]string{"", "", "", ""},
	}
}

func (cache *Cache) setData(pos int, res int){
	cache.Data[pos] = res
}

func (cache *Cache) setAddress(pos int, res int){
	cache.Address[pos] = res
}

func (cache *Cache) setState(pos int, res string){
	cache.Status[pos] = res
}

func (cache *Cache) getData(pos int) int{
	return cache.Data[pos]
}

func (cache *Cache) getAddress(pos int) int{
	return cache.Address[pos]
}

func (cache *Cache) getState(pos int) string{
	return cache.Status[pos]
}

