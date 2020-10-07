# ParallelDots_Assignment

## Steps to run

1. `cd src/github.com/ParallelDots`
2. Run `go build ./...` . This will create executables for all the files in their respective directory.
3. Run `server` and `client` executables in seperate terminals.
4. For server use `server -cache` to enable caching as by default it is disabled and for client use `client -requests=value` to specify number of concurrent requests to send.

### Approach

I thought of a hybrid approach where cache can live as a persistent layer as well as on RAM. To make it persistent I first thought of saving the data on a file
after a certain period of time which would require event handling if the server is still serving up requests. So, instead I implemented it in such a way that when
the server is shutdown the whole data of the cache is saved into a file, and when next time the server is spawned it can load that file and use the data stored to
serve requests.

I used maps to store data as they provide constant lookup time. Also I used gob package to encode and decode data as the data will reside in RAM as well as hard disk.

There are many improvements that can be made like adding a timer to cache so that the unwanted/least used data can be removed from cache, serving multiple servers at once(that
will require more synchronization) , etc.

### Benchmarking Results

1. Without caching for 100 requests each request takes around 30.75 - 30.95 seconds.
2. With caching enabled for 80 requests the time is 54-70 milli seconds and for the remaining 20 requests it is
    between 30.077 - 30.089 seconds.

We see a significant improvement here from ~30 seconds to 54 milli seconds for cached requests.