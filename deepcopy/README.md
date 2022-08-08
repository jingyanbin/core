//小对象测试
goos: windows

goarch: amd64     

cpu: Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz

BenchmarkCopyByJson

BenchmarkCopyByJson-16            300951              3624 ns/op            1104 B/op          9 allocs/op

BenchmarkCopyByJson-16            333441              3637 ns/op            1104 B/op          9 allocs/op

BenchmarkCopyByJson-16            333435              3619 ns/op            1104 B/op          9 allocs/op

BenchmarkCopyByMsgPack

BenchmarkCopyByMsgPack-16         214356              5489 ns/op             538 B/op         16 allocs/op

BenchmarkCopyByMsgPack-16         214350              5485 ns/op             538 B/op         16 allocs/op

BenchmarkCopyByMsgPack-16         222127              5468 ns/op             538 B/op         16 allocs/op

BenchmarkCopyByGob

BenchmarkCopyByGob-16              18813             63765 ns/op           12893 B/op        369 allocs/op

BenchmarkCopyByGob-16              19053             63382 ns/op           12894 B/op        369 allocs/op

BenchmarkCopyByGob-16              18784             63385 ns/op           12894 B/op        369 allocs/op

BenchmarkCopy

BenchmarkCopy-16                 2269152               524.7 ns/op            88 B/op          2 allocs/op

BenchmarkCopy-16                 2281935               523.5 ns/op            88 B/op          2 allocs/op

BenchmarkCopy-16                 2282097               526.6 ns/op            88 B/op          2 allocs/op

BenchmarkCopyByGoGo

BenchmarkCopyByGoGo-16           4820786               247.6 ns/op            96 B/op          2 allocs/op

BenchmarkCopyByGoGo-16           4859839               247.7 ns/op            96 B/op          2 allocs/op

BenchmarkCopyByGoGo-16           4835924               251.9 ns/op            96 B/op          2 allocs/op


//大对象faster测试 使用gogofaster编译出的pb对象

goos: windows

goarch: amd64

cpu: Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz

BenchmarkCopyByJson

BenchmarkCopyByJson-16               999           1168981 ns/op          181627 B/op       2162 allocs/op

BenchmarkCopyByJson-16              1016           1184009 ns/op          181691 B/op       2162 allocs/op

BenchmarkCopyByJson-16              1005           1174328 ns/op          181759 B/op       2162 allocs/op

BenchmarkCopyByMsgPack

BenchmarkCopyByMsgPack-16            792           1491351 ns/op          203560 B/op       1528 allocs/op

BenchmarkCopyByMsgPack-16            794           1499519 ns/op          203561 B/op       1528 allocs/op

BenchmarkCopyByMsgPack-16            794           1484412 ns/op          203560 B/op       1528 allocs/op

BenchmarkCopyByGob

BenchmarkCopyByGob-16               5001            239639 ns/op           82653 B/op        963 allocs/op

BenchmarkCopyByGob-16               5002            240326 ns/op           82655 B/op        963 allocs/op

BenchmarkCopyByGob-16               4898            237643 ns/op           82654 B/op        963 allocs/op

BenchmarkCopy

BenchmarkCopy-16                    6313            180076 ns/op           35416 B/op        357 allocs/op

BenchmarkCopy-16                    6856            178741 ns/op           35416 B/op        357 allocs/op

BenchmarkCopy-16                    7060            180655 ns/op           35416 B/op        357 allocs/op

BenchmarkCopyByGoGo

BenchmarkCopyByGoGo-16             10000            114563 ns/op           72105 B/op        817 allocs/op

BenchmarkCopyByGoGo-16             10000            114563 ns/op           72104 B/op        817 allocs/op

BenchmarkCopyByGoGo-16             10000            115263 ns/op           72104 B/op        817 allocs/op

PASS

ok      command-line-arguments  19.847s

//大对象测试 使用gogofast编译出的pb对象

goos: windows

goarch: amd64

cpu: Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz

BenchmarkCopyByJson            

BenchmarkCopyByJson-16               974           1214192 ns/op          190305 B/op       2625 allocs/op

BenchmarkCopyByJson-16               975           1218107 ns/op          190300 B/op       2625 allocs/op

BenchmarkCopyByJson-16               967           1204373 ns/op          190477 B/op       2625 allocs/op

BenchmarkCopyByMsgPack           

BenchmarkCopyByMsgPack-16            730           1637848 ns/op          220040 B/op       2345 allocs/op

BenchmarkCopyByMsgPack-16            722           1626772 ns/op          220040 B/op       2345 allocs/op

BenchmarkCopyByMsgPack-16            726           1641516 ns/op          220040 B/op       2345 allocs/op

BenchmarkCopyByGob

BenchmarkCopyByGob-16               3807            312971 ns/op           94880 B/op       1510 allocs/op

BenchmarkCopyByGob-16               3870            312739 ns/op           94877 B/op       1510 allocs/op

BenchmarkCopyByGob-16               3807            316471 ns/op           94882 B/op       1510 allocs/op

BenchmarkCopy

BenchmarkCopy-16                    4208            271862 ns/op           43768 B/op        820 allocs/op

BenchmarkCopy-16                    4285            266926 ns/op           43768 B/op        820 allocs/op

BenchmarkCopy-16                    4363            264360 ns/op           43768 B/op        820 allocs/op

BenchmarkCopyByGoGo

BenchmarkCopyByGoGo-16              8574            129770 ns/op           80456 B/op       1280 allocs/op

BenchmarkCopyByGoGo-16              8574            129963 ns/op           80456 B/op       1280 allocs/op

BenchmarkCopyByGoGo-16              9234            129480 ns/op           80456 B/op       1280 allocs/op

PASS

ok      command-line-arguments  19.131s

总结: 
大对象拷贝: gogo > reflect > gob > json > msgpack

小对象拷贝: gogo > reflect > json > msgpack > gob

1.对pb对象进行拷贝使用gogo序列化与反序列化的速度是最快的,但是不能用与非pb对象的拷贝(faster生成的协议序列化比fast生成的协议更快)

2.其他对象拷贝使用反射实现的Copy是最快的 CopyAll 可以拷贝私有成员变量, 其他放置均不支持私有成员变量的拷贝


