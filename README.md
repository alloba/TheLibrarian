# TheLibrarian
Personal File Tracker

**Super Duper Work in Progress Mode** 

# Initial Idea Statement

Comphrehensive personal data backup solution. 

This tool would be a way to keep a full backup of all files in my life. 
All saved media, all note files, all home directory configs... all of it in one place. 

The scope of the project is a little fuzzy when it comes to the border between simple file backups and additional associated behavior. 
Hopefully through user stories I can get a better idea of where to draw the line. 


## User Stories

- I want to specify certain directories on my computer for the Librarian to keep track of, and do periodic backups. 
- I want to be able to specify where the Librarian keeps its data. This will be a central location
- I want the librarian to be able to change its data storage location
- I want to be able to define multiple locations as data backups 
- I want to be able to restore the central Librarian data from backups 
- I want the Librarian to only back files up when they have changed. 
- I want the librarian to keep a history of file changes within some range, and allow restoring them. 
- I want the librarian to run quietly in the background, and automatically perform its functions. 
- I want all customization for the librarian to be stored in a standard location (XDG standard).
- I want a linear history of all operations the librarian has taken
- I want the librarian to be able to exist on multiple machines and be cross-platform.
- I want backup directories to respect relative directories (specifically the home directory).
- I want Librarian instances on separate machines to cooperate and preserve a unified data source. 
- I want to be able to specify directories that the Librarian should ignore. 
- I want the librarian to differentiate between current-machine settings and global user settings for some tasks (sync frequency/directories/ignores).
- I want the librarian to track file metadata internally.
- I want to be able to review the metadata that the librarian is keeping, in a human-readable format. 
- I want the Librarian to be able to blind-backup itself, as a final failsafe. Meaning no delta tracking or cleverness, just a pure dumb copy. 
- I want the Librarian to be secure. Encrypted files and whatnot. 
- I want to be able to add modules to the librarian that allow backing up foreign data sources (bitwarden/google photos/email/etc)
- I want to be ale to run a simplified version of the Librarian on my phone.
- I want to be able to communicate with the Librarian in order to make changes without restarting the process


