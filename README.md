# logstrata

HLL 
Bloom Filter
Skip List → Memtable → WAL (reuse W1)
Data Blocks → SSTable → Index Block
Bloom Filter → SSTable (attach per SSTable)
Memtable + SSTable → Block Cache
Block Cache + SSTable → Compaction
Compaction → Write Stall
Everything → Concurrent read safety