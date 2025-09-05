# Filters

- Bloom Filter
- Counting Bloom Filter
- Cuckoo Filter
- Quotient Filter
- XOR Filter
- Count-Min Sketch
- HyperLogLog
How to Choose: A Simple Decision Guide
Here is a step-by-step way to know which filter to use.

Step 1: What problem are you solving?
First, determine the exact question you need to answer.

A) "Is this item in my set?" (Membership)

Example: "Has this user seen this article before?"

If yes, continue to Step 2. This is the most common use case.

B) "How many times have I seen this item?" (Frequency)

Example: "How many times has this IP address tried to log in?"

If yes, your answer is the Count-Min Sketch.

C) "How many unique items are in my set?" (Cardinality)

Example: "How many unique people visited my website today?"

If yes, your answer is the HyperLogLog.

Step 2: Does your data change? (Static vs. Dynamic)
This is the most important question for choosing a membership filter.

A) My data is STATIC. (You write the set once and don't change it.)

Example: A dictionary of all English words, a list of all known malware signatures.

Your best choice is an XOR Filter. It's the fastest and most memory-efficient for static data. A standard Bloom Filter is also a simple and classic choice.

B) My data is DYNAMIC. (You need to add and remove items.)

Example: A list of users currently logged into a chat room, items in a shopping cart.

If your data is dynamic, continue to Step 3.

Step 3: What is your main priority?
If you need to add and remove items, your choice comes down to a final trade-off.

A) "I need the best all-around performance and space savings."

Your best choice is a Cuckoo Filter. It has excellent space efficiency, the fastest lookups, and is the modern default for dynamic sets. You just need to handle the rare chance an insert might fail when it's very full.

B) "I need simplicity and absolutely guaranteed insertions."

Your choice is a Counting Bloom Filter. It's a simple extension of the Bloom filter, but be prepared for it to use 3-4x more memory than a Cuckoo filter.

C) "I need to squeeze every drop of CPU cache performance."

Your choice is a Quotient Filter. This is a more advanced option but can be faster than a Cuckoo filter in specific, cache-sensitive applications.

The 30-Second Rule of Thumb
When in doubt, use this simple summary:

Counting unique items? ➡️ HyperLogLog.

Checking for items that are never deleted? ➡️ XOR Filter.

Checking for items that can be added and deleted? ➡️ Cuckoo Filter.