### version1

- Environment

MacBook Pro (Mid 2015), OSX 10.14.1, 2.5 GHz Intel Core i7, 16 GB 1600 MHz DDR3

- Sample

```dash
$ du -sh generateLog/test.log
  5.9G	generateLog/test.log
$ wc -l generateLog/test.log
  243706857 generateLog/test.log
```

- duration

rehash file: 30 min

![ver1_1](../image/ver1_1.png)

calTopN: 2min

![ver1_2](../image/ver1_2.png)

- result

![ver1_result](../image/ver1_result.png)

- pprof

heap(coutTable):
![ver1_heap](../image/ver1_heap.png)

![ver1_source](../image/ver1_source.png)

rehash profile(file write):
![ver1_rehash_flamegraph](../image/ver1_rehash_flamegraph.png)

![ver1_rehash_top](../image/ver1_rehash_top.png)

