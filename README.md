
# b-plus-tree
B+ tree disk implementation for storing key value pairs

# Operations supported
- Get : gets a KeyValuePair by a key
- Put : puts a KeyValuePair

# Open
- No support for meta pages which means index file does not get read again
- No support for concurrency

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
