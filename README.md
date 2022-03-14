### English translation
in this task should be written something like unix pipeline
```bash
grep 127.0.0.1 | awk '{print $2}' | sort | uniq -c | sort -nr
```

STDOUT of one program passes to STDIN of other program 

Everything is concurrent inside go app

Steps:
1. Write ExecutePipeline function which uses the assembly of worker functions that do something. 
2. Write some helper functions, which are mocks of some calculations, like count hashes

Chain of hash calculation:
1. ```SingleHash``` calc value ```crc32(data)+"~"+crc32(md5(data))``` ( string concatenation via ```~```), 
where data is what was received as parameter (basically it values from first function)
2. ```MultiHash``` calc value ```crc32(th+data))``` ( concatenation integer as string and string), where th=0..5 ( i.e. 6 hashes per incoming value ), 
then take concatenation of results in the same order (0..5), where data is something that comes as parameter (and returns as result from SingleHash)
3. ```CombineResults``` receive all results, sort it, and merge into one string via ```_```
4. ```crc32``` calc via ```DataSignerCrc32```
5. ```md5``` calc via ```DataSignerMd5```

#### Important notes

* ```DataSignerMd5``` can't be called in parallel, it makes calculations 10 ms. If it is called in parallel, there will be an overheat for 1 second
* ```DataSignerCrc32```, executing 1 sec
* Executing deadline 3 sec.

If run without concurrency - it takes ~57 second for 7 elements.

run via 
```bash
go test -v -race
```

Do not use ```select```, ```context```, ```time.Sleep```

Example results for sending 2 values, (commented in tests):
```
0 SingleHash data 0
0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
0 SingleHash crc32(md5(data)) 502633748
0 SingleHash crc32(data) 4108050209
0 SingleHash result 4108050209~502633748
4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555

1 SingleHash data 1
1 SingleHash md5(data) c4ca4238a0b923820dcc509a6f75849b
1 SingleHash crc32(md5(data)) 709660146
1 SingleHash crc32(data) 2212294583
1 SingleHash result 2212294583~709660146
2212294583~709660146 MultiHash: crc32(th+step1)) 0 495804419
2212294583~709660146 MultiHash: crc32(th+step1)) 1 2186797981
2212294583~709660146 MultiHash: crc32(th+step1)) 2 4182335870
2212294583~709660146 MultiHash: crc32(th+step1)) 3 1720967904
2212294583~709660146 MultiHash: crc32(th+step1)) 4 259286200
2212294583~709660146 MultiHash: crc32(th+step1)) 5 2427381542
2212294583~709660146 MultiHash result: 4958044192186797981418233587017209679042592862002427381542

CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
```

### Task description
В этом задании мы пишем аналог unix pipeline, что-то вроде:
```
grep 127.0.0.1 | awk '{print $2}' | sort | uniq -c | sort -nr
```

Когда STDOUT одной программы передаётся как STDIN в другую программу

Но в нашем случае эти роли выполняют каналы, которые мы передаём из одной функции в другую.

Само задание по сути состоит из двух частей
* Написание функции ExecutePipeline которая обеспечивает нам конвейерную обработку функций-воркеров, которые что-то делают.
* Написание нескольких функций, которые считают нам какую-то условную хеш-сумму от входных данных

Расчет хеш-суммы реализован следующей цепочкой:
* SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~), где data - то что пришло на вход (по сути - числа из первой функции)
* MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)
* CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/), объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
* crc32 считается через функцию DataSignerCrc32
* md5 считается через DataSignerMd5

В чем подвох:
* DataSignerMd5 может одновременно вызываться только 1 раз, считается 10 мс. Если одновременно запустится несколько - будет перегрев на 1 сек
* DataSignerCrc32, считается 1 сек
* На все расчеты у нас 3 сек.
* Если делать в лоб, линейно - для 7 элементов это займёт почти 57 секунд, следовательно надо это как-то распараллелить

Результаты, которые выводятся если отправить 2 значения (закомментировано в тесте):

```
0 SingleHash data 0
0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
0 SingleHash crc32(md5(data)) 502633748
0 SingleHash crc32(data) 4108050209
0 SingleHash result 4108050209~502633748
4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555

1 SingleHash data 1
1 SingleHash md5(data) c4ca4238a0b923820dcc509a6f75849b
1 SingleHash crc32(md5(data)) 709660146
1 SingleHash crc32(data) 2212294583
1 SingleHash result 2212294583~709660146
2212294583~709660146 MultiHash: crc32(th+step1)) 0 495804419
2212294583~709660146 MultiHash: crc32(th+step1)) 1 2186797981
2212294583~709660146 MultiHash: crc32(th+step1)) 2 4182335870
2212294583~709660146 MultiHash: crc32(th+step1)) 3 1720967904
2212294583~709660146 MultiHash: crc32(th+step1)) 4 259286200
2212294583~709660146 MultiHash: crc32(th+step1)) 5 2427381542
2212294583~709660146 MultiHash result: 4958044192186797981418233587017209679042592862002427381542

CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
```

Код писать в signer.go. В этот файл не надо добавлять ничего из common.go, он уже будет на сервере.

