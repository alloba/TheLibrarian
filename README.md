# TheLibrarian

*This project is on hold until inspiration strikes me again.*

The Librarian is an attempt at a custom file backup solution. 

All files in a provided target directory are hashed and stored in a central location. 
This can be repeated as often as needed, with each new submission of the same target directory creating a new snapshot of the files contained within. 
The benefit comes from the hashing of the targeted files - new snapshots tend to be very light due to only pulling in changed files. 

Data is tracked in a SQLite table, which means a very small footprint overall. 
One single folder to contain both the database and the files themselves, which can then be replicated on other machines. 

The codebase is largely based around the metaphor of a library, categorizing data into books (directories), pages (files), and editions (snapshots). 
It's also my first attempt at writing something in Go, beyond simple one-off scripts. 
Which is kind of my excuse for any code quality or organization issues. 
I'm much more interested in getting a feeling for the language than making anything deeply robust. 


