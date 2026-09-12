---
title: How to implement a DBMS?
date: 2023-08-30
---

> **Note:** This article was originally published on my Medium while I was taking database classes at the University of São Paulo

DBMSs (Database Management Systems) are software created to facilitate the storage and retrieval of data in a computer system

Usually, when we develop an application that accesses a DBMS, we don't think much about what happens under the hood when we interact with the stored data

The main goal of this article is to present the core ideas behind implementing a DBMS, along with a bonus that gives a simplified explanation of how SQLite was implemented

## Data comes and goes But what does it represent?

Before understanding what happens during the execution of a DBMS, we need to understand what data is

Data can be understood conceptually or physically

A **conceptual data** is an isolated set of qualitative or quantitative values For example, in a shopping list, the expressions _3_ and _lettuce_ are data The combination of this data can generate information, such as "buy 3 lettuce leaves"

When we talk about **physical data**, on the other hand, we are referring to any sequence of bytes that represents a conceptual data For example, using UTF-8 binary encoding, we can represent the expression _lettuce_ as the hexadecimal sequence _6C 65 74 74 75 63 65_

In the context of DBMSs, data is represented physically, since it is stored in **storage devices** that only understand binary values

## A bit about storage devices

Physical data is stored in storage devices, which are hardware specialized in storing, searching for, and retrieving bytes

Examples of these devices are:

- _HD_ (Hard disk)
- _RAM_ (Random Access Memory)
- _SSD_ (Solid-state drive)
- _DVD_ (Digital versatile disk)

Each device has a unique way of managing these bytes, whether through electrical charge (SSD), magnetism (HD), light sensitivity (DVD), or any other physical medium

These devices can be **volatile** or **non-volatile** Volatile devices, unlike non-volatile ones, need power to keep the stored data This means that if, for example, we turn off the server, all its data will be lost RAM is a volatile device, while HD, SSD, and DVD are non-volatile

## They're all the same, yet different…

From the very definition given earlier about a DBMS, we can conclude that all DBMSs have in common the storage of physical data in storage devices

However, there are big differences between each of them, ranging from the data model (relational, document-oriented, …) to the type of storage device used (volatile or non-volatile)

The first step to learning how to build a DBMS is to understand that, as different as the purposes of _MySQL_, _MongoDB_, and _Redis_ may be, all of them are based on storing physical data

## And where does the operating system come into this story?

Although it might seem like DBMSs access storage devices directly, it is actually the operating system that does so

This happens because of the extent and complexity of interacting with this hardware Imagine you wanted to implement a DBMS and, to do so, needed to understand how to communicate with a specific device through the memory bus Hard, right?

The operating system is a facilitator It provides ways for us to communicate indirectly with these devices through system calls (syscalls)

The way this communication happens is through **file systems**, which will be explained in the next topic

## Systems of what!?

File systems are an abstraction offered by the operating system to communicate with storage devices without having to access the hardware directly

Each file represents a set of bytes stored on a device (such as an HD) Files have identifiers (usually names) and can generally be searched through a directory tree

When we use a programming language to modify a file, under the hood we are accessing data on a storage device

The second step to learning how to build a DBMS is understanding that it uses files (usually binary) as the fundamental basis for data management

However, one caveat should be made When we talk about DBMSs that use memory (RAM) for data storage (as Redis does), we don't necessarily need to use files We can simply create data structures such as arrays or linked lists for temporary storage

For the sake of simplicity, let's assume until the end of this article that the only abstraction used by DBMSs is the file system

## A bit about files…

As mentioned earlier, a file is nothing more than a set of bytes located on some storage device

Accessing a file is usually relatively slow This is because we need to make a system call to the operating system, which must retrieve the data contained in that file from the hardware For example, for an HD, there is a read-write arm that needs to move to find the file's bytes, which is a slow operation

It's expected that byte organization should be optimized so we can insert, remove, and search for specific data with the fewest possible file accesses

There are many ways to organize data in a file, and this organization is one of the points that can differ from one DBMS to another

Examples of byte organization in a file include B+ trees and hash tables Although we usually think of these data structures only in memory (RAM), we can also represent and serialize them in a file

## I trust the storage engine

DBMSs can be split into different modules, where each module is responsible for a specific functionality We can have a module that performs parsing on a query language (such as SQL), another that optimizes queries, another that interacts with files, and so on

The module that interacts with files and defines how data is organized in those files is called the **storage engine** (also known as _database engine_)

Examples of storage engines include:

