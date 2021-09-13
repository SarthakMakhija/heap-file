
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

# References
- https://www.cs.utexas.edu/users/djimenez/utsa/cs3343/lecture17.html
- https://github.com/spy16/kiwi
