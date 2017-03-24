### Zipf's law in go

Zipf's law states:

> A pattern of distribution in certain data sets, notably words in a linguistic corpus, by which the frequency of an item is inversely proportional to its ranking by frequency.

The programs here attempt to show Zipf's law in action by parsing books from guttenberg.org and counting the frequency of words. There are two versions: one uses channels and one does not. The exercise being to see what effect using channels has on time and effort to process the data.

Both programs parse the files in the `books` directory, loop through the lines, split (and roughly tokenize) the words, create a map of the words to their frequency, sort the map and then generate a bar chart for a graphical view of Zipf's law in action.

The program using channels (`channels.go` if you were in any doubt) uses multiple go routines and channels at certain points.

##### Usage

Open a couple of terminal windows and run
```
go run channels.go
```
in one and
```
go run no_channels.go
```
in the other. The output will be the time taken and the word count.

Notice the programs don't exit after completion, this is due to the profiler running.

To view the heap profile in another terminal use the command
```
go tool pprof http://localhost:8000/debug/pprof/heap
```
to view the channels version and change the port to `8080` to view the non-channels version.

Once in the interactive `(pprof)` terminal enter `top` to get a view of the top memory usage of the heap. The disclaimer here is that I know nothing about profiling Go programs, it's quite interesting to take a look though.

Hit `ctrl + c` to quit.

You can drop more .txt files into the books folder and see how each program compares with a larger and larger corpus.

MIT license.