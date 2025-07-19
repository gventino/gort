# GoRt - Go Rainbow Table
GoRt is a command-line tool written in Go to demonstrate the concept of a pre-computed hash rainbow table for cracking SHA-512 hashes. This project is intended for educational purposes to illustrate how a time-memory tradeoff attack works.

***Disclaimer***: This tool is for educational use only. Do not use it for malicious activities.

## What is GoRt?
GoRt implements a pre-computation attack. It works in two main phases:

1. Generation Phase: It generates a large number of potential passwords, computes the SHA-512 hash for each one, and stores these (hash, password) pairs in a lookup table. This table is then saved to a binary file.

2. Lookup Phase: Given a SHA-512 hash, the tool performs a quick search in the pre-computed table. If the hash is found, the corresponding original password is returned.

## Features
1. Password Generation: Create a list of random passwords based on a defined charset and length.

2. Table Generation: Build a lookup table from a password list, hashing each entry with SHA-512.

3. Hash Cracking: Look up a given SHA-512 hash in the generated table to find the original password.

4. Testing: Test the effectiveness of the table against randomly generated passwords.

## How to Use
### Prerequisites
Go (version 1.18 or newer) must be installed on your system.

### Running the Application
Clone the repository:

```
git clone <repository-url>
cd gort
```

Run the application:

```
go run .
```

### Configuration
You can configure the tool by editing the constants in the consts.go file before running it.

* PASSWORD_LENGTH: The fixed length for all generated passwords.

* NUM_PASSWORDS: The total number of passwords to include in the table. Increasing this number improves the chances of finding a match but also increases the size of rainbow_table.bin and the time it takes to generate it.
