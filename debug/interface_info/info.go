package interfaceinfo

// https://blog.csdn.net/weiguang102/article/details/129299867
// interface 是一个特殊结构体, 包含类型和值, 当interface进行比较的时候, 需要满足两边都是interface, 不满足的时候会将非interface的对象转成interface的对象, 然后需要类型相等, 值也相等
// 不可比较的类型有切片, map, 函数
