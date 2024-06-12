# Use sync.RWMutex in Message Tracker

- Status: accepted
- Date: 2024-06-12
- Tags: doc

## Context and Problem Statement

In our P2P messaging system usually a peer receive messages concurrently from many other peers to which is connected.
For that reason we should make sure that only one goroutine access the shared data ( linked list with messages and the map) 
at a given time.

Using sync.RWMutex in Go is a good choice for managing concurrent access to a shared resource, particularly in scenarios 
where read operations significantly outnumber write operations. Here are the key reasons why sync.RWMutex is suitable and beneficial`

- Multiple Readers: sync.RWMutex allows multiple readers to hold the lock simultaneously, which significantly improves 
performance in read scenarios.This is in contrast to a standard sync.Mutex, which allows only one goroutine to access 
the critical section at a time, even for reads.
- Single Writer: Write operations are exclusive; only one writer can hold the lock, ensuring data integrity during modifications.