package redis

type Option func(*RedisComponent)

func defaultOption() RedisComponent {
	return RedisComponent{}
}