- [InnoDB](https://dev.mysql.com/doc/refman/en/innodb-storage-engine.html) (used by default in MySQL)
- [Bitcask](https://github.com/basho/bitcask)
- [TileDB](https://github.com/TileDB-Inc/TileDB)
- [PumpkinDB](https://github.com/PumpkinDB/PumpkinDB)

Of all the modules a DBMS can provide, the storage engine is one of the most important, since without it we would have no way to store and retrieve data on devices

## Performing operations on files

Understanding how we can, through code, access and manipulate bytes in a file is essential for analyzing how storage engines are implemented

As mentioned earlier, files are abstractions provided by the operating system for manipulating bytes on a storage device without having to access it directly

The operating system provides us with several operations that can be performed on files Examples of these operations include:

- _Creation_: Creates a new file
- _Removal_: Removes a file
- _Read_: Reads a fixed amount of bytes
- _Write_: Writes a sequence of bytes
- _Seek_: Searches for a specific byte

For teaching purposes, we can understand a file as an array of bytes Each index of this array is called an **offset**, and holds a single associated byte

The file also has what we call a **file cursor** (or file pointer), which is a fundamental concept for file manipulation It represents the current read or write position within a file

For example, if we have a file with 100 bytes, with the current position at offset 0 (the first byte of the file), and we read 10 bytes, then the cursor moves to offset 10 (the eleventh byte of the file) If, from that position, we write another 10 bytes, the cursor then moves to offset 20

We can also seek a specific offset, which means we will be changing the cursor's position This operation is usually costly, since the operating system needs to locate where that offset is on the storage device

Programming languages are responsible for providing an interface that communicates with the operating system for file manipulation For example, the C language provides the [_fread_](https://learn.microsoft.com/en-us/cpp/c-runtime-library/reference/fread), [_fwrite_](https://learn.microsoft.com/en-us/cpp/c-runtime-library/reference/fwrite), and [_fseek_](https://learn.microsoft.com/en-us/cpp/c-runtime-library/reference/fseek) functions to read, write, and seek bytes in a file, respectively The Java language, through the [_RandomAccessFile_](https://docs.oracle.com/javase/7/docs/api//java/io/RandomAccessFile.html) class, provides the _read_, _write_, and _seek_ methods

## Unveiling records and blocks

So far, we've only discussed sequences of bytes in files, but we haven't covered how these bytes actually represent conceptual data

Two concepts that can be found in storage engine implementations are records and blocks

A **record** within a binary file is a sequence of bytes that represents different fields with different types and values The concept of a record is normally used in relational DBMSs

An example of a record with the fields _name_ and _favorite food_ is the byte sequence (in hexadecimal) _74 68 61 69 73 74 6F 6D 61 74 65_, which can be deserialized into the values _thais_ and _tomato_

There are several techniques for serializing and deserializing bytes into conceptual data One of them is [Protocol Buffers](https://github.com/protocolbuffers/protobuf), a format used by Google for communication between different services, but which can also be used for storage in files Another format, specific to the Go programming language, is provided in the [encoding/gob](https://pkg.go.dev/encoding/gob) package

A **block** within a binary file is a fixed-size sequence of bytes that stores one or more data records Reading binary files can be done in blocks instead of individual bytes, increasing the number of records retrieved in a single read

In some storage engines, a block represents a node (or page) of a B+ tree, as happens in SQLite

Note that records and blocks are just techniques used for managing data in binary files, and are not necessarily mandatory There are several storage engine implementations, and each one can follow a different approach

## A bit about data models…

Now that we understand how data is managed at the lowest level by a DBMS, we can abstract the discussion a bit

Data models are theoretical models created for working with conceptual data following some defined system

Examples of data models include:

- Relational
- Document-oriented
- Graph-oriented

When implementing a DBMS, one of the first decisions that must be made is choosing which model or models the DBMS will follow

DBMSs like PostgreSQL and MySQL follow the relational model, where data is represented by rows and columns in a table, and these tables can relate to each other

Under the hood, we are still dealing with bytes in files, but what differentiates one data model from another is how these bytes are represented

## Other DBMS modules

Earlier, we talked about the storage engine, which is a DBMS module responsible for managing data on a storage device

There are several other modules a DBMS can implement

Examples of other modules include:

- _Compiler_: When the DBMS uses a language such as SQL to perform queries, a compiler is essential for interpreting these queries
- _Optimizer_: Responsible for optimizing any query made to the DBMS, which may include creating equivalent queries and choosing algorithms to access the data
- _Interface_: Some DBMSs provide an interface (command line or any other means) to interact with different functionalities
- _Driver_: A driver allows a programming language to interact with different DBMS functionalities

A DBMS usually contains several other modules to ensure its operation An in-depth explanation of each module is beyond the scope of this article

## Bonus: How was SQLite implemented?

SQLite is a relational DBMS created with the C programming language All data is stored in a single binary file structured with multiple B+ trees, each representing a table or index

The first bytes of the file contain metadata about the database, which assist in different operations SQLite performs on the file

Fortunately, SQLite's source code is public domain and can be seen at [github.com/sqlite/sqlite](https://github.com/sqlite/sqlite) If you want to dig deeper into the file's binary format, the official website also provides a detailed explanation that can be accessed [at this link](https://www.sqlite.org/fileformat.html)
