# Use LRU Least Recently Used implementation for the message tracker 

- Status: accepted
- Date: 2024-06-11
- Tags: doc

## Context and Problem Statement

The Requirements for the Message Tracker are:
1. Efficient Addition: New messages should be added to the front of the cache.
2. Efficient Search: We need to quickly check if a message exists in the cache.
3. Efficient Update: If a message is repeated, it should be moved to the front of the cache. 
4. Efficient Deletion: If the cache exceeds its capacity, the oldest message should be removed.

## Considered Options

- [LRU] - Last Recently Used

## Decision Outcome

Chosen option: "LRU". An LRU cache is a data structure that evicts the least recently used items first when it reaches 
its capacity. This behavior is ideal for a message tracker in a P2P system where:

  - Frequently accessed messages need to stay in the cache.
  - Less frequently accessed messages can be removed when new messages arrive and the cache reaches its capacity.

## Data Structures Used

  - Linked List: To maintain the order of messages and support efficient updates, deletions, and insertions.
  - Hash Map: To allow for O(1) average-time complexity for search operations.

## How The LRU Cache Works

  - Linked List: This doubly linked list keeps track of the message order. The head of the list points to the most 
recently used message, while the tail points to the least recently used message.
  - Hash Map: This map stores pointers (references) to the nodes in the linked list, allowing for O(1) time complexity 
for search operations.

## Go Implementation Details

  - In Go, this can be efficiently implemented using the "container/list" package for the doubly linked list and a 
built-in map for the hash map.

## Positive Consequences

  - Search (O(1)): The hash map allows you to quickly find if a message exists in the cache.
  - Update (O(1)): When a message is repeated, you can use the hash map to find its node in the linked list and move it 
to the front.
  - Delete (O(1)): If the cache is full, you can remove the node at the tail of the linked list and update the hash map.
  - Add    (O(1)): Adding a new message involves inserting a node at the front of the linked list and updating the hash map.

