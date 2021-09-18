
# heap-ille
Unordered collection of records

# Operations supported
- Get : gets a Tuple by a key
- Put : puts a Tuple

# Open
- No support for meta pages which means heap file does not get read again
- No support for concurrency
- No support for handling large values which are bigger than the page size  


# Implementation details
A Tuple is a collection of fields. Currently, only string and unsigned-16 bit integer fields are supported. Heap file is organized into a collection of ```Slotted``` pages of fixed size. 

A ```Slotted``` page contains
- a page id
- slots, each slot contains a tuple offset into the page and the tuple size
- tuples

```Put(Tuple)``` into the heap file returns a ```TupleId``` which contains ```PageNo``` and ```SlotNo```. After the ```Put(Tuple)``` operation is done, the content of the page is written to the underlying file using standard IO.

Given this is a learning implementation, ```Tuple``` returns a ```Key``` which is the last field in the tuple. After the ```Put(Tuple)``` into the heap file, ```Key``` and ```TupleId``` are inserted into B+Tree.


# b-plus-tree
B+ tree disk implementation for storing key value pairs

# Operations supported
- Get : gets a KeyValuePair by a key
- Put : puts a KeyValuePair

# Open
- No support for meta pages which means index file does not get read again
- No support for concurrency
- No support for handling large values which are bigger than the page size  
- Optimization for Put, currently the entire page get marshalled and written to memory mapped file after each put

# Implementation details
B+Tree is modelled as a collection of pages which is represented in an abstraction called ```PageHierarchy```. 
By default, the size of each page is ```os.GetPageSize```.

B+Tree needs an instance of ```Options``` to initialize. As a part of Options one can provide the following -
- PageSize 
- FileName 
- PreAllocatedPagePoolSize 
- AllowedPageOccupancyPercentage  

As a part of initialization, B+Tree relies on ```PagePool``` to fetch the predefined number of pages from the underlying file, and this
file gets stored as a memory mapped file.

# Splitting a page
Technically, a page must contain at least ```t-1``` entries where ```t``` is the minimum degree and must be >=2. 
Also, a page can contain a maximum of ```2t-1``` entries, after which it should be split, but the current implementation splits a page if the ```size of the page``` goes beyond the ```allowed occupancy of a page```.

# References
- https://www.cs.utexas.edu/users/djimenez/utsa/cs3343/lecture17.html
- https://github.com/spy16/kiwi
