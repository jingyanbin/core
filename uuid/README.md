# uuid
唯一ID生成
PS F:\workspace\test\github.com\jingyanbin\core\uuid> go test -bench="." -parallel 500000 -count 1 -benchmem -cpu 2,4,8,16


goos: windows
goarch: amd64
pkg: github.com/jingyanbin/core/uuid
cpu: 11th Gen Intel(R) Core(TM) i7-11700 @ 2.50GHz
BenchmarkUUID-2                         87706474                13.37 ns/op            0 B/op          0 allocs/op
BenchmarkUUID-4                         86407395                13.37 ns/op            0 B/op          0 allocs/op
BenchmarkUUID-8                         92521915                13.48 ns/op            0 B/op          0 allocs/op
BenchmarkUUID-16                        89262468                13.80 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFast-2                     95632770                12.17 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFast-4                     95222979                12.01 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFast-8                     96041488                12.07 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFast-16                    92469157                12.14 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDParallel-2                 64140253                18.93 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDParallel-4                 34676668                29.14 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDParallel-8                 29850745                39.23 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDParallel-16                30052140                40.99 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFastParallel-2             99086610                12.08 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFastParallel-4             100000000               12.51 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFastParallel-8             86889470                12.56 ns/op            0 B/op          0 allocs/op
BenchmarkUUIDFastParallel-16            95468431                13.72 ns/op            0 B/op          0 allocs/op
