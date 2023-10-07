package CacheController

type Cache struct{
	data	[4]int
	address	[4]int
	status	[4]string
}

func NewCache() *Cache{
	return &Cache{
		data: [4]int{0, 0, 0, 0},
		address: [4]int{-1, -1, -1, -1},
		status: [4]string{"I", "I", "I", "I"},
	}
}

func (cache *Cache) SetData(pos int, res int){
	cache.data[pos] = res
}

func (cache *Cache) SetAddress(pos int, res int){
	cache.address[pos] = res
}

func (cache *Cache) SetState(pos int, res string){
	cache.status[pos] = res
}

func (cache *Cache) GetData(pos int) int{
	return cache.data[pos]
}

func (cache *Cache) GetAddress(pos int) int{
	return cache.address[pos]
}

func (cache *Cache) GetState(pos int) string{
	return cache.status[pos]
}
